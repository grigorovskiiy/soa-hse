package application

import (
	"auth/users_service/internal/infrastructure/repository"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

type UsersApp struct {
	R repository.Repository
}

func NewUsersApp(r repository.Repository) *UsersApp {
	return &UsersApp{r}
}

func (a *UsersApp) Register(user *repository.UserGetRegisterLogin) error {
	return a.R.Register(user)
}

func (a *UsersApp) Login(user *repository.UserGetRegisterLogin) (string, error) {
	err := a.R.Login(user)
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

func (a *UsersApp) UpdateUserInfo(user *repository.UserUpdate) error {
	return a.R.UpdateUserInfo(user)
}

func (a *UsersApp) GetUserInfo(user *repository.UserGetRegisterLogin) (*repository.UserUpdate, error) {
	return a.R.GetUserInfo(user)
}
