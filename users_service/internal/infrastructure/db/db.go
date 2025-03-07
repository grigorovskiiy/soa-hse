package db

import (
	"auth/users_service/internal/infrastructure"
	"auth/users_service/internal/infrastructure/repository"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/uptrace/bun"
	"os"
)

func InitDb() *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@postgres:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*repository.DbUser)(nil)).
		Exec(context.Background())

	if err != nil {
		infrastructure.Logger.Error(err.Error())
	}

	return db
}
