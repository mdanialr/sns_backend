package database

import (
	"context"
	"fmt"
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
