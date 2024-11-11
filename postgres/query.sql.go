// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getLinks = `-- name: GetLinks :many
select Id, Title, uri, created, tags, draft from links
order by created desc limit $1 offset $2
`

type GetLinksParams struct {
	Limit  int32
	Offset int32
}

type GetLinksRow struct {
	ID      int32
	Title   string
	Uri     string
	Created pgtype.Timestamptz
	Tags    pgtype.Text
	Draft   pgtype.Bool
}

func (q *Queries) GetLinks(ctx context.Context, arg GetLinksParams) ([]GetLinksRow, error) {
	rows, err := q.db.Query(ctx, getLinks, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLinksRow
	for rows.Next() {
		var i GetLinksRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Uri,
			&i.Created,
			&i.Tags,
			&i.Draft,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}