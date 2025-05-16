package application

import (
	"context"
	"encoding/json"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/models"
	"io"
	"net/http"
	"time"
)

type UsersService interface {
	Register(*models.RegisterRequest) (int, error)
	Login(*models.GetLoginRequest) (string, error)
	UpdateUserInfo(*models.UserUpdateRequest, string) error
	GetUserInfo(string) (*models.DbUser, error)
}

type EventsService interface {
	SendEvent(context.Context, string, any) error
}
type UsersApp struct {
	UsersService  UsersService
	EventsService EventsService
	cfg           *config.Config
}

func NewUsersApp(uS UsersService, eS EventsService, cfg *config.Config) *UsersApp {
	return &UsersApp{UsersService: uS, EventsService: eS, cfg: cfg}
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
	if err = json.Unmarshal(d, &req); err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	userID, err := a.UsersService.Register(&req)
	if err != nil {
		logger.Error("service register error", "error", err.Error())
		writeRes(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = a.EventsService.SendEvent(context.Background(), a.cfg.ClientsTopic, models.ClientUpdate{UserId: userID, Time: time.Now()}); err != nil {
		logger.Error("register send event error", "error", err.Error())
		writeRes(w, http.StatusInternalServerError, err.Error())
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

	token, err := a.UsersService.Login(&req)
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

	err = a.UsersService.UpdateUserInfo(&req, login)
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

	user, err := a.UsersService.GetUserInfo(login)
	if err != nil {
		logger.Error("service get user info error", "error", err.Error())
		writeRes(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeRes(w, http.StatusOK, user)
}
