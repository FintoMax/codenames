package api

import (
	"codenames/auth"
	"codenames/config"
	"context"
	"log/slog"
	"net/http"
)

type API struct {
	cfg         *config.Config
	authService *auth.AuthService

	server *http.Server
}

func NewAPI(cfg *config.Config, authService *auth.AuthService) *API {
	api := &API{
		cfg:         cfg,
		authService: authService,
	}

	mux := api.mux()

	api.server = &http.Server{
		Handler: mux,
		Addr:    "localhost:8080",
	}

	return api
}

func (a *API) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		slog.Info("shutting down server")
		a.server.Shutdown(context.Background())
	}()

	slog.Info("starting server")
	return a.server.ListenAndServe()
}

func (a *API) mux() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
