package answer

import (
	"context"

	"dochq.co.uk.answerservice/internal/domain"
	errors "dochq.co.uk.answerservice/internal/error"

	"github.com/go-kit/log"
)

type service struct {
	repository      domain.AnswerRepository
	eventRepository domain.AnswerEventRepository
	queueService    domain.QueueService
	eventQueueName  string
}

// NewService creates a new service with necessary dependencies.
func NewService(repository domain.AnswerRepository,
	eventRepository domain.AnswerEventRepository,
	queueService domain.QueueService,
	eventQueueName string,
	logger log.Logger) domain.AnswerService {
	var service domain.AnswerService
	{
		service = newBasicService(repository, eventRepository, queueService, eventQueueName)
		service = LoggingServiceMiddleware(logger)(service)
	}
	return service
}

// Returns a naive, stateless implementation of service.
func newBasicService(repository domain.AnswerRepository,
	eventRepository domain.AnswerEventRepository,
	queueService domain.QueueService,
	eventQueueName string) domain.AnswerService {
	return &service{
		repository:      repository,
		eventRepository: eventRepository,
		queueService:    queueService,
		eventQueueName:  eventQueueName,
	}
}

func (s *service) CreateAnswer(ctx context.Context, answer *domain.Answer) error {

	// Check for nil.
	//
	if answer == nil {
		return errors.NewErrInvalidArgument("Answer required")
	}

	// Validate fields.
	//
	if err := answer.Validate(); err != nil {
		return errors.NewErrInvalidArgument(err.Error())
	}

	// If the answer already exist, we must return an error.
	//
	foundAnswer, _ := s.repository.Get(answer.Key)
	if foundAnswer != nil {
		return errors.NewErrAlreadyExist("Answer with the provided key is already in use")
	}

	// Create answer.
	//
	if err := s.repository.Create(answer); err != nil {
		return err
	}

	// Send event message.
	//
	_, err := s.queueService.SendMessage(ctx, s.eventQueueName, &domain.AnswerEventMessage{
		Event: &domain.AnswerEvent{
			EventType: domain.CreateAnswerEventType,
			Data:      answer,
		},
	})
	return err
}

func (s *service) UpdateAnswer(ctx context.Context, answer *domain.Answer) error {

	// Check for nil.
	//
	if answer == nil {
		return errors.NewErrInvalidArgument("Answer required")
	}

	// Validate fields.
	//
	if err := answer.Validate(); err != nil {
		return errors.NewErrInvalidArgument(err.Error())
	}

	// If the answer does not exist, we must return an error.
	//
	foundAnswer, _ := s.repository.Get(answer.Key)
	if foundAnswer == nil {
		return errors.NewErrNotFound("Answer with the provided key not found")
	}

	// Update answer.
	//
	if err := s.repository.Update(answer); err != nil {
		return err
	}

	// Send event message.
	//
	_, err := s.queueService.SendMessage(ctx, s.eventQueueName, &domain.AnswerEventMessage{
		Event: &domain.AnswerEvent{
			EventType: domain.UpdateAnswerEventType,
			Data:      answer,
		},
	})
	return err
}

func (s *service) DeleteAnswer(ctx context.Context, key domain.AnswerKey) error {

	// Check key.
	//
	if len(key) == 0 {
		return errors.NewErrInvalidArgument("AnswerKey required")
	}

	// If the answer does not exist, we must return an error.
	//
	foundAnswer, _ := s.repository.Get(key)
	if foundAnswer == nil {
		return errors.NewErrNotFound("Answer with the provided key not found")
	}

	// Delete answer.
	//
	if err := s.repository.Delete(key); err != nil {
		return err
	}

	// Send event message.
	//
	_, err := s.queueService.SendMessage(ctx, s.eventQueueName, &domain.AnswerEventMessage{
		Event: &domain.AnswerEvent{
			EventType: domain.DeleteAnswerEventType,
			Data:      foundAnswer,
		},
	})
	return err
}

func (s *service) GetAnswer(ctx context.Context, key domain.AnswerKey) (*domain.Answer, error) {

	// Check key.
	//
	if len(key) == 0 {
		return nil, errors.NewErrInvalidArgument("AnswerKey required")
	}

	// Return result.
	//
	return s.repository.Get(key)
}

func (s *service) GetAnswerHistory(ctx context.Context, key domain.AnswerKey) ([]*domain.AnswerEvent, error) {

	// Check key.
	//
	if len(key) == 0 {
		return nil, errors.NewErrInvalidArgument("AnswerKey required")
	}

	// Return result.
	//
	return s.eventRepository.ListEvents(key)
}
