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

type Pagination[T any] struct {
	Total  int64 `json:"total"`
	Result []T   `json:"result"`
}

func NewResponse[T any](data T) Response[T] {
	return Response[T]{
		OK:   true,
		Data: data,
	}
}

func NewPaginatedResponse[T any](total int64, result []T) Response[Pagination[T]] {
	return NewResponse(Pagination[T]{
		Total:  total,
		Result: result,
	})
}
