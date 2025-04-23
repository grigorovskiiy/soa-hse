package repository

import (
	"context"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/errors"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/models"
	"github.com/uptrace/bun"
)

type UsersRepository struct {
	db *bun.DB
}

func NewUsersRepository(db *bun.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Register(userInfo *models.DbUser) error {
	exists, err := r.db.NewSelect().
		Model((*models.DbUser)(nil)).
		Where("login = ?", userInfo.Login).
		Exists(context.Background())
	if err != nil {
		return err
	}
	if exists {
		return errors.AlreadyRegisteredError{}
	}

	_, err = r.db.NewInsert().Model(userInfo).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) Login(userInfo *models.GetLoginRequest) error {
	exists, err := r.db.NewSelect().
		Model((*models.DbUser)(nil)).
		Where("login = ? and password = ?", userInfo.Login, userInfo.Password).
		Exists(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		return errors.LoginError{}
	}

	return nil
}

func (r *UsersRepository) UpdateUserInfo(userInfo *models.DbUser, login string) error {
	_, err := r.db.NewUpdate().Model(userInfo).Where("login = ?", login).OmitZero().Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) GetUserInfo(userLogin string) (*models.DbUser, error) {
	var user models.DbUser
	err := r.db.NewSelect().
		Model(&user).
		Where("login = ?", userLogin).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) GetUserID(login string) (int, error) {
	var user models.DbUser
	err := r.db.NewSelect().
		Model(&user).
		Where("login = ?", login).
		Scan(context.Background())

	if err != nil {
		return 0, err
	}

	return user.Id, nil
}
