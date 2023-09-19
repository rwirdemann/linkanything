package adapter

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/rwirdemann/linkanything/domain"
	"os"
)

type LinkRepository struct {
	connection *pgx.Conn
}

func NewLinkRepository() *LinkRepository {
	c, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &LinkRepository{connection: c}
}

func (n LinkRepository) Create(link domain.Link) domain.Link {
	fmt.Printf("%v", link)
	return link
}
