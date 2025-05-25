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
// @Param 		 post_id query int true "ID поста"
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
// @Param        page query int true "Номер страницы"
// @Param        page_size query int true "Количество элементов на странице"
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
	_, err = a.GRPCClients.PostsServiceClient.CreatePost(ctx, req.ToPostsProto())
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
// @Param 		 post_id body models.PostID true "ID поста"
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
	_, err = a.GRPCClients.PostsServiceClient.DeletePost(ctx, req.ToPostsProto())
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
	_, err = a.GRPCClients.PostsServiceClient.UpdatePost(ctx, req.ToPostsProto())
	if err != nil {
		logger.Error("error grpc request UpdatePost", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "Post is updated")
}

// PostComment godoc
// @Summary      Добавить комментарий к посту
// @Description  Добавить комментарий к посту
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param 		 comment_info body models.PostCommentRequest true "Информация о комментарии"
// @Success      200  {string} string
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /post_comment [post]
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
	_, err = a.GRPCClients.PostsServiceClient.PostComment(ctx, req.ToPostsProto())
	if err != nil {
		logger.Error("error  grpc request PostComment", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "Comment is posted")
}

// PostLike godoc
// @Summary      Добавить лайк к посту
// @Description  Добавить лайк к посту
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param 		 post_id body models.PostID true "ID поста"
// @Success      200  {string} string
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /post_like [post]
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
	_, err = a.GRPCClients.PostsServiceClient.PostLike(ctx, req.ToPostsProto())
	if err != nil {
		logger.Error("error grpc request PostLike", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "Like is posted")
}

// PostView godoc
// @Summary      Добавить просмотр к посту
// @Description  Добавить просмотр к посту
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param 		 post_id body models.PostID true "ID поста"
// @Success      200  {string} string
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /post_view [post]
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
	_, err = a.GRPCClients.PostsServiceClient.PostView(ctx, req.ToPostsProto())
	if err != nil {
		logger.Error("error grpc request PostView", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, "View is posted")
}

// GetCommentList godoc
// @Summary      Получить пагинированный список комментариев
// @Description  Получить пагинированный список комментариев
// @Tags         Post
// @Security BearerAuth
// @Produce      json
// @Param        page query int true "Номер страницы"
// @Param        page_size query int true "Количество элементов на странице"
// @Success      200  {object} models.GetCommentListResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_comment_list [get]
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
		logger.Error("error grpc request GetCommentList", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoListCommentResponse(res))
}

// GetViewsCount godoc
// @Summary      Получить количество просмотров по посту
// @Description  Получить количество просмотров по посту
// @Tags         Statistic
// @Produce      json
// @Param 		 post_id query int true "ID поста"
// @Success      200  {object} models.CountResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_views_count [get]
func (a *GatewayApp) GetViewsCount(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	postIdStr := query.Get("post_id")

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		logger.Error("postId is empty")
		writeRes(w, http.StatusBadRequest, "postId is empty")
	}

	req := models.PostID{PostID: postId}

	res, err := a.GRPCClients.StatisticServiceClient.GetViewsCount(r.Context(), req.ToStatisticProto())
	if err != nil {
		logger.Error("error grpc request GetViewsCount", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoCountResponse(res))
}

// GetLikesCount godoc
// @Summary      Получить количество лайков по посту
// @Description  Получить количество лайков по посту
// @Tags         Statistic
// @Produce      json
// @Param 		 post_id query int true "ID поста"
// @Success      200  {object} models.CountResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_likes_count [get]
func (a *GatewayApp) GetLikesCount(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	postIdStr := query.Get("post_id")

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		logger.Error("postId is empty")
		writeRes(w, http.StatusBadRequest, "postId is empty")
	}

	req := models.PostID{PostID: postId}

	res, err := a.GRPCClients.StatisticServiceClient.GetLikesCount(r.Context(), req.ToStatisticProto())
	if err != nil {
		logger.Error("error grpc request GetLikesCount", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoCountResponse(res))
}

// GetCommentsCount godoc
// @Summary      Получить количество комментариев по посту
// @Description  Получить количество комментариев по посту
// @Tags         Statistic
// @Produce      json
// @Param 		 post_id query int true "ID поста"
// @Success      200  {object} models.CountResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_comments_count [get]
func (a *GatewayApp) GetCommentsCount(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	postIdStr := query.Get("post_id")

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		logger.Error("postId is empty")
		writeRes(w, http.StatusBadRequest, "postId is empty")
	}

	req := models.PostID{PostID: postId}

	res, err := a.GRPCClients.StatisticServiceClient.GetCommentsCount(r.Context(), req.ToStatisticProto())
	if err != nil {
		logger.Error("error grpc request GetCommentsCount", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoCountResponse(res))
}

