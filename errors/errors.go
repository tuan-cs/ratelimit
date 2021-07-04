package errors

type Error string

func (e Error) Error() string {
	return string(e)
}

var (
	ErrTimeout = Error("timeout white executing")

	ErrContextCanceled = Error("context canceled, logic not executed")

	ErrTimeoutWaitingExecution = Error("timeout while waiting for execution")
)
