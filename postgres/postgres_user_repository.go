package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rwirdemann/linkanything"
	log "github.com/sirupsen/logrus"
)

type PostgresUserRepository struct {
	dbpool *pgxpool.Pool
}

func NewPostgresUserRepository(dbpool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{dbpool: dbpool}
}

func (r PostgresUserRepository) Create(user linkanything.User) (linkanything.User, error) {
	rows, err := r.dbpool.Query(context.Background(), "select id from users where name=$1", user.Name)
	if err != nil {
		return linkanything.User{}, err
	}
	if rows.Err() != nil {
		return linkanything.User{}, rows.Err()
	}

	if rows.Next() {
		return linkanything.User{}, fmt.Errorf("user exists: %s", user.Name)
	}

	hash, err := HashPassword(user.Password)
	if err != nil {
		return linkanything.User{}, err
	}
	_, err = r.dbpool.Exec(context.Background(), "insert into users(name,password) values($1, $2)", user.Name, hash)
	if err != nil {
		return linkanything.User{}, err
	}
	return user, nil
}

func (r PostgresUserRepository) GetHash(user string) (string, error) {
	var dbHash string
	err := r.dbpool.QueryRow(context.Background(), "SELECT password FROM users WHERE name=$1", user).Scan(&dbHash)
	return dbHash, err
}

func (r PostgresUserRepository) DeleteAll() error {
	_, err := r.dbpool.Exec(context.Background(), "delete from users")
	return err
}

func (r PostgresUserRepository) ByName(name string) (linkanything.User, error) {
	rows, err := r.dbpool.Query(context.Background(), "select id, name, password from users where name=$1", name)
	if err != nil {
		log.Error(err)
		return linkanything.User{}, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		log.Error(err)
		return linkanything.User{}, err
	}

	if rows.Next() {
		var u linkanything.User
		err := rows.Scan(&u.Id, &u.Name, &u.Password)
		if err != nil {
			return linkanything.User{}, err
		}
		return u, nil
	}

	return linkanything.User{}, fmt.Errorf("does not user exists: %s", name)
}
