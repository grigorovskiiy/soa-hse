package usersservice

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
	Register(*models.DbUser) (int, error)
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

func (a *UService) Register(req *models.RegisterRequest) (int, error) {
	userInfo := models.DbUser{
		Email:     req.Email,
		Login:     req.Login,
		Password:  req.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	id, err := a.repository.Register(&userInfo)
	if err != nil {
		logger.Logger.Error("register user info error", "error", err.Error())
		return 0, err
	}

	return id, nil
}

func (a *UService) Login(req *models.GetLoginRequest) (string, error) {
	if err := a.repository.Login(req); err != nil {
		logger.Logger.Error("login error", "error", err.Error())
		return "", err
	}

	userID, err := a.repository.GetUserID(req.Login)
	if err != nil {
		logger.Logger.Error("get user id error", "error", err.Error())
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
		logger.Logger.Error("update user info error", "error", err.Error())
		return err
	}

	return nil
}

func (a *UService) GetUserInfo(login string) (*models.DbUser, error) {
	userInfo, err := a.repository.GetUserInfo(login)
	if err != nil {
		logger.Logger.Error("get user info error", "error", err.Error())
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return userInfo, nil
}
