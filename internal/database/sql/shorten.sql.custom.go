package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// DBOrder just a base type to define as enums for ordering query result.
type DBOrder string

const (
	ASC         DBOrder = "ASC"  // enums type for Ascending
	DESC        DBOrder = "DESC" // enums type for Descending
	listShorten         = `-- name: ListShorten :many
SELECT id, url, target, permanent, created_at, updated_at FROM Shorten
ORDER BY
`
)

// ListShorten custom implementation that take additional arguments which are column name and order of the result.
func (q *Queries) ListShorten(ctx context.Context, colName string, order DBOrder) ([]Shorten, error) {
	stmt := fmt.Sprintf("%s %s %s", listShorten, colName, order)
	rows, err := q.db.QueryContext(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Shorten
	for rows.Next() {
		var i Shorten
		if err := rows.Scan(
			&i.ID,
			&i.Url,
			&i.Target,
			&i.Permanent,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateShorten = `-- name: UpdateShorten :one
UPDATE shorten
SET url = $1,
	target = $2,
	permanent = $3,
	updated_at = $4
WHERE id = #1
RETURNING id, url, target, permanent, created_at, updated_at
`

type UpdateShortenParams struct {
	ID        int64        `json:"id"`
	Url       string       `json:"url"`
	Target    string       `json:"target"`
	Permanent bool         `json:"permanent"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

func (q *Queries) UpdateShorten(ctx context.Context, arg UpdateShortenParams) (Shorten, error) {
	// replace #1 with arg.ID
	idStr := strconv.FormatInt(arg.ID, 10)
	newUpdateShorten := strings.Replace(updateShorten, "#1", idStr, 1)
	row := q.db.QueryRowContext(ctx, newUpdateShorten,
		arg.Url,
		arg.Target,
		arg.Permanent,
		arg.UpdatedAt,
	)
	var i Shorten
	err := row.Scan(
		&i.ID,
		&i.Url,
		&i.Target,
		&i.Permanent,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
