package repository

import (
	"github.com/mdanialr/sns_backend/pkg/pagination"
	"gorm.io/gorm"
)

// IOptions signature that should be used to additionally add query to each
// repository layer implementations.
type IOptions interface {
	Set(*gorm.DB) *gorm.DB
}

type (
	columns struct{ cols []string }
	order   struct{ order []string }
	where   struct{ cons []string }
)

func (c *columns) Set(db *gorm.DB) *gorm.DB { return db.Select(c.cols) }

func (o *order) Set(db *gorm.DB) *gorm.DB {
	for _, ord := range o.order {
		db = db.Order(ord)
	}
	return db
}

func (w *where) Set(db *gorm.DB) *gorm.DB {
	for _, con := range w.cons {
		db = db.Where(con)
	}
	return db
}

// Cols add query Select.
// Example:
//
//	repository.Cols("id","created_at","updated_at")
func Cols(cols ...string) IOptions { return &columns{cols} }

// Order add query Order.
//
// Example:
//
//	repository.Order("created_at DESC")
func Order(orders ...string) IOptions { return &order{orders} }

// Cons add query Where for each given cons. Each given conditions will be
// combined by GORM using AND.
//
// Example:
//
//	repository.Cons("id IS NULL"), repository.Cons("name IS NOT NULL")
func Cons(cons ...string) IOptions { return &where{cons} }

// Paginate add query Limit & Offset accordingly by given paginate.M.
//
// Example:
//
//	repository.Paginate(&paginate.M{Limit: 4})
func Paginate(p *paginate.M) IOptions {
	if p == nil {
		return new(paginate.M)
	}
	return p
}
