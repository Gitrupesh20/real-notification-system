package error

type AppError struct {
	Msg string
	Err string
}

func (e *AppError) Error() string {
	return e.Msg + ": " + e.Err
}

var NoConnectionError = &AppError{Err: "Error while connecting to RabbitMQ"}

func NewAppError(msg string) *AppError {
	return &AppError{Err: msg}
}
