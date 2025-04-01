package application

import (
	"auth/users_service/internal/infrastructure"
	"auth/users_service/internal/infrastructure/models"
	"encoding/json"
	"io"
	"net/http"
)

type UsersService interface {
	Register(user *models.UserGetRegisterLogin) error
	Login(user *models.UserGetRegisterLogin) (string, error)
	UpdateUserInfo(user *models.UserUpdate) error
	GetUserInfo(user *models.UserGetRegisterLogin) (*models.UserUpdate, error)
}
type UsersApp struct {
	Service UsersService
}

func NewUsersApp(service UsersService) *UsersApp {
	return &UsersApp{Service: service}
}

func (a *UsersApp) Register(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.UserGetRegisterLogin
	err = json.Unmarshal(d, &req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.Service.Register(&req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode("Пользователь зарегистрирован")
	w.WriteHeader(http.StatusOK)
}

func (a *UsersApp) Login(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.UserGetRegisterLogin
	err = json.Unmarshal(d, &req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.Service.Login(&req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusOK)
}

func (a *UsersApp) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("Login")
	password := r.Header.Get("Password")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.UserUpdate
	err = json.Unmarshal(d, &req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if password != req.Password || login != req.Login {
		infrastructure.Logger.Error("No privileges")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode("No privileges")
		return
	}

	err = a.Service.UpdateUserInfo(&req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *UsersApp) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("Login")
	password := r.Header.Get("Password")

	req := models.UserGetRegisterLogin{Login: login, Password: password}

	user, err := a.Service.GetUserInfo(&req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}
