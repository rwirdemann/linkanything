package adapter

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rwirdemann/linkanything/core/domain"
	log "github.com/sirupsen/logrus"
	"strings"
)

type PostgresLinkRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresLinkRepository(dbpool *pgxpool.Pool) *PostgresLinkRepository {
	return &PostgresLinkRepository{dbpool: dbpool}
}

func (r PostgresLinkRepository) Create(link domain.Link) (domain.Link, error) {
	err := r.dbpool.QueryRow(context.Background(),
		"insert into links(title,uri,draft,tags) values($1, $2, $3, $4) RETURNING id",
		link.Title, link.URI, link.Draft, strings.Join(lower(link.Tags), ",")).Scan(&link.Id)
	if err != nil {
		return domain.Link{}, err
	}
	return link, nil
}

func (r PostgresLinkRepository) Update(link domain.Link) (domain.Link, error) {
	_, err := r.dbpool.Exec(context.Background(), "update links SET (title,uri,draft,tags) = ($1, $2, $3, $4) where id=$5", link.Title, link.URI, link.Draft, strings.Join(lower(link.Tags), ","), link.Id)
	if err != nil {
		return domain.Link{}, err
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

func (r PostgresLinkRepository) GetLinks(tagList []string) ([]domain.Link, error) {
	var rows pgx.Rows
	var err error
	if len(tagList) > 0 {
		rows, err = r.dbpool.Query(context.Background(),
			"select id, title, uri, created, tags from links where tags like $1 order by created desc", "%"+tagList[0]+"%")
	} else {
		rows, err = r.dbpool.Query(context.Background(), "select id, title, uri, created, tags from links order by created desc")
	}

	if err != nil {
		log.Error(err)
		return []domain.Link{}, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		log.Error(err)
		return []domain.Link{}, err
	}

	var links []domain.Link
	var tags string
	for rows.Next() {
		log.Printf("GetLinks: Adding row to return array")
		var l domain.Link
		err := rows.Scan(&l.Id, &l.Title, &l.URI, &l.Created, &tags)
		if len(tags) > 0 {
			l.Tags = strings.Split(tags, ",")
		}
		if err != nil {
			return []domain.Link{}, err
		}
		links = append(links, l)
	}
	return links, nil
}

func (r PostgresLinkRepository) Get(id int) (domain.Link, error) {
	var l domain.Link
	var tags string
	err := r.dbpool.QueryRow(context.Background(), "SELECT id, title, uri, created, tags from links where id=$1", id).Scan(&l.Id, &l.Title, &l.URI, &l.Created, &tags)
	if len(tags) > 0 {
		l.Tags = strings.Split(tags, ",")
	}
	return l, err
}
