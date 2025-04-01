package repository

import (
	"auth/posts_service/internal/errors"
	"auth/posts_service/internal/infrastructure"
	"auth/posts_service/internal/infrastructure/models"
	"context"
	"github.com/uptrace/bun"
)

type PRepository struct {
	db *bun.DB
}

func NewPRepository(db *bun.DB) *PRepository {
	return &PRepository{db: db}
}

func (r *PRepository) GetPost(postId int32) (*models.DbPost, error) {
	exists, err := r.db.NewSelect().
		Model((*models.DbPost)(nil)).
		Where("id = ?", postId).
		Exists(context.Background())
	if err != nil {
		infrastructure.Logger.Error("get post db error", err.Error())
		return nil, err
	}
	if exists {
		infrastructure.Logger.Info(errors.PostNotFoundError{}.Error())
		return nil, errors.PostNotFoundError{}
	}

	var post models.DbPost
	err = r.db.NewSelect().Model(&post).Where("id = ?", postId).Scan(context.Background())
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
		return nil, err
	}

	return posts, nil
}

func (r *PRepository) UpdatePost(post *models.DbPost) error {
	exists, err := r.db.NewSelect().
		Model((*models.DbPost)(nil)).
		Where("id = ?", post.Id).
		Exists(context.Background())
	if err != nil {
		infrastructure.Logger.Error("update post db error", err.Error())
		return err
	}
	if exists {
		infrastructure.Logger.Info(errors.PostNotFoundError{}.Error())
		return errors.PostNotFoundError{}
	}

	_, err = r.db.NewUpdate().Model(&post).Where("id = ?", post.Id).OmitZero().Exec(context.Background())
	if err != nil {
		infrastructure.Logger.Error("update post db error", err.Error())
		return err
	}

	return nil
}

func (r *PRepository) DeletePost(postId int32) error {
	exists, err := r.db.NewSelect().
		Model((*models.DbPost)(nil)).
		Where("id = ?", postId).
		Exists(context.Background())
	if err != nil {
		infrastructure.Logger.Error("delete post db error", err.Error())
		return err
	}
	if exists {
		infrastructure.Logger.Info(errors.PostNotFoundError{}.Error())
		return errors.PostNotFoundError{}
	}

	_, err = r.db.NewDelete().Model(&models.DbPost{}).Where("id = ?", postId).Exec(context.Background())
	if err != nil {
		infrastructure.Logger.Error("delete post db error", err.Error())
		return err
	}
	return nil
}

func (r *PRepository) CreatePost(post *models.DbPost) error {
	_, err := r.db.NewInsert().Model(&post).Exec(context.Background())
	if err != nil {
		infrastructure.Logger.Error("create post db error", err.Error())
		return err
	}

	return nil
}
