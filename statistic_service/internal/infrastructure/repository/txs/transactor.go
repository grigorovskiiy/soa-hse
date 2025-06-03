package txs

import (
	"context"
	"database/sql"
	"errors"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/logger"
)

type Querier interface {
	Exec(sql string, arguments ...any) (sql.Result, error)
	Query(sql string, args ...any) (*sql.Rows, error)
	QueryRow(sql string, args ...any) *sql.Row
}

type txKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetQuerier(ctx context.Context, defaultQuerier Querier) Querier {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}

	return defaultQuerier
}

type TxBeginner struct {
	db *sql.DB
}

func NewTxBeginner(db *sql.DB) *TxBeginner {
	return &TxBeginner{db: db}
}

func (r *TxBeginner) WithTransaction(ctx context.Context, txFunc func(ctx context.Context) error) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Logger.Error("error starting tx", "error", err.Error())
		return err
	}

	defer func() {
		if err != nil {
			logger.Logger.Error(err.Error())
			err = errors.Join(err, tx.Rollback())
		}
	}()

	err = txFunc(injectTx(ctx, tx))
	if err != nil {
		logger.Logger.Error(err.Error())
		return err
	}

	return tx.Commit()
}

func (r *TxBeginner) WithTransactionWithValue(ctx context.Context, txFunc func(ctx context.Context) (any, error)) (val any, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Logger.Error("error starting tx", "error", err.Error())
		return nil, err
	}

	defer func() {
		if err != nil {
			logger.Logger.Error(err.Error())
			err = errors.Join(err, tx.Rollback())
		}
	}()

	val, err = txFunc(injectTx(ctx, tx))
	if err != nil {
		logger.Logger.Error(err.Error())
		return nil, err
	}

	return val, tx.Commit()
}
