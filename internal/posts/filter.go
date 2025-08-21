package posts

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FilterParams struct {
	Limit  int  `query:"limit" validate:"omitempty,min=1,max=100"`
	Offset int  `query:"offset" validate:"omitempty,min=0"`
	Desc   bool `query:"desc" validate:"omitempty,boolean"`
}

func DefaultFilterParams(f FilterParams) FilterParams {
	fp := FilterParams{
		Limit:  20,
		Offset: 0,
	}

	if f.Limit > 0 {
		fp.Limit = f.Limit
	}

	if f.Offset > 0 {
		fp.Offset = f.Offset
	}

	if f.Desc {
		fp.Desc = f.Desc
	}

	return fp
}

func (f *FilterParams) Apply(db *gorm.DB) *gorm.DB {
	query := db

	query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: f.Desc})

	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}

	if f.Offset > 0 {
		query = query.Offset(f.Offset)
	}

	return query
}
