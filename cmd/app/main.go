package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pkgApi "dochq.co.uk.answerservice/api/generated/dochq.co.uk/answerserviceapi/v1"
	pkgAnswer "dochq.co.uk.answerservice/internal/answer"
	"dochq.co.uk.answerservice/internal/domain"
	pkgDynamodb "dochq.co.uk.answerservice/internal/dynamodb"
	pkgHelpers "dochq.co.uk.answerservice/internal/helpers"
	"dochq.co.uk.answerservice/internal/sqsqueue"

	"github.com/aws/aws-sdk-go/service/sqs"
	kitzapadapter "github.com/go-kit/kit/log/zap"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {

	var (
		answerTableName      = os.Getenv(domain.EnvAnswerTableName)
		answerEventTableName = os.Getenv(domain.EnvAnswerEventTableName)
		answerEventQueueName = os.Getenv(domain.EnvAnswerEventQueueName)
	)

	// Create a single logger, which we'll use and give to other components.
	//
	zapLogger, _ := zap.NewProduction()
	defer func() {
		_ = zapLogger.Sync()
	}()

	var logger log.Logger
	logger = kitzapadapter.NewZapSugarLogger(zapLogger, zapcore.InfoLevel)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	// Logging helper function
	logFatal := func(args ...interface{}) {
		_ = logger.Log(args...)
		os.Exit(1)
	}

	// Define our flags.
	//
	fs := flag.NewFlagSet("", flag.ExitOnError)
	grpcAddr := fs.String("grpc-addr", ":6565", "gRPC listen address")
	httpAddr := fs.String("http-addr", ":8000", "HTTP listen address")
	err := fs.Parse(os.Args[1:])
	if err != nil {
		logFatal(err)
	}

	// Setup AWS session.
	//
	awsSession := pkgHelpers.GetAwsSession()
	sqsClient := sqs.New(awsSession)

	// Repository layer.
	//
	answerRepository := pkgDynamodb.NewAnswerRepository(awsSession, answerTableName)
	answerEventRepository := pkgDynamodb.NewAnswerEventRepository(awsSession, answerEventTableName)

	// Service layer.
	//
	queueService := sqsqueue.NewQueueService(sqsClient, logger)
	answerService := pkgAnswer.NewService(answerRepository, answerEventRepository, queueService, answerEventQueueName, logger)

	// Endpoints layer.
	//
	answerEndpoints := pkgAnswer.NewEndpoint(answerService, logger)

	// GRPC Server layer.
	//
	answerGrpcServer := pkgAnswer.NewGRPCServer(answerEndpoints, logger)

	// Setup base grpc-server.
	//
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
	pkgApi.RegisterAnswerServiceServer(grpcServer, answerGrpcServer)

	// gRPC Gateway setup.
	//
	ctx := context.Background()
	rmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}))
	mux := http.NewServeMux()

	// Serve the swagger specs.
	//
	mux.Handle("/", rmux)
	mux.HandleFunc("/swagger", pkgHelpers.ServeSwagger)
	{
		err := pkgApi.RegisterAnswerServiceHandlerServer(ctx, rmux, answerGrpcServer)
		if err != nil {
			logFatal("during", "Setup", "err", err)
		}
	}
	var g run.Group
	// Startup the gRPC listener
	{
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logFatal("transport", "gRPC", "during", "Listen", "err", err)
		}

		g.Add(func() error {
			_ = logger.Log("transport", "gRPC", "addr", *grpcAddr)
			return grpcServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	// Startup the HTTP listener
	{
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logFatal("transport", "HTTP", "during", "Listen", "err", err)
		}

		g.Add(func() error {
			_ = logger.Log("transport", "HTTP", "addr", *httpAddr)
			return http.Serve(httpListener, mux)
		}, func(err error) {
			httpListener.Close()
		})
	}
	// This function just sits and waits for ctrl-C.
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	_ = logger.Log("exit", g.Run())
}
