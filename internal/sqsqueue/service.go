package sqsqueue

import (
	"context"
	"encoding/json"
	"sync"

	"dochq.co.uk.answerservice/internal/domain"
	errors "dochq.co.uk.answerservice/internal/error"
	"dochq.co.uk.answerservice/internal/helpers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kit/log"
)

type service struct {
	queueAPI    domain.QueueAPI
	queueURLMap sync.Map
}

// NewQueueService creates a service with necessary dependencies.
func NewQueueService(queueAPI domain.QueueAPI, logger log.Logger) domain.QueueService {
	var service domain.QueueService
	{
		service = newBasicQueueService(queueAPI)
		service = loggingServiceMiddleware(logger)(service)
	}
	return service
}

func newBasicQueueService(queueAPI domain.QueueAPI) domain.QueueService {
	return &service{
		queueAPI: queueAPI,
	}
}

func (s *service) SendMessage(
	ctx context.Context,
	queueName string,
	message domain.QueueMessage,
) (emptyMessageID string, err error) {

	// Check if the queue name defined.
	//
	if len(queueName) == 0 {
		return emptyMessageID,
			errors.NewErrInvalidArgument("Queue name required")
	}

	// Check message for nil.
	//
	if message == nil {
		return emptyMessageID,
			errors.NewErrInvalidArgument("Message cannot be nil")
	}

	// Validate message.
	//
	if err := message.Validate(); err != nil {
		return emptyMessageID,
			errors.NewErrInvalidArgument(err.Error())
	}

	// Marshal message body to JSON.
	//
	messageBody, err := json.Marshal(message)
	if err != nil {
		return emptyMessageID, err
	}

	// Get queue url.
	//
	queueURL, err := s.getOrCreateQueueURL(queueName)
	if err != nil {
		return emptyMessageID, err
	}

	// Send message.
	//
	resp, err := s.queueAPI.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(messageBody)),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			domain.MessageTypeAttributeKey: {
				DataType:    aws.String("String"),
				StringValue: aws.String(message.GetMessageType().String()),
			},
		},
	})
	if err != nil {
		return emptyMessageID, err
	}

	// Return result.
	//
	return aws.StringValue(resp.MessageId), nil
}

func (s *service) getOrCreateQueueURL(queueName string) (string, error) {
	if queueURL, ok := s.queueURLMap.Load(queueName); ok {
		return queueURL.(string), nil
	}
	queueURL, err := helpers.GetOrCreateSQSQueue(s.queueAPI, queueName)
	if err != nil {
		return "", err
	}
	s.queueURLMap.Store(queueName, queueURL)
	return queueURL, nil
}
