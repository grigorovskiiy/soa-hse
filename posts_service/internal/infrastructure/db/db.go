package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/fx"
)

func CreateTable(db *bun.DB, model interface{}) error {
	_, err := db.NewCreateTable().
		IfNotExists().
		Model(model).
		Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func CreateTables(db *bun.DB) error {
	if err := CreateTable(db, (*models.DbPost)(nil)); err != nil {
		logger.Logger.Error("create posts table error", "error", err.Error())
		return err
	}
	if err := CreateTable(db, (*models.DbComment)(nil)); err != nil {
		logger.Logger.Error("create comments table error", "error", err.Error())
		return err
	}
	if err := CreateTable(db, (*models.DbLike)(nil)); err != nil {
		logger.Logger.Error("create likes table error", "error", err.Error())
		return err
	}
	if err := CreateTable(db, (*models.DbView)(nil)); err != nil {
		logger.Logger.Error("create views table error", "error", err.Error())
		return err
	}

	return nil
}

func InitDb(lc fx.Lifecycle, cfg *config.Config) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s%s/%s?sslmode=disable",
		cfg.PostsPostgresUser, cfg.PostsPostgresPassword, cfg.PostsPostgresHost, cfg.PostsPostgresPort, cfg.PostsPostgresDb)

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Logger.Error("open database error", "error", err.Error())
		return nil
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	err = CreateTables(db)
	if err != nil {
		logger.Logger.Error("create tables", "error", err.Error())
		return nil
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}
