package application

import (
	"auth/api_gateway_service/internal/infrastructure"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type GatewayApp struct{}

func NewGatewayApp() *GatewayApp {
	return &GatewayApp{}
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

func JWTVerify(r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		infrastructure.Logger.Error("token is empty")
		return errors.New("token is empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			infrastructure.Logger.Error("token sign error")
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		infrastructure.Logger.Error("token is invalid")
		return errors.New("token is invalid")
	}

	r.Header.Set("Login", claims.Login)
	r.Header.Set("Password", claims.Password)

	return nil
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
func (a *GatewayApp) Register(w http.ResponseWriter, r *http.Request) {
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
func (s *GatewayApp) Login(w http.ResponseWriter, r *http.Request) {
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
// @Router       /update_user [put]
func (s *GatewayApp) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		_ = json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
	}

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
// @Router       /get_user [get]
func (s *GatewayApp) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		_ = json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
	}

	proxy := CreateProxy()
	if proxy == nil {
		infrastructure.Logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxy.ServeHTTP(w, r)
}

func (s *GatewayApp) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		_ = json.NewEncoder(w).Encode(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
	}
	conn, err := grpc.NewClient("localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()))

}
