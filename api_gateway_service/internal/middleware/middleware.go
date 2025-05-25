package middleware

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/logger"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type Claims struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	UserID   int    `json:"user_id"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("secret-key")

func JWTVerify(r *http.Request) error {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		logger.Logger.Error("token is empty")
		return errors.New("token is empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Logger.Error("token sign error")
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		logger.Logger.Error("token is invalid")
		return errors.New("token is invalid")
	}

	r.Header.Set("Login", claims.Login)
	r.Header.Set("Password", claims.Password)
	r.Header.Set("UserID", strconv.Itoa(claims.UserID))

	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := JWTVerify(r)
		if err != nil {
			logger.Logger.Error("jwt verify erorr", "error", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func MethodMiddleware(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ProxyMiddleware(targetHost, targetPort string) func(http.Handler) http.Handler {
	targetURL := fmt.Sprintf("%s%s", targetHost, targetPort)
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   targetURL,
	})

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Logger.Error("proxy error", "error", err.Error())
		w.WriteHeader(http.StatusBadGateway)
	}

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = targetURL
		req.Host = targetURL
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Logger.Info("proxying request",
				"method", r.Method,
				"path", r.URL.Path,
				"target", targetURL)
			proxy.ServeHTTP(w, r)
		})
	}
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logger.Logger.With(
			"path", r.URL.Path,
			"method", r.Method,
		)
		logger.Info("request started")
		defer logger.Info("request finished")
		next.ServeHTTP(w, r)
	})
}
