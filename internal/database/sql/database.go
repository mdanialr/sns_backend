package database

import "database/sql"

// SNS signatures that related to all database operation about Shorten and Send.
type SNS interface {
	Querier
}

// DB provides all functionality to interacting with database.
type DB struct {
	db *sql.DB
	*Queries
}

// NewSNS creates a new SNS interface instance.
func NewSNS(db *sql.DB) SNS {
	return &DB{
		db:      db,
		Queries: New(db),
	}
}
