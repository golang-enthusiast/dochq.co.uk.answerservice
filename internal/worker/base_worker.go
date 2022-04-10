package worker

import (
	"context"
	"sync"

	"dochq.co.uk.answerservice/internal/domain"
	errors "dochq.co.uk.answerservice/internal/error"
	"dochq.co.uk.answerservice/internal/helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kit/log"
)

// Props struct.
type Props struct {
	WorkerName          string
	QueueName           string
	MaxNumberOfMessages int64
}

// Worker struct.
type Worker struct {
	Props       *Props
	QueueAPI    domain.QueueAPI
	Logger      log.Logger
	QueueURLMap sync.Map
}

// new - sets up a new worker.
func new(props *Props, queueAPI domain.QueueAPI, logger log.Logger) *Worker {
	return &Worker{
		Props:    props,
		QueueAPI: queueAPI,
		Logger:   logger,
	}
}

// Start - starts the polling and will continue polling till the application is forcibly stopped.
func (worker *Worker) Start(ctx context.Context, h domain.QueueHandler) {
	for {
		select {
		case <-ctx.Done():
			_ = worker.Logger.Log("Stopping polling because a context kill signal was sent")
			return
		default:

			// Get queue url.
			//
			queueURL, err := worker.getOrCreateQueueURL(worker.Props.QueueName)
			if err != nil {
				_ = worker.Logger.Log("err", err.Error())
				continue
			}

			// Setup worker parameters.
			//
			params := &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL), // Required
				MaxNumberOfMessages: aws.Int64(worker.Props.MaxNumberOfMessages),
				AttributeNames: []*string{
					aws.String("All"), // Required
				},
				MessageAttributeNames: []*string{
					aws.String("All"), // Required
				},
			}

			// Receive message from queue.
			//
			resp, err := worker.QueueAPI.ReceiveMessage(params)
			if err != nil {
				_ = worker.Logger.Log("err", err.Error())
				continue
			}
			if len(resp.Messages) > 0 {
				worker.run(ctx, h, resp.Messages)
			}
		}
	}
}

// run - launches goroutine per received message and wait for all message to be processed.
func (worker *Worker) run(ctx context.Context, h domain.QueueHandler, messages []*sqs.Message) {
	numMessages := len(messages)
	_ = worker.Logger.Log("Received messages", numMessages)

	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *sqs.Message) {
			// launch goroutine
			defer wg.Done()

			// Hadle message
			//
			err := worker.handleMessage(ctx, m, h)
			if err != nil {
				_ = worker.Logger.Log("Failed to handle message", "err", err)
			}
		}(messages[i])
	}

	wg.Wait()
}

func (worker *Worker) handleMessage(ctx context.Context, m *sqs.Message, h domain.QueueHandler) (err error) {

	// Handle message.
	//
	err = h.HandleMessage(ctx, m)
	if _, ok := err.(*errors.ErrInvalidArgument); ok {
		_ = worker.Logger.Log("err", err.Error())
	} else if err != nil {
		return err
	}

	// Get queue url.
	//
	queueURL, err := worker.getOrCreateQueueURL(worker.Props.QueueName)
	if err != nil {
		return err
	}

	// Delete message.
	//
	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL), // Required
		ReceiptHandle: m.ReceiptHandle,      // Required
	}
	_, err = worker.QueueAPI.DeleteMessage(params)
	if err != nil {
		return err
	}
	return nil
}

func (worker *Worker) getOrCreateQueueURL(queueName string) (string, error) {
	if queueURL, ok := worker.QueueURLMap.Load(queueName); ok {
		return queueURL.(string), nil
	}
	queueURL, err := helpers.GetOrCreateSQSQueue(worker.QueueAPI, queueName)
	if err != nil {
		return "", err
	}
	worker.QueueURLMap.Store(queueName, queueURL)
	return queueURL, nil
}
