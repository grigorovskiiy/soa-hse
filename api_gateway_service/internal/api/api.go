package api

import (
	"auth/api_gateway_service/internal/application"
	"auth/api_gateway_service/internal/infrastructure"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type GatewayService struct {
	a *application.GatewsayApp
}

func NewApiGatewayService(a *application.GatewsayApp) *GatewayService {
	return &GatewayService{
		a: a,
	}
}

type Claims struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret-key")

func CreateProxy() *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("users-service%s", os.Getenv("USERS_SERVICE_PORT")),
	})

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("users-service%s", os.Getenv("USERS_SERVICE_PORT"))
		req.Host = fmt.Sprintf("users-service%s", os.Getenv("USERS_SERVICE_PORT"))
	}

	return proxy
}

// Register godoc
// @Summary      Регистрация
// @Description  Зарегистрироваться в сервисе
// @Tags         Auth
// @Accept		 json
// @Produce      json
// @Param 		 user body repository.UserGetRegisterLogin true "Зарегистрировать пользователя"
// @Success      200
// @Failure		 400 {string} string
// @Router       /register [post]
func (s *GatewayService) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	proxy := CreateProxy()
	if proxy == nil {
		infrastructure.Logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy.ServeHTTP(w, r)
}

// Login godoc
// @Summary      Войти
// @Description  Войти в систему
// @Tags         Auth
// @Accept		 json
// @Produce      json
// @Param 		 user body repository.UserGetRegisterLogin true "Войти в систему"
// @Success      200  {string} string
// @Failure 	 400 {string} string
// @Router       /login [post]
func (s *GatewayService) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	proxy := CreateProxy()
	if proxy == nil {
		infrastructure.Logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy.ServeHTTP(w, r)
}

// UpdateUserInfo godoc
// @Summary      Обновить пользователя
// @Description  Обновить данные о пользователе
// @Tags         User
// @Accept		 json
// @Security BearerAuth
// @Produce      json
// @Param 		 user body repository.UserUpdate true "Обновить пользователя"
// @Success      200
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500
// @Router       /update_user_info [put]
func (s *GatewayService) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		infrastructure.Logger.Error("token is empty")
		_ = json.NewEncoder(w).Encode("token is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			infrastructure.Logger.Error("token sign error")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		infrastructure.Logger.Error("token is invalid")
		_ = json.NewEncoder(w).Encode("token is invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.Header.Set("Login", claims.Login)
	r.Header.Set("Password", claims.Password)

	proxy := CreateProxy()
	if proxy == nil {
		infrastructure.Logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy.ServeHTTP(w, r)
}

// GetUserInfo godoc
// @Summary      Получить пользователя
// @Description  Получить пользователя
// @Tags         User
// @Accept		 application/x-www-form-urlencoded
// @Security BearerAuth
// @Produce      json
// @Success      200  {object} repository.UserUpdate
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500
// @Router       /get_user_info [get]
func (s *GatewayService) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		infrastructure.Logger.Error("token is empty")
		_ = json.NewEncoder(w).Encode("token is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			infrastructure.Logger.Error("token sign error")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		infrastructure.Logger.Error("token is invalid")
		_ = json.NewEncoder(w).Encode("token is invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	r.Header.Set("Login", claims.Login)
	r.Header.Set("Password", claims.Password)

	proxy := CreateProxy()
	if proxy == nil {
		infrastructure.Logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy.ServeHTTP(w, r)
}
