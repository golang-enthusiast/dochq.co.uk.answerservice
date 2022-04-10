package domain

// Message type key.
const (
	MessageTypeAttributeKey = "MessageType"
)

// MessageType - message type.
type MessageType string

// String - string representation.
func (mt MessageType) String() string {
	return string(mt)
}

// IsValid - checks if the message type is valid.
func (mt MessageType) IsValid() bool {
	ok := validMessageTypes[mt]
	return ok
}

// IsEmpty - checks if the message type is empty.
func (mt MessageType) IsEmpty() bool {
	return len(mt) == 0
}

// Message types.
const (
	AnswerEventMessageType = MessageType("ANSWER_EVENT")
)

// List of valid message types.
var (
	validMessageTypes = map[MessageType]bool{
		AnswerEventMessageType: true,
	}
)
