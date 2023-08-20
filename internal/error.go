package internal

type APIError struct {
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e APIError) Error() string {
	return e.Message
}
