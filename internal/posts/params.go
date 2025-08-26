package posts

type FilterParams struct {
	Limit   int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset  int    `query:"offset" validate:"omitempty,min=0"`
	Sort    string `query:"sort" validate:"omitempty,oneof=asc desc"`
	OrderBy string `query:"order_by" validate:"omitempty,oneof=created_at"`
}
