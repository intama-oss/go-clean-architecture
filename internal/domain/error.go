package domain

type Error struct {
	Code    int      `json:"code"`
	Errors  []string `json:"errors,omitempty"`
	Message string   `json:"message"`
}

func NewError(code int, message string) Error {
	return Error{Code: code, Message: message}
}

func (e Error) Error() string {
	return e.Message
}
