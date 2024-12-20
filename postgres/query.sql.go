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
select distinct l.Id, l.title, l.uri, l.created, l.draft
from links l
         left join public.tags_links tl on l.id = tl.link_id
         left join public.tags tag on tag.id = tl.tag_id
where tag.name = ANY($1::varchar[]) or $1 is NULL
order by created desc
limit $2 offset $3
`

type GetLinksParams struct {
	Column1 []string
	Limit   int32
	Offset  int32
}

type GetLinksRow struct {
	ID      int32
	Title   string
	Uri     string
	Created pgtype.Timestamptz
	Draft   pgtype.Bool
}

func (q *Queries) GetLinks(ctx context.Context, arg GetLinksParams) ([]GetLinksRow, error) {
	rows, err := q.db.Query(ctx, getLinks, arg.Column1, arg.Limit, arg.Offset)
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

const linkCount = `-- name: LinkCount :one
select count(*) from links
`

func (q *Queries) LinkCount(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, linkCount)
	var count int64
	err := row.Scan(&count)
	return count, err
}
