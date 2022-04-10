package helpers

import (
	"os"

	"dochq.co.uk.answerservice/internal/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// GetAwsSession - returns an aws session.
func GetAwsSession() *session.Session {
	// Don't use mock server in production otherwise it will override the real aws endpoints.
	mockServerAddress := os.Getenv("AWS_MOCK_SERVER_ADDRESS")
	if len(mockServerAddress) > 0 {
		return session.Must(session.NewSession(&aws.Config{
			Endpoint:         aws.String(mockServerAddress),
			S3ForcePathStyle: aws.Bool(true),
		}))
	}
	return session.Must(session.NewSession())
}

// GetSQSQueueURL - returns the URL of an existing Amazon SQS queue.
// An error will be returned if the queue does not exist.
func GetSQSQueueURL(queueAPI domain.QueueAPI, queueName string) (url string, err error) {
	resp, err := queueAPI.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return url, err
	}
	return aws.StringValue(resp.QueueUrl), nil
}

// CreateSQSQueue - creates a new standard Amazon SQS queue.
func CreateSQSQueue(queueAPI domain.QueueAPI, queueName string) (url string, err error) {
	resp, err := queueAPI.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"VisibilityTimeout": aws.String("60"),
		},
	})
	if err != nil {
		return url, err
	}
	return aws.StringValue(resp.QueueUrl), nil
}

// GetOrCreateSQSQueue - returns an existing queue or creates a new one.
func GetOrCreateSQSQueue(queueAPI domain.QueueAPI, queueName string) (url string, err error) {
	url, err = GetSQSQueueURL(queueAPI, queueName)
	if err == nil {
		return url, nil
	}
	// The queue does not exist, create it.
	//
	return CreateSQSQueue(queueAPI, queueName)
}
