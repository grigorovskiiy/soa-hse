package server

import (
	"context"
	"errors"
	"github.com/grigorovskiiy/soa-hse/users_service/internal/application"
	"go.uber.org/fx"
	"net/http"
	"os"
)

func NewServer(app *application.UsersApp) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", http.HandlerFunc(app.Register))
	mux.HandleFunc("/login", http.HandlerFunc(app.Login))
	mux.HandleFunc("/get_user_info", http.HandlerFunc(app.GetUserInfo))
	mux.HandleFunc("/update_user_info", http.HandlerFunc(app.UpdateUserInfo))

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
