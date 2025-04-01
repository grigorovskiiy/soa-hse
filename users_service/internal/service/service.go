package service

import (
	"auth/users_service/internal/infrastructure/models"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

type Repository interface {
	Register(userInfo *models.UserGetRegisterLogin) error
	Login(userInfo *models.UserGetRegisterLogin) error
	UpdateUserInfo(userInfo *models.UserUpdate) error
	GetUserInfo(userinfo *models.UserGetRegisterLogin) (*models.UserUpdate, error)
}

type UService struct {
	repository Repository
}

func (a *UService) Register(user *models.UserGetRegisterLogin) error {
	return a.repository.Register(user)
}

func (a *UService) Login(user *models.UserGetRegisterLogin) (string, error) {
	err := a.repository.Login(user)
	if err != nil {
		return "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":    user.Login,
		"password": user.Password,
		"email":    user.Email,
	})

	token, err := claims.SignedString(secretKey)
	if err != nil {
		return "", nil
	}

	return token, nil
}

func (a *UService) UpdateUserInfo(user *models.UserUpdate) error {
	return a.repository.UpdateUserInfo(user)
}

func (a *UService) GetUserInfo(user *models.UserGetRegisterLogin) (*models.UserUpdate, error) {
	return a.repository.GetUserInfo(user)
}
