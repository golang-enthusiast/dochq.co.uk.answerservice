package domain

import (
	"errors"
)

// QueueMessage interface.
type QueueMessage interface {

	// GetMessageType - returns message type.
	GetMessageType() MessageType

	// Validate - validates message type.
	Validate() error
}

// AnswerEventMessage - event message.
type AnswerEventMessage struct {
	Event *AnswerEvent
}

// GetMessageType - returns AnswerEventMessageType.
func (aem *AnswerEventMessage) GetMessageType() MessageType {
	return AnswerEventMessageType
}

// Validate - validates message type.
func (aem *AnswerEventMessage) Validate() error {
	if !aem.GetMessageType().IsValid() {
		return errors.New("Message type is not valid")
	}

	return nil
}

var (
	_ QueueMessage = &AnswerEventMessage{}
)
