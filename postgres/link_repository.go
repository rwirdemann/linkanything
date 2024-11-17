package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rwirdemann/linkanything/core"
	"strings"
)

type LinkRepository struct {
	dbpool  *pgxpool.Pool
	queries *Queries
}

func NewPostgresLinkRepository(dbpool *pgxpool.Pool) *LinkRepository {
	return &LinkRepository{dbpool: dbpool, queries: New(dbpool)}
}

func (r LinkRepository) Create(link core.Link) (core.Link, error) {
	err := r.dbpool.QueryRow(context.Background(),
		"insert into links(title,uri,draft,tags) values($1, $2, $3, $4) RETURNING id",
		link.Title, link.URI, link.Draft, strings.Join(lower(link.Tags), ",")).Scan(&link.Id)
	if err != nil {
		return core.Link{}, err
	}
	return link, nil
}

func (r LinkRepository) Update(link core.Link) (core.Link, error) {
	_, err := r.dbpool.Exec(context.Background(), "update links SET (title,uri,draft,tags) = ($1, $2, $3, $4) where id=$5", link.Title, link.URI, link.Draft, strings.Join(lower(link.Tags), ","), link.Id)
	if err != nil {
		return core.Link{}, err
	}
	return link, nil
}

func lower(tags []string) []string {
	var result []string
	for _, t := range tags {
		result = append(result, strings.ToLower(t))
	}
	return result
}

func (r LinkRepository) GetLinks(tagList []string, includeDrafts bool, page, limit int) ([]core.Link, error) {
	offset := 0
	if page > 0 && limit > 0 {
		offset = (page - 1) * limit
	}
	if limit == 0 {
		limit = 1000
	}

	rows, err := r.queries.GetLinks(context.Background(), GetLinksParams{
		Column1: tagList,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return []core.Link{}, err
	}

	var links []core.Link
	for _, row := range rows {
		links = append(links, core.Link{
			Id:      int(row.ID),
			Title:   row.Title,
			URI:     row.Uri,
			Created: row.Created.Time,
			Draft:   false,
		})
	}
	return links, nil
}

func (r LinkRepository) Get(id int) (core.Link, error) {
	var l core.Link
	var tags string
	err := r.dbpool.QueryRow(context.Background(), "SELECT id, title, uri, created, tags from links where id=$1", id).Scan(&l.Id, &l.Title, &l.URI, &l.Created, &tags)
	if len(tags) > 0 {
		l.Tags = strings.Split(tags, ",")
	}
	return l, err
}
