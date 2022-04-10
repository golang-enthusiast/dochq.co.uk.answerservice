package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"dochq.co.uk.answerservice/internal/domain"
	pkgDynamodb "dochq.co.uk.answerservice/internal/dynamodb"
	errors "dochq.co.uk.answerservice/internal/error"
	pkgHelpers "dochq.co.uk.answerservice/internal/helpers"
	workers "dochq.co.uk.answerservice/internal/worker"

	"github.com/aws/aws-sdk-go/service/sqs"
	kitzapadapter "github.com/go-kit/kit/log/zap"
	"github.com/go-kit/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {

	var (
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

	// Setup AWS services.
	//
	awsSession := pkgHelpers.GetAwsSession()
	sqsClient := sqs.New(awsSession)

	// Setup repository.
	//
	eventRepository := pkgDynamodb.NewAnswerEventRepository(awsSession, answerEventTableName)

	// Define worker properties.
	//
	workerProps := &workers.Props{
		WorkerName:          "answer-event-worker",
		QueueName:           answerEventQueueName,
		MaxNumberOfMessages: 10,
	}

	w := workers.NewAnswerWorker(
		workerProps,
		sqsClient,
		eventRepository,
		logger,
	)
	w.Start(context.Background(), domain.QueueHandlerFunc(func(ctx context.Context, msg *sqs.Message) error {
		// Define message type.
		//
		messageTypeValue, ok := msg.MessageAttributes[domain.MessageTypeAttributeKey]
		if !ok {
			return errors.NewErrNotFound("Message type not found in attributes")
		}
		messageType := domain.MessageType(*messageTypeValue.StringValue)

		switch messageType {
		case domain.AnswerEventMessageType:
			// Unmarshal payload.
			//
			payload := &domain.AnswerEventMessage{}
			err := json.Unmarshal([]byte(*msg.Body), payload)
			if err != nil {
				return err
			}
			return w.HandleAnswerEventMessage(ctx, payload)
		default:
			return errors.NewErrNotFound(fmt.Sprintf("Unsupported message type %v", messageType))
		}
	}))
}
