package adapter

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rwirdemann/linkanything/core/domain"
	log "github.com/sirupsen/logrus"
	"strings"
)

type PostgresRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresRepository(dbpool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{dbpool: dbpool}
}

func (r PostgresRepository) Create(link domain.Link) (domain.Link, error) {
	_, err := r.dbpool.Exec(context.Background(), "insert into links(title,uri,draft,tags) values($1, $2, $3, $4)", link.Title, link.URI, link.Draft, strings.Join(lower(link.Tags), ","))
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

func (r PostgresRepository) GetLinks(tagList []string) ([]domain.Link, error) {
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
