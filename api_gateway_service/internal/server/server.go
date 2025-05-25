package server

import (
	"context"
	"errors"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/application"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/config"
	"github.com/grigorovskiiy/soa-hse/api_gateway_service/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
	"net/http"
)

func NewServer(a *application.GatewayApp, cfg *config.Config) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/register",
		middleware.ProxyMiddleware(cfg.UsersServiceHost, cfg.UsersServicePort)(
			middleware.LoggerMiddleware(
				middleware.MethodMiddleware(http.MethodPost, http.HandlerFunc(a.Register)),
			)))

	mux.Handle("/login",
		middleware.ProxyMiddleware(cfg.UsersServiceHost, cfg.UsersServicePort)(
			middleware.LoggerMiddleware(
				middleware.MethodMiddleware(http.MethodPost, http.HandlerFunc(a.Register)),
			)))

	mux.Handle("/get_user_info",
		middleware.ProxyMiddleware(cfg.UsersServiceHost, cfg.UsersServicePort)(
			middleware.LoggerMiddleware(
				middleware.MethodMiddleware(http.MethodPost, http.HandlerFunc(a.GetUserInfo)),
			)))

	mux.Handle("/update_user_info",
		middleware.ProxyMiddleware(cfg.UsersServiceHost, cfg.UsersServicePort)(
			middleware.LoggerMiddleware(
				middleware.MethodMiddleware(http.MethodPost, http.HandlerFunc(a.UpdateUserInfo)),
			)))

	mux.Handle("/create_post",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodPost,
				middleware.AuthMiddleware(http.HandlerFunc(a.CreatePost)))))

	mux.Handle("/delete_post",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodDelete,
				middleware.AuthMiddleware(http.HandlerFunc(a.DeletePost)))))

	mux.Handle("/update_post",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodPut,
				middleware.AuthMiddleware(http.HandlerFunc(a.UpdatePost)))))

	mux.Handle("/get_post",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet,
				middleware.AuthMiddleware(http.HandlerFunc(a.GetPost)))))

	mux.Handle("/get_post_list",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet,
				middleware.AuthMiddleware(http.HandlerFunc(a.GetPostList)))))

	mux.Handle("/post_comment",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodPost,
				middleware.AuthMiddleware(http.HandlerFunc(a.PostComment)))),
	)

	mux.Handle("/post_like",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodPost,
				middleware.AuthMiddleware(http.HandlerFunc(a.PostLike)))))

	mux.Handle("/post_view",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodPost,
				middleware.AuthMiddleware(http.HandlerFunc(a.PostView)))))

	mux.Handle("/get_comment_list",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet,
				middleware.AuthMiddleware(http.HandlerFunc(a.GetCommentList)))))

	mux.Handle("/get_comments_count",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet, http.HandlerFunc(a.GetCommentsCount))))

	mux.Handle("/get_likes_count",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet, http.HandlerFunc(a.GetLikesCount))))

	mux.Handle("/get_views_count",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet, http.HandlerFunc(a.GetViewsCount))))

	mux.Handle("/get_comments_dynamic",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet, http.HandlerFunc(a.GetCommentsDynamic))))

	mux.Handle("/get_likes_dynamic",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet, http.HandlerFunc(a.GetLikesDynamic))))

	mux.Handle("/get_views_dynamic",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet, http.HandlerFunc(a.GetViewsDynamic))))

	mux.Handle("/get_top_ten_posts",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet, http.HandlerFunc(a.GetTopTenPosts))))

	mux.Handle("/get_top_ten_users",
		middleware.LoggerMiddleware(
			middleware.MethodMiddleware(http.MethodGet, http.HandlerFunc(a.GetTopTenUsers))))

	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("swagger/swagger/doc.json")))

	return &http.Server{
		Addr:    cfg.GatewayServicePort,
		Handler: mux,
	}

}

func RunServer(lc fx.Lifecycle, server *http.Server) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	return nil
}