// GetCommentsDynamic godoc
// @Summary      Получить динамику комментариев по посту
// @Description  Получить динамику комментариев по посту
// @Tags         Statistic
// @Produce      json
// @Param 		 post_id query int true "ID поста"
// @Success      200  {object} models.DynamicListResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_comments_dynamic [get]
func (a *GatewayApp) GetCommentsDynamic(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	postIdStr := query.Get("post_id")

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		logger.Error("postId is empty")
		writeRes(w, http.StatusBadRequest, "postId is empty")
	}

	req := models.PostID{PostID: postId}

	res, err := a.GRPCClients.StatisticServiceClient.GetCommentsDynamic(r.Context(), req.ToStatisticProto())
	if err != nil {
		logger.Error("error grpc request GetCommentsDynamic", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoDynamuicListResponse(res))
}

// GetLikesDynamic godoc
// @Summary      Получить динамику лайков по посту
// @Description  Получить динамику лайков по посту
// @Tags         Statistic
// @Produce      json
// @Param 		 post_id query int true "ID поста"
// @Success      200  {object} models.DynamicListResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_likes_dynamic [get]
func (a *GatewayApp) GetLikesDynamic(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	postIdStr := query.Get("post_id")

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		logger.Error("postId is empty")
		writeRes(w, http.StatusBadRequest, "postId is empty")
	}

	req := models.PostID{PostID: postId}

	res, err := a.GRPCClients.StatisticServiceClient.GetLikesDynamic(r.Context(), req.ToStatisticProto())
	if err != nil {
		logger.Error("error grpc request GetLikesDynamic", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoDynamuicListResponse(res))
}

// GetViewsDynamic godoc
// @Summary      Получить динамику просмотров по посту
// @Description  Получить динамику просмотров по посту
// @Tags         Statistic
// @Produce      json
// @Param 		 post_id query int true "ID поста"
// @Success      200  {object} models.DynamicListResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_views_dynamic [get]
func (a *GatewayApp) GetViewsDynamic(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	postIdStr := query.Get("post_id")

	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		logger.Error("postId is empty")
		writeRes(w, http.StatusBadRequest, "postId is empty")
	}

	req := models.PostID{PostID: postId}

	res, err := a.GRPCClients.StatisticServiceClient.GetViewsDynamic(r.Context(), req.ToStatisticProto())
	if err != nil {
		logger.Error("error grpc request GetViewsDynamic", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoDynamuicListResponse(res))
}

// GetTopTenPosts godoc
// @Summary      Получить топ 10 постов по параметру
// @Description  Получить топ 10 постов по параметру
// @Tags         Statistic
// @Produce      json
// @Param 		 top_parameter query string true "Параметер топа"
// @Success      200  {object} models.TopTenResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_top_ten_posts [get]
func (a *GatewayApp) GetTopTenPosts(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	par := query.Get("top_parameter")

	req := models.TopParameter{Parameter: par}

	res, err := a.GRPCClients.StatisticServiceClient.GetTopTenPosts(r.Context(), req.ToStatisticProto())
	if err != nil {
		logger.Error("error grpc request GetTopTenPosts", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoTopTenPostsResponse(res))
}

// GetTopTenUsers godoc
// @Summary      Получить топ 10 пользователей по параметру
// @Description  Получить топ 10 пользователей по параметру
// @Tags         Statistic
// @Produce      json
// @Param 		 top_parameter query string true "Параметер топа"
// @Success      200  {object} models.TopTenResponse
// @Failure 	 400 {string} string
// @Failure 	 401  {string} string
// @Failure 	 500 {string} string
// @Router       /get_top_ten_users [get]
func (a *GatewayApp) GetTopTenUsers(w http.ResponseWriter, r *http.Request) {
	logger := logger.Logger.With("path", r.URL.Path, "method", r.Method)

	query := r.URL.Query()
	par := query.Get("top_parameter")

	req := models.TopParameter{Parameter: par}

	res, err := a.GRPCClients.StatisticServiceClient.GetTopTenUsers(r.Context(), req.ToStatisticProto())
	if err != nil {
		logger.Error("error grpc request GetTopTenUsers", status.Convert(err).Message())
		writeRes(w, http.StatusInternalServerError, status.Convert(err).Message())
		return
	}

	writeRes(w, http.StatusOK, models.FromProtoTopTenUsersResponse(res))
}
