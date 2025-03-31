package server

import (
	"auth/api_gateway_service/internal/api"
	"context"
	"errors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/fx"
	"net/http"
	"os"
)

func NewServer(s *api.GatewayService) *http.Server {
	mux := http.NewServeMux()

	mux.Handle("/register", http.HandlerFunc(s.Register))
	mux.Handle("/login", http.HandlerFunc(s.Login))
	mux.Handle("/get_user_info", http.HandlerFunc(s.GetUserInfo))
	mux.Handle("/update_user_info", http.HandlerFunc(s.UpdateUserInfo))
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
