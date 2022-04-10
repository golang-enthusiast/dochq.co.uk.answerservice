package domain

// AnswerEventType - answer event type.
type AnswerEventType string

// IsValid - checks if the event type is valid.
func (t AnswerEventType) IsValid() bool {
	ok := validAnswerEventTypes[t]
	return ok
}

// Event types.
const (
	CreateAnswerEventType = AnswerEventType("create")
	UpdateAnswerEventType = AnswerEventType("update")
	DeleteAnswerEventType = AnswerEventType("delete")
)

// JSON fields.
const (
	JSONFieldEventType = "eventType"
)

// List of valid answer event types.
var (
	validAnswerEventTypes = map[AnswerEventType]bool{
		CreateAnswerEventType: true,
		UpdateAnswerEventType: true,
		DeleteAnswerEventType: true,
	}
)

// AnswerEvent - represents an answer event struct.
type AnswerEvent struct {
	EventType AnswerEventType `json:"eventType"`
	Data      *Answer         `json:"data"`
}

// AnswerEventRepository - provides access to a storage.
type AnswerEventRepository interface {
	Create(answerEvent *AnswerEvent) error
}
