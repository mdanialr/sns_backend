package database

import (
	"context"
	"database/sql"
)

// SNS signatures that related to all database operation about Shorten and Send.
type SNS interface {
	Querier
	ListShorten(ctx context.Context, colName string, order DBOrder) ([]Shorten, error)
	UpdateShorten(ctx context.Context, arg UpdateShortenParams) (Shorten, error)
}

// DB provides all functionality to interacting with database.
type DB struct {
	db *sql.DB
	SNS
}

// NewSNS creates a new SNS interface instance.
func NewSNS(db *sql.DB) SNS {
	return &DB{
		db:  db,
		SNS: New(db),
	}
}
