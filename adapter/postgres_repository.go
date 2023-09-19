package adapter

import (
	"context"
	"fmt"
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
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &PostgresRepository{connection: c}
}

func (n PostgresRepository) Create(link domain.Link) domain.Link {
	log.Printf("%v", link)
	return link
}
