package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/fx"
	"os"
)

func InitDb(lc fx.Lifecycle) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@posts-postgres%s/%s?sslmode=disable",
		os.Getenv("POSTS_POSTGRES_USER"), os.Getenv("POSTS_POSTGRES_PASSWORD"), os.Getenv("POSTS_POSTGRES_PORT"), os.Getenv("POSTS_POSTGRES_DB"))

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		infrastructure.Logger.Error("open database error", "error", err.Error())
		return nil
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*models.DbPost)(nil)).
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
