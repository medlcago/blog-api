package errors

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e Error) Error() string {
	return e.Msg
}

func New(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}
