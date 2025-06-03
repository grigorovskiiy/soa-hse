package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/errors"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/models"
	"github.com/grigorovskiiy/soa-hse/statistic_service/internal/infrastructure/repository/txs"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetViewsCount(ctx context.Context, postID int) (int, error) {
	querier := txs.GetQuerier(ctx, r.db)
	var count int
	err := querier.QueryRow("SELECT COUNT(*) FROM views WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		logger.Logger.Error("query get views count db error", "error", err.Error())
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetCommentsCount(ctx context.Context, postID int) (int, error) {
	querier := txs.GetQuerier(ctx, r.db)
	var count int
	err := querier.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		logger.Logger.Error("query get comments count db error", "error", err.Error())
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetLikesCount(ctx context.Context, postID int) (int, error) {
	querier := txs.GetQuerier(ctx, r.db)
	var count int
	err := querier.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		logger.Logger.Error("query get likes count db error", "error", err.Error())
		return 0, err
	}

	return count, nil
}

func (r *Repository) GetViewsDynamic(ctx context.Context, postID int) ([]*models.Dynamic, error) {

	querier := txs.GetQuerier(ctx, r.db)
	query := `
		SELECT 
			toDate(time) as date,
			COUNT(*) as count
		FROM views
		WHERE post_id = ?
		GROUP BY date
		ORDER BY date
	`

	rows, err := querier.Query(query, postID)
	if err != nil {
		logger.Logger.Error("query get views dynamic db error", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var dynamics []*models.Dynamic
	for rows.Next() {
		var d models.Dynamic
		if err := rows.Scan(&d.Date, &d.Count); err != nil {
			logger.Logger.Error("scan rows get views dynamic db error", "error", err.Error())
			return nil, err
		}
		dynamics = append(dynamics, &d)
	}

	return dynamics, nil
}

func (r *Repository) GetCommentsDynamic(ctx context.Context, postID int) ([]*models.Dynamic, error) {
	querier := txs.GetQuerier(ctx, r.db)
	query := `
		SELECT 
			toDate(time) as date,
			COUNT(*) as count
		FROM comments
		WHERE post_id = ?
		GROUP BY date
		ORDER BY date
	`

	rows, err := querier.Query(query, postID)
	if err != nil {
		logger.Logger.Error("query get comments dynamic db error", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var dynamics []*models.Dynamic
	for rows.Next() {
		var d models.Dynamic
		if err := rows.Scan(&d.Date, &d.Count); err != nil {
			logger.Logger.Error("scan rows get comments dynamic db error", "error", err.Error())
			return nil, err
		}

		dynamics = append(dynamics, &d)
	}

	return dynamics, nil
}

func (r *Repository) GetLikesDynamic(ctx context.Context, postID int) ([]*models.Dynamic, error) {
	querier := txs.GetQuerier(ctx, r.db)
	query := `
		SELECT 
			toDate(time) as date,
			COUNT(*) as count
		FROM likes
		WHERE post_id = ?
		GROUP BY date
		ORDER BY date
	`

	rows, err := querier.Query(query, postID)
	if err != nil {
		logger.Logger.Error("query get likes dynamic db error", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var dynamics []*models.Dynamic
	for rows.Next() {
		var d models.Dynamic
		if err := rows.Scan(&d.Date, &d.Count); err != nil {
			logger.Logger.Error("scan rows get likes dynamic db error", "error", err.Error())
			return nil, err
		}
		dynamics = append(dynamics, &d)
	}

	return dynamics, nil
}

func (r *Repository) GetTopTenPosts(ctx context.Context, par string) ([]int, error) {
	querier := txs.GetQuerier(ctx, r.db)
	if par != "likes" && par != "comments" && par != "views" {
		logger.Logger.Error(errors.InvalidTopParameterError{}.Error(), "par", par)
		return nil, errors.InvalidTopParameterError{}
	}

	query := fmt.Sprintf(`
		SELECT post_id
		FROM %s
		GROUP BY post_id
		ORDER BY COUNT(*) DESC
		LIMIT 10
	`, par)

	rows, err := querier.Query(query)
	if err != nil {
		logger.Logger.Error("query get top ten posts db error", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var postIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.Logger.Error("scan rows get top ten posts db error", "error", err.Error())
			return nil, err
		}
		postIDs = append(postIDs, id)
	}

	return postIDs, nil
}

func (r *Repository) GetTopTenUsers(ctx context.Context, par string) ([]int, error) {
	querier := txs.GetQuerier(ctx, r.db)
	if par != "likes" && par != "comments" && par != "views" {
		logger.Logger.Error(errors.InvalidTopParameterError{}.Error(), "par", par)
		return nil, errors.InvalidTopParameterError{}
	}

	query := fmt.Sprintf(`
		SELECT user_id
		FROM %s
		GROUP BY user_id
		ORDER BY COUNT(*) DESC
		LIMIT 10
	`, par)

	rows, err := querier.Query(query)
	if err != nil {
		logger.Logger.Error("query get top ten users db error", "error", err.Error())
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.Logger.Error("scan rows get top ten users db error", "error", err.Error())
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}
