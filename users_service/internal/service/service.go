package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/infrastructure/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var secretKey = []byte("secret-key")

type Repository interface {
	Register(*models.DbUser) error
	Login(*models.GetLoginRequest) error
	UpdateUserInfo(*models.DbUser, string) error
	GetUserInfo(string) (*models.DbUser, error)
	GetUserID(string) (int, error)
}

type UService struct {
	repository Repository
}

func NewUService(repository Repository) *UService {
	return &UService{
		repository: repository,
	}
}

func (a *UService) Register(req *models.RegisterRequest) error {
	userInfo := models.DbUser{
		Email:     req.Email,
		Login:     req.Login,
		Password:  req.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.repository.Register(&userInfo); err != nil {
		logger.Logger.Error("db register user info error", "error", err.Error())
		return err
	}

	return nil
}

func (a *UService) Login(req *models.GetLoginRequest) (string, error) {
	err := a.repository.Login(req)
	if err != nil {
		logger.Logger.Error("db login error", "error", err.Error())
		return "", err
	}

	userID, err := a.repository.GetUserID(req.Login)
	if err != nil {
		logger.Logger.Error("db get user id error", "error", err.Error())
		return "", err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":    req.Login,
		"password": req.Password,
		"user_id":  userID,
	})

	token, err := claims.SignedString(secretKey)
	if err != nil {
		logger.Logger.Error("signing jwt error", "error", err.Error())
		return "", nil
	}

	return token, nil
}

func (a *UService) UpdateUserInfo(req *models.UserUpdateRequest, login string) error {
	userInfo := models.DbUser{
		Email:     req.Email,
		UpdatedAt: time.Now(),
		Name:      req.Name,
		Surname:   req.Surname,
	}

	if err := a.repository.UpdateUserInfo(&userInfo, login); err != nil {
		logger.Logger.Error("db update user info error", "error", err.Error())
		return err
	}

	return nil
}

func (a *UService) GetUserInfo(login string) (*models.DbUser, error) {
	userInfo, err := a.repository.GetUserInfo(login)
	if err != nil {
		logger.Logger.Error("db get user info error", "error", err.Error())
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return userInfo, nil
}
