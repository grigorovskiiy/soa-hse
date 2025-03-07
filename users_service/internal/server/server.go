package server

import (
	"auth/users_service/internal/api"
	"context"
	"errors"
	"go.uber.org/fx"
	"net/http"
	"os"
)

func NewServer(s *api.UsersService) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", http.HandlerFunc(s.Register))
	mux.HandleFunc("/login", http.HandlerFunc(s.Login))
	mux.HandleFunc("/get_user_info", http.HandlerFunc(s.GetUserInfo))
	mux.HandleFunc("/update_user_info", http.HandlerFunc(s.UpdateUserInfo))

	return &http.Server{
		Addr:    os.Getenv("USERS_SERVICE_PORT"),
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
