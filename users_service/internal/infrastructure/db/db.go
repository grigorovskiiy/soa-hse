package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/fx"

	"github.com/uptrace/bun"
	"os"
)

func InitDb(lc fx.Lifecycle) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@users-postgres%s/%s?sslmode=disable",
		os.Getenv("USERS_POSTGRES_USER"), os.Getenv("USERS_POSTGRES_PASSWORD"), os.Getenv("USERS_POSTGRES_PORT"), os.Getenv("USERS_POSTGRES_DB"))

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		infrastructure.Logger.Error("open database error", "error", err.Error())
		return nil
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*models.DbUser)(nil)).
		Exec(context.Background())

	if err != nil {
		infrastructure.Logger.Error("create table error", "error", err.Error())
		return nil
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}
