package worker

import (
	"context"

	"dochq.co.uk.answerservice/internal/domain"
	errors "dochq.co.uk.answerservice/internal/error"

	"github.com/go-kit/log"
)

// AnswerWorker struct.
type AnswerWorker struct {
	*Worker
	eventRepository domain.AnswerEventRepository
}

// NewAnswerWorker - sets up a new worker.
func NewAnswerWorker(
	props *Props,
	queueAPI domain.QueueAPI,
	eventRepository domain.AnswerEventRepository,
	logger log.Logger) *AnswerWorker {
	var (
		worker = new(props, queueAPI, logger)
	)
	return &AnswerWorker{
		Worker:          worker,
		eventRepository: eventRepository,
	}
}

// HandleAnswerEventMessage - handle message.
func (w *AnswerWorker) HandleAnswerEventMessage(
	ctx context.Context,
	m *domain.AnswerEventMessage,
) error {
	if m == nil {
		return errors.NewErrInvalidArgument("AnswerEventMessage required")
	}
	if m.Event == nil {
		return errors.NewErrInvalidArgument("Event required")
	}
	if !m.Event.EventType.IsValid() {
		return errors.NewErrInvalidArgument("EventType is not valid")
	}
	if m.Event.Data == nil {
		return errors.NewErrInvalidArgument("EventData required")
	}
	if err := m.Event.Data.Validate(); err != nil {
		return errors.NewErrInvalidArgument(err.Error())
	}
	return w.eventRepository.Create(m.Event)
}
