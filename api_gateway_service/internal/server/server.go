package server

import (
	"auth/api_gateway_service/internal/api"
	"auth/api_gateway_service/internal/application"
	"context"
	"errors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
	"net/http"
	"os"
)

func NewServer(a *application.GatewayApp) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/register", http.HandlerFunc(a.Register))
	mux.Handle("/login", http.HandlerFunc(a.Login))
	mux.Handle("/get_user_info", http.HandlerFunc(a.GetUserInfo))
	mux.Handle("/update_user_info", http.HandlerFunc(a.UpdateUserInfo))
	mux.Handle("/create_post", http.HandlerFunc(a.CreatePost))
	mux.Handle("/delete_post", http.HandlerFunc(a.DeletePost))
	mux.Handle("/update_post", http.HandlerFunc(a.UpdatePost))
	mux.Handle("/get_post", http.HandlerFunc(a.GetPost))
	mux.Handle("/post_list", http.HandlerFunc(a.GetPostList))
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("swagger/swagger/doc.json")))

	return &http.Server{
		Addr:    os.Getenv("API_GATEWAY_PORT"),
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
