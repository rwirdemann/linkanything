package adapter

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rwirdemann/linkanything/core/domain"
)

type PostgresUserRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresUserRepository(dbpool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{dbpool: dbpool}
}

func (r PostgresUserRepository) Create(user domain.User) (domain.User, error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}
	_, err = r.dbpool.Exec(context.Background(), "insert into users(name,password) values($1, $2)", user.Name, hash)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
