// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Link struct {
	ID      int32
	Title   string
	Uri     string
	Tags    pgtype.Text
	Draft   pgtype.Bool
	Created pgtype.Timestamptz
}

type Tag struct {
	ID   int32
	Name string
}

type TagsLink struct {
	TagID  pgtype.Int4
	LinkID pgtype.Int4
}

type User struct {
	ID       int32
	Name     string
	Password string
	Created  pgtype.Timestamptz
}
