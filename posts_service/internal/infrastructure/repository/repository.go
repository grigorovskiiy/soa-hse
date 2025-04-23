package repository

import (
	"context"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/errors"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure"
	"github.com/grigorovskiiy/soa-hse/posts_service/internal/infrastructure/models"
	"github.com/uptrace/bun"
)

type PRepository struct {
	db *bun.DB
}

func NewPRepository(db *bun.DB) *PRepository {
	return &PRepository{db: db}
}

func (r *PRepository) GetPost(postId int32, userId int32) (*models.DbPost, error) {
	exists, err := r.db.NewSelect().
		Model((*models.DbPost)(nil)).
		Where("id = ? and user_id = ?", postId, userId).
		Exists(context.Background())
	if err != nil {
		infrastructure.Logger.Error("get post db error", err.Error())
		return nil, err
	}
	if !exists {
		infrastructure.Logger.Info(errors.PostNotFoundError{}.Error())
		return nil, errors.PostNotFoundError{}
	}

	var post models.DbPost
	err = r.db.NewSelect().Model(&post).Where("id = ? and user_id = ?", postId, userId).Scan(context.Background())
	if err != nil {
		infrastructure.Logger.Error("get post db error", err.Error())
		return nil, err
	}

	return &post, nil
}

func (r *PRepository) GetPostList(page int32, limit int32, userId int32) ([]*models.DbPost, error) {
	var posts []*models.DbPost

	offset := (page - 1) * limit

	query := r.db.NewSelect().
		Model(&posts).
		Where("security_flag = ?", false).
		OrderExpr("created_at DESC").
		Limit(int(limit)).
		Offset(int(offset))

	query = query.WhereOr("user_id = ?", userId)
	err := query.Scan(context.Background(), &posts)
	if err != nil {
		infrastructure.Logger.Error("scan get post list error", "error", err.Error())
		return nil, err
	}

	return posts, nil
}

func (r *PRepository) UpdatePost(post *models.DbPost) error {
	exists, err := r.db.NewSelect().
		Model((*models.DbPost)(nil)).
		Where("id = ? and user_id = ?", post.Id, post.UserId).
		Exists(context.Background())
	if err != nil {
		infrastructure.Logger.Error("exists update post db error", "error", err.Error())
		return err
	}
	if !exists {
		infrastructure.Logger.Info(errors.PostNotFoundError{}.Error())
		return errors.PostNotFoundError{}
	}

	_, err = r.db.NewUpdate().Model(post).Where("id = ? and user_id = ?", post.Id, post.UserId).OmitZero().Exec(context.Background())
	if err != nil {
		infrastructure.Logger.Error("execing update post db error", "error", err.Error())
		return err
	}

	return nil
}

func (r *PRepository) DeletePost(postId int32, userId int32) error {
	exists, err := r.db.NewSelect().
		Model((*models.DbPost)(nil)).
		Where("id = ? and user_id = ?", postId, userId).
		Exists(context.Background())
	if err != nil {
		infrastructure.Logger.Error("exists delete post db error", "error", err.Error())
		return err
	}
	if !exists {
		infrastructure.Logger.Info(errors.PostNotFoundError{}.Error())
		return errors.PostNotFoundError{}
	}

	_, err = r.db.NewDelete().Model(&models.DbPost{}).Where("id = ? and user_id = ?", postId, userId).Exec(context.Background())
	if err != nil {
		infrastructure.Logger.Error("execing delete post db error", "error", err.Error())
		return err
	}
	return nil
}

func (r *PRepository) CreatePost(post *models.DbPost) error {
	_, err := r.db.NewInsert().Model(post).Exec(context.Background())
	if err != nil {
		infrastructure.Logger.Error("execing create post db error", "error", err.Error())
		return err
	}

	return nil
}
