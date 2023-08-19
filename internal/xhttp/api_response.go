package xhttp

type APIError struct {
	Message string `json:"message"`
	status  int
}

func NewAPIError(message string, statusCode int) *APIError {
	return &APIError{
		Message: message,
		status:  statusCode,
	}
}

func (e APIError) Error() string {
	return e.Message
}

func (e APIError) Status() int {
	return e.status
}
