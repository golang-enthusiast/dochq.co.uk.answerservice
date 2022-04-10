package sqsqueue

import (
	"context"

	"dochq.co.uk.answerservice/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(domain.QueueService) domain.QueueService

func loggingServiceMiddleware(logger log.Logger) Middleware {
	return func(next domain.QueueService) domain.QueueService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.QueueService
}

func (mw loggingMiddleware) SendMessage(
	ctx context.Context,
	queueName string,
	message domain.QueueMessage,
) (messageID string, err error) {
	defer func() {
		_ = mw.logger.Log("method", "SendMessage",
			"queueName", queueName,
			"message", message,
			"messageID", messageID,
			"err", err,
		)
	}()
	return mw.next.SendMessage(ctx, queueName, message)
}
