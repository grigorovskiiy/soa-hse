package db

import (
	"auth/posts_service/internal/infrastructure"
	"auth/posts_service/internal/infrastructure/models"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"os"
)

func InitDb() *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost%s/%s?sslmode=disable",
		os.Getenv("POSTS_POSTGRES_USER"), os.Getenv("POSTS_POSTGRES_PASSWORD"), os.Getenv("POSTS_POSTGRES_PORT"), os.Getenv("POSTS_POSTGRES_DB"))

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*models.DbPost)(nil)).
		Exec(context.Background())

	if err != nil {
		infrastructure.Logger.Error(err.Error())
	}

	return db
}
