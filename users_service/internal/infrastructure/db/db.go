package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/fx"

	"github.com/uptrace/bun"
)

func InitDb(lc fx.Lifecycle, cfg *config.Config) *bun.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s%s/%s?sslmode=disable",
		cfg.UsersPostgresUser, cfg.UsersPostgresPassword, cfg.UsersPostgresHost, cfg.UsersPostgresPort, cfg.UsersPostgresDb)

	sqldb, err := sql.Open("pgx", dsn)
	if err != nil {
		logger.Logger.Error("open database error", "error", err.Error())
		return nil
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	_, err = db.NewCreateTable().
		IfNotExists().
		Model((*models.DbUser)(nil)).
		Exec(context.Background())

	if err != nil {
		logger.Logger.Error("create table error", "error", err.Error())
		return nil
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close()
		},
	})

	return db
}
