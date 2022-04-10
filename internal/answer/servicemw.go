package answer

import (
	"context"

	"dochq.co.uk.answerservice/internal/domain"

	"github.com/go-kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(domain.AnswerService) domain.AnswerService

// LoggingServiceMiddleware takes a logger as a dependency
// and returns a service Middleware.
func LoggingServiceMiddleware(logger log.Logger) Middleware {
	return func(next domain.AnswerService) domain.AnswerService {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   domain.AnswerService
}

func (mw loggingMiddleware) CreateAnswer(ctx context.Context, answer *domain.Answer) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "CreateAnswer",
			"answer", answer,
			"err", err,
		)
	}()
	return mw.next.CreateAnswer(ctx, answer)
}

func (mw loggingMiddleware) UpdateAnswer(ctx context.Context, answer *domain.Answer) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "UpdateAnswer",
			"answer", answer,
			"err", err,
		)
	}()
	return mw.next.UpdateAnswer(ctx, answer)
}

func (mw loggingMiddleware) DeleteAnswer(ctx context.Context, key domain.AnswerKey) (err error) {
	defer func() {
		_ = mw.logger.Log("method", "DeleteAnswer",
			"key", key,
			"err", err,
		)
	}()
	return mw.next.DeleteAnswer(ctx, key)
}

func (mw loggingMiddleware) GetAnswer(ctx context.Context, key domain.AnswerKey) (found *domain.Answer, err error) {
	defer func() {
		_ = mw.logger.Log("method", "GetAnswer",
			"key", key,
			"found", found,
			"err", err,
		)
	}()
	return mw.next.GetAnswer(ctx, key)
}

func (mw loggingMiddleware) GetAnswerHistory(ctx context.Context, key domain.AnswerKey) (list []*domain.AnswerEvent, err error) {
	defer func() {
		_ = mw.logger.Log("method", "GetAnswerHistory",
			"key", key,
			"list", list,
			"err", err,
		)
	}()
	return mw.next.GetAnswerHistory(ctx, key)
}
