package domain

import (
	"context"
	"errors"
)

// Answer JSON fields.
const (
	JSONFieldAnswerKey   = "key"
	JSONFieldAnswerValue = "value"
)

// Answer key/value types.
type (
	AnswerKey   string
	AnswerValue string
)

// Answer - represents an answer struct.
type Answer struct {
	Key   AnswerKey   `json:"key"`
	Value AnswerValue `json:"value"`
}

// Validate - validates struct.
func (answer *Answer) Validate() error {
	if len(answer.Key) == 0 {
		return errors.New("Key required")
	}
	if len(answer.Value) == 0 {
		return errors.New("Value required")
	}
	return nil
}

// AnswerRepository - provides access to a storage.
type AnswerRepository interface {
	Create(answer *Answer) error
	Update(answer *Answer) error
	Delete(key AnswerKey) error
	Get(key AnswerKey) (*Answer, error)
}

// AnswerService - provides access to a business logic.
type AnswerService interface {

	// CreateAnswer - creates a new answer.
	CreateAnswer(ctx context.Context, answer *Answer) error

	// UpdateAnswer - updates an existing answer.
	UpdateAnswer(ctx context.Context, answer *Answer) error

	// DeleteAnswer - deletes an existing answer.
	DeleteAnswer(ctx context.Context, key AnswerKey) error

	// GetAnswer - returns an existing answer by the provided key.
	GetAnswer(ctx context.Context, key AnswerKey) (*Answer, error)
}
