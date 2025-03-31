package api

import (
	"auth/users_service/internal/application"
	"auth/users_service/internal/infrastructure"
	"auth/users_service/internal/infrastructure/repository"
	"encoding/json"
	"io"
	"net/http"
)

type UsersService struct {
	A *application.UsersApp
}

func NewUsersService(a *application.UsersApp) *UsersService {
	return &UsersService{A: a}
}

func (s *UsersService) Register(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req repository.UserGetRegisterLogin
	err = json.Unmarshal(d, &req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.A.Register(&req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode("Пользователь зарегистрирован")
	w.WriteHeader(http.StatusOK)
}

func (s *UsersService) Login(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req repository.UserGetRegisterLogin
	err = json.Unmarshal(d, &req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.A.Login(&req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(token)
	w.WriteHeader(http.StatusOK)
}

func (s *UsersService) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("Login")
	password := r.Header.Get("Password")

	d, err := io.ReadAll(r.Body)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req repository.UserUpdate
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

	err = s.A.UpdateUserInfo(&req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *UsersService) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("Login")
	password := r.Header.Get("Password")

	req := repository.UserGetRegisterLogin{Login: login, Password: password}

	user, err := s.A.GetUserInfo(&req)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	_ = json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}
