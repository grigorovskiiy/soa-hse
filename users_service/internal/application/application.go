package application

import (
	"encoding/json"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/models"
	"io"
	"net/http"
)

type UsersService interface {
	Register(*models.RegisterRequest) error
	Login(*models.GetLoginRequest) (string, error)
	UpdateUserInfo(*models.UserUpdateRequest, string) error
	GetUserInfo(string) (*models.DbUser, error)
}
type UsersApp struct {
	Service UsersService
}

func NewUsersApp(service UsersService) *UsersApp {
	return &UsersApp{Service: service}
}

func (a *UsersApp) Register(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/register")
	logger.Info("request started")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.RegisterRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.Service.Register(&req)
	if err != nil {
		logger.Error("service register error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("User is registered")
	logger.Info("request finished")
}

func (a *UsersApp) Login(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/login")
	logger.Info("request started")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.GetLoginRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := a.Service.Login(&req)
	if err != nil {
		logger.Error("service login error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(token)
	logger.Info("request finished")
}

func (a *UsersApp) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/updateUserInfo")
	logger.Info("request started")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.UserUpdateRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	login := r.Header.Get("Login")

	err = a.Service.UpdateUserInfo(&req, login)
	if err != nil {
		logger.Error("service update error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("User info is updated")
	logger.Info("request finished")
}

func (a *UsersApp) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/getUserInfo")
	logger.Info("request started")
	login := r.Header.Get("Login")

	user, err := a.Service.GetUserInfo(login)
	if err != nil {
		logger.Error("service get user info error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user)
	logger.Info("request finished")
}
