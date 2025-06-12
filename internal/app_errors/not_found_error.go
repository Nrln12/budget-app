package app_errors

type NotFoundError struct {
	Message string
}

func (receiver NotFoundError) Error() string {
	return receiver.Message
}

func NewNotFoundError(message string) NotFoundError {
	if message == "" {
		message = "Resource Not Found"
	}
	return NotFoundError{
		Message: message,
	}
}
