package application

import (
	"encoding/json"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/logger"
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

func writeRes(w http.ResponseWriter, code int, val any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(val)
}

func (a *UsersApp) Register(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)
	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	var req models.RegisterRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	err = a.Service.Register(&req)
	if err != nil {
		logger.Error("service register error", "error", err.Error())
		writeRes(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeRes(w, http.StatusOK, "user is registered")
}

func (a *UsersApp) Login(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	var req models.GetLoginRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	token, err := a.Service.Login(&req)
	if err != nil {
		logger.Error("service login error", "error", err.Error())
		writeRes(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeRes(w, http.StatusOK, token)
}

func (a *UsersApp) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	var req models.UserUpdateRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}
	login := r.Header.Get("Login")

	err = a.Service.UpdateUserInfo(&req, login)
	if err != nil {
		logger.Error("service update error", "error", err.Error())
		writeRes(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeRes(w, http.StatusOK, "user is updated")
}

func (a *UsersApp) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)
	login := r.Header.Get("Login")

	user, err := a.Service.GetUserInfo(login)
	if err != nil {
		logger.Error("service get user info error", "error", err.Error())
		writeRes(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeRes(w, http.StatusOK, user)
}
