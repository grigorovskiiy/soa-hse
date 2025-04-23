package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/clients"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/models"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

type GatewayApp struct {
	GRPCClients *clients.GRPCClients
}

func NewGatewayApp(GRPCClients *clients.GRPCClients) *GatewayApp {
	return &GatewayApp{
		GRPCClients: GRPCClients,
	}
}

type Claims struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	UserID   int    `json:"user_id"`
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
	r.Header.Set("UserID", strconv.Itoa(claims.UserID))

	return nil
}

// Register godoc
// @Summary      Регистрация
// @Description  Зарегистрироваться в сервисе
// @Tags         Auth
// @Accept		 json
// @Produce      json
// @Param 		 user body models.RegisterRequest true "Зарегистрировать пользователя"
// @Success      200
// @Failure		 400 {string} string
// @Router       /register [post]
func (a *GatewayApp) Register(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/register")
	logger.Info("request started")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	proxy := CreateProxy()
	if proxy == nil {
		logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.Info("request finished")
	proxy.ServeHTTP(w, r)
}

// Login godoc
// @Summary      Войти
// @Description  Войти в систему
// @Tags         Auth
// @Accept		 json
// @Produce      json
// @Param 		 user body models.GetLoginRequest true "Войти в систему"
// @Success      200  {string} string
// @Failure 	 400 {string} string
// @Router       /login [post]
func (a *GatewayApp) Login(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/login")
	logger.Info("request started")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	proxy := CreateProxy()
	if proxy == nil {
		logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("request finished")
	proxy.ServeHTTP(w, r)
}

// UpdateUserInfo godoc
// @Summary      Обновить пользователя
// @Description  Обновить данные о пользователе
// @Tags         User
// @Accept		 json
// @Security BearerAuth
// @Produce      json
// @Param 		 user body models.UserUpdateRequest true "Обновить пользователя"
// @Success      200
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500
// @Router       /update_user_info [put]
func (a *GatewayApp) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/update_user_info")
	logger.Info("request started")
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	proxy := CreateProxy()
	if proxy == nil {
		logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("request finished")
	proxy.ServeHTTP(w, r)
}

// GetUserInfo godoc
// @Summary      Получить пользователя
// @Description  Получить пользователя
// @Tags         User
// @Accept		 application/x-www-form-urlencoded
// @Security BearerAuth
// @Produce      json
// @Success      200  {object} models.GetLoginRequest
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500
// @Router       /get_user_info [get]
func (a *GatewayApp) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/get_user_info")
	logger.Info("request started")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		logger.Error("jwt verify error", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	proxy := CreateProxy()
	if proxy == nil {
		logger.Error("proxy error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("request finished")
	proxy.ServeHTTP(w, r)
}

// GetPost godoc
// @Summary      Получить пост
// @Description  Получить пост
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param 		 post_id body models.GetDeletePostRequest true "ID поста"
// @Success      200  {object} models.GetPostResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_post [get]
func (a *GatewayApp) GetPost(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/get_post")
	logger.Info("request started")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		logger.Error("jwt verify error", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.GetDeletePostRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)

	res, err := a.GRPCClients.PostsServiceClient.GetPost(ctx, req.ToProto())
	if err != nil {
		st := status.Convert(err)
		logger.Error("grpc request GetPost error", "error", st.Message())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(st.Message())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.FromProtoPostResponse(res))
	logger.Info("request finished")
}

// GetPostList godoc
// @Summary      Получить пагинированный список постов
// @Description  Получить пагинированный список постов
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param 		 list_params body models.GetPostListRequest true "Параметры списка"
// @Success      200  {object} models.GetPostListResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_post_list [get]
func (a *GatewayApp) GetPostList(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/get_post_list")
	logger.Info("request started")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		logger.Error("jwt verify error", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.GetPostListRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	res, err := a.GRPCClients.PostsServiceClient.GetPostList(ctx, req.ToProto())
	if err != nil {
		st := status.Convert(err)
		logger.Error("rpc request GetPostList", "error", st.Message())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(st.Message())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(models.FromProtoListPostResponse(res))
	logger.Info("request finished")
}

// CreatePost godoc
// @Summary      Создать пост
// @Description  Создать пост
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param 		 post_info body models.CreatePostRequest true "Информация о посте"
// @Success      200  {string} string
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /create_post [post]
func (a *GatewayApp) CreatePost(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/create_post")
	logger.Info("request started")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		infrastructure.Logger.Error(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.CreatePostRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.CreatePost(ctx, req.ToProto())
	if err != nil {
		st := status.Convert(err)
		logger.Error("grpc request CreatePost error", "error", st.Message())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(st.Message())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Post is created")
	logger.Info("request finished")
}

// DeletePost godoc
// @Summary      Удалить пост
// @Description  Удалить пост
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param 		 post_id body models.GetDeletePostRequest true "ID поста"
// @Success      200  {string} string
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /delete_post [delete]
func (a *GatewayApp) DeletePost(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/delete_post")
	logger.Info("request started")
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		logger.Error("jwt verify error", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req models.GetDeletePostRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("user_id is empty")
		return
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.DeletePost(ctx, req.ToProto())
	if err != nil {
		st := status.Convert(err)
		logger.Error("grpc request DeletePost error", "error", st.Message())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(st.Message())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Post is deleted")
	logger.Info("request finished")
}

// UpdatePost godoc
// @Summary      Обновить пост
// @Description  Обновить пост
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param 		 post_info body models.UpdatePostRequest true "Информация о посте"
// @Success      200  {string} string
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /update_post [put]
func (a *GatewayApp) UpdatePost(w http.ResponseWriter, r *http.Request) {
	logger := infrastructure.Logger.With("path", "/update_post")
	logger.Info("request started")
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := JWTVerify(r)
	if err != nil {
		logger.Error("jwt verify error", "error", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req models.UpdatePostRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode("user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.UpdatePost(ctx, req.ToProto())
	if err != nil {
		st := status.Convert(err)
		logger.Error("error api grpc request UpdatePost", st.Message())
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(st.Message())
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode("Post is updated")
	logger.Info("request finished")
}
