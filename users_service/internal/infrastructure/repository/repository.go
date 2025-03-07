package repository

import (
	"auth/users_service/internal/errors"
	"context"
	"github.com/uptrace/bun"
	"time"
)

type Repository interface {
	Register(userInfo *UserGetRegisterLogin) error
	Login(userInfo *UserGetRegisterLogin) error
	UpdateUserInfo(userInfo *UserUpdate) error
	GetUserInfo(userinfo *UserGetRegisterLogin) (*UserUpdate, error)
}

type UsersRepository struct {
	db *bun.DB
}

func NewUsersRepository(db *bun.DB) Repository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Register(userInfo *UserGetRegisterLogin) error {
	exists, err := r.db.NewSelect().
		Model((*DbUser)(nil)).
		Where("login = ?", userInfo.Login).
		Exists(context.Background())
	if err != nil {
		return err
	}
	if exists {
		return errors.AlreadyRegisteredError{}
	}

	user := DbUser{
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

func (r *UsersRepository) Login(userInfo *UserGetRegisterLogin) error {
	exists, err := r.db.NewSelect().
		Model((*DbUser)(nil)).
		Where("login = ?", userInfo.Login).
		Exists(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		return errors.LoginError{}
	}

	return nil
}

func (r *UsersRepository) UpdateUserInfo(userInfo *UserUpdate) error {
	user := DbUser{
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

func (r *UsersRepository) GetUserInfo(userInfo *UserGetRegisterLogin) (*UserUpdate, error) {
	var user UserUpdate
	err := r.db.NewSelect().
		Model(&user).
		Where("login = ?", userInfo.Login).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return &user, nil
}
