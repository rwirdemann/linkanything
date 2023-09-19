package adapter

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/rwirdemann/linkanything/core/domain"
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
	_, err := r.connection.Exec(context.Background(), "insert into links(title,uri) values($1, $2)", link.Title, link.URI)
	if err != nil {
		return domain.Link{}, err
	}
	return link, nil
}
