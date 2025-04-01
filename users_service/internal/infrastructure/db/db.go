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
	dsn := fmt.Sprintf("postgres://%s:%s@users-postgres%s/%s?sslmode=disable",
		os.Getenv("USERS_POSTGRES_USER"), os.Getenv("USERS_POSTGRES_PASSWORD"), os.Getenv("USERS_POSTGRES_PORT"), os.Getenv("USERS_POSTGRES_DB"))

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
