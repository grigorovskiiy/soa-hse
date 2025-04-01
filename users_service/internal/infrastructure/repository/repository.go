package repository

import (
	"auth/users_service/internal/errors"
	"auth/users_service/internal/infrastructure/models"
	"context"
	"github.com/uptrace/bun"
	"time"
)

type UsersRepository struct {
	db *bun.DB
}

func NewUsersRepository(db *bun.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Register(userInfo *models.UserGetRegisterLogin) error {
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

	user := models.DbUser{
		Email:     userInfo.Email,
		Login:     userInfo.Login,
		Password:  userInfo.Password,
		CreatedAt: time.Now(),
	}

	_, err = r.db.NewInsert().Model(&user).Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) Login(userInfo *models.UserGetRegisterLogin) error {
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

func (r *UsersRepository) UpdateUserInfo(userInfo *models.UserUpdate) error {
	user := models.DbUser{
		Email:     userInfo.Email,
		Login:     userInfo.Login,
		Password:  userInfo.Password,
		UpdatedAt: time.Now(),
		Name:      userInfo.Name,
		Surname:   userInfo.Surname,
	}

	_, err := r.db.NewUpdate().Model(&user).Where("login = ?", user.Login).OmitZero().Exec(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) GetUserInfo(userInfo *models.UserGetRegisterLogin) (*models.UserUpdate, error) {
	var user models.UserUpdate
	err := r.db.NewSelect().
		Model(&user).
		Where("login = ?", userInfo.Login).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return &user, nil
}
