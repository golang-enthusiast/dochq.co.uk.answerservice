package domain

import (
	"context"

	"github.com/aws/aws-sdk-go/service/sqs"
)

// QueueAPI - is the minimum interface required from a queue implementation.
type QueueAPI interface {
	CreateQueue(*sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error)
	GetQueueUrl(*sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error)
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

// QueueHandlerFunc is used to define the Handler that is run on for each message.
type QueueHandlerFunc func(ctx context.Context, msg *sqs.Message) error

// HandleMessage wraps a function for handling sqs messages.
func (f QueueHandlerFunc) HandleMessage(ctx context.Context, msg *sqs.Message) error {
	return f(ctx, msg)
}

// QueueHandler interface.
type QueueHandler interface {
	HandleMessage(ctx context.Context, msg *sqs.Message) error
}

// QueueService - high-level service that provides access to the queue.
type QueueService interface {

	// SendMessage - sends a message to the queue.
	SendMessage(ctx context.Context, queueName string, message QueueMessage) (messageID string, err error)
}
