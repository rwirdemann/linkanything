package cdc

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rwirdemann/linkanything/adapter"
	"github.com/rwirdemann/linkanything/core/domain"
	"github.com/rwirdemann/linkanything/core/service"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func TestCreate(t *testing.T) {
	dbpool := newDbPool()
	defer dbpool.Close()

	linkRepository := adapter.NewPostgresLinkRepository(dbpool)
	linkService := service.NewLinkService(linkRepository)
	link, err := linkService.Create(domain.Link{
		Title: "Hello",
		URI:   "https://hello.de",
		Tags:  []string{"event"},
		Draft: true,
	})
	assert.Nil(t, err)
	assert.True(t, link.Id != 0)
	assert.True(t, link.Draft)
}

func TestPatch(t *testing.T) {
	dbpool := newDbPool()
	defer dbpool.Close()

	linkRepository := adapter.NewPostgresLinkRepository(dbpool)
	linkService := service.NewLinkService(linkRepository)
	link, err := linkService.Create(domain.Link{
		Title: "Hello",
		URI:   "https://hello.de",
		Tags:  []string{"event"},
		Draft: true,
	})

	err = linkService.Patch(domain.Patch{
		Id:    link.Id,
		Field: "draft",
		Value: "false",
	})
	assert.Nil(t, err)

	link, err = linkService.Get(link.Id)
	assert.False(t, link.Draft)
}

func newDbPool() *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}
