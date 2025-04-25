package application

import (
	"context"
	"encoding/json"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/clients"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/logger"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/infrastructure/models"
	pb "github.com/grigorovskiiy/soa-hse/protos"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"net/http"
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

func writeRes(w http.ResponseWriter, code int, val any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(val)
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
	logger.Logger.Info("request proxied", "path", "/register")
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
	logger.Logger.Info("request proxied", "path", "/login")
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
	logger.Logger.Info("request proxied", "path", "/update_user_info")
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
	logger.Logger.Info("request proxied", "path", "/get_user_info")
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
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	postIDStr := query.Get("post_id")

	if postIDStr == "" {
		writeRes(w, http.StatusBadRequest, "post_id is empty")
		return
	}

	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		writeRes(w, http.StatusBadRequest, "post_id is invalid")
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)

	res, err := a.GRPCClients.PostsServiceClient.GetPost(ctx, &pb.PostID{PostId: int32(postID)})
	if err != nil {
		logger.Error("grpc request GetPost error", "error", status.Convert(err))
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoPostResponse(res))
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
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	pageStr := query.Get("page")
	pageSizeStr := query.Get("page_size")

	if pageStr == "" || pageSizeStr == "" {
		writeRes(w, http.StatusBadRequest, "page or page size is empty")
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	res, err := a.GRPCClients.PostsServiceClient.GetPostList(ctx, &pb.PaginatedListRequest{Page: int32(page), PageSize: int32(pageSize)})
	if err != nil {
		logger.Error("rpc request GetPostList", "error", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoListPostResponse(res))
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
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	var req models.CreatePostRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.CreatePost(ctx, req.ToProto())
	if err != nil {
		logger.Error("grpc request CreatePost error", "error", status.Convert(err))
		writeRes(w, http.StatusInternalServerError, status.Convert(err))
		return
	}

	writeRes(w, http.StatusOK, "Post is created")
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
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	var req models.PostID
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
		return
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.DeletePost(ctx, req.ToProto())
	if err != nil {
		logger.Error("grpc request DeletePost error", "error", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "Post is deleted")
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
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}
	var req models.UpdatePostRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.UpdatePost(ctx, req.ToProto())
	if err != nil {
		logger.Error("error api grpc request UpdatePost", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "Post is updated")
}

func (a *GatewayApp) PostComment(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	var req models.PostCommentRequest
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.PostComment(ctx, req.ToProto())
	if err != nil {
		logger.Error("error api grpc request PostComment", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "Comment is posted")
}

func (a *GatewayApp) PostLike(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	var req models.PostID
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.PostLike(ctx, req.ToProto())
	if err != nil {
		logger.Error("error api grpc request PostLike", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "Like is posted")
}

func (a *GatewayApp) PostView(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	d, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("read body error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	var req models.PostID
	err = json.Unmarshal(d, &req)
	if err != nil {
		logger.Error("unmarshal error", "error", err.Error())
		writeRes(w, http.StatusBadRequest, nil)
		return
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	_, err = a.GRPCClients.PostsServiceClient.PostView(ctx, req.ToProto())
	if err != nil {
		logger.Error("error api grpc request PostView", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "View is posted")
}

func (a *GatewayApp) GetCommentList(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	pageStr := query.Get("page")
	pageSizeStr := query.Get("page_size")

	if pageStr == "" || pageSizeStr == "" {
		writeRes(w, http.StatusBadRequest, "page or page size is empty")
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	userID := r.Header.Get("UserID")
	if userID == "" {
		logger.Error("user_id is empty")
		writeRes(w, http.StatusBadRequest, "user_id is empty")
		return
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "user_id", userID)
	res, err := a.GRPCClients.PostsServiceClient.GetCommentList(ctx, &pb.PaginatedListRequest{Page: int32(page), PageSize: int32(pageSize)})
	if err != nil {
		logger.Error("error api grpc request GetCommentList", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoListCommentResponse(res))
}
