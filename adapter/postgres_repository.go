package adapter

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/rwirdemann/linkanything/core/domain"
	"log"
	"os"
	"strings"
)

type PostgresRepository struct {
	connection *pgx.Conn
}

func NewPostgresRepository() *PostgresRepository {
	c, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return &PostgresRepository{connection: c}
}

func (r PostgresRepository) Create(link domain.Link) (domain.Link, error) {
	_, err := r.connection.Exec(context.Background(), "insert into links(title,uri,draft,tags) values($1, $2, $3, $4)", link.Title, link.URI, link.Draft, strings.Join(link.Tags, ","))
	if err != nil {
		return domain.Link{}, err
	}
	return link, nil
}

func (r PostgresRepository) GetLinks() ([]domain.Link, error) {
	rows, err := r.connection.Query(context.Background(), "select id, title, uri, created from links order by created desc")
	if err != nil {
		return []domain.Link{}, err
	}
	defer rows.Close()

	var links []domain.Link
	for rows.Next() {
		var l domain.Link
		err := rows.Scan(&l.Id, &l.Title, &l.URI, &l.Created)
		if err != nil {
			return []domain.Link{}, err
		}
		links = append(links, l)
	}
	return links, nil
}
