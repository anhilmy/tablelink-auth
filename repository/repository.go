package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/anhilmy/tablelink-auth/entity"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	GetUser(ctx context.Context, username string) (*entity.User, error)
}

type repository struct {
	db        *sql.DB
	rdsClient *redis.Client
}

func NewRepository(db *sql.DB, rds *redis.Client) Repository {
	return &repository{db, rds}
}

func (r *repository) GetUser(ctx context.Context, username string) (*entity.User, error) {
	res := entity.User{}

	rows, err := r.db.Query("SELECT * FROM users WHERE username = $1 LIMIT 1", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&res.Id, &res.Name, &res.Email, &res.Password, &res.CreatedAt, &res.LastAccess)
		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		log.Fatal(err)
	}

	return &res, nil
}
