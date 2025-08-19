package response

type Response[T any] struct {
	OK   bool   `json:"ok"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func (r *Response[T]) Error() string {
	if r.OK {
		return ""
	}
	return r.Msg
}
