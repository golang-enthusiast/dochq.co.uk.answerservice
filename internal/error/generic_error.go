package error

// ErrInvalidArgument - invalid argument
type ErrInvalidArgument struct {
	Msg string
}

// ErrAlreadyExist - already exist.
type ErrAlreadyExist struct {
	Msg string
}

// ErrNotFound - not found.
type ErrNotFound struct {
	Msg string
}

// ErrFailedPrecondition - failed pre condition.
type ErrFailedPrecondition struct {
	Msg string
}

// ErrInternal - internal error.
type ErrInternal struct {
	Msg string
}

// ErrUnauthorized - unauthorized access.
type ErrUnauthorized struct {
	Msg string
}

// ErrPermissionDenied - permission or access denied.
type ErrPermissionDenied struct {
	Msg string
}

// NewErrInvalidArgument creates a new error.
func NewErrInvalidArgument(msg string) error {
	return &ErrInvalidArgument{msg}
}

// NewErrAlreadyExist already exist.
func NewErrAlreadyExist(msg string) error {
	return &ErrAlreadyExist{msg}
}

// NewErrNotFound not found.
func NewErrNotFound(msg string) error {
	return &ErrNotFound{msg}
}

// NewErrFailedPrecondition failed pre condition.
func NewErrFailedPrecondition(msg string) error {
	return &ErrFailedPrecondition{msg}
}

// NewErrInternal internal error.
func NewErrInternal(msg string) error {
	return &ErrInternal{msg}
}

// NewErrUnauthorized unauthorized access.
func NewErrUnauthorized(msg string) error {
	return &ErrUnauthorized{msg}
}

// NewErrPermissionDenied permission denied.
func NewErrPermissionDenied(msg string) error {
	return &ErrPermissionDenied{msg}
}

func (e *ErrInvalidArgument) Error() string {
	return e.Msg
}

func (e *ErrAlreadyExist) Error() string {
	return e.Msg
}

func (e *ErrNotFound) Error() string {
	return e.Msg
}

func (e *ErrFailedPrecondition) Error() string {
	return e.Msg
}

func (e *ErrInternal) Error() string {
	return e.Msg
}

func (e *ErrUnauthorized) Error() string {
	return e.Msg
}

func (e *ErrPermissionDenied) Error() string {
	return e.Msg
}
