package routes

import (
	"net/http"
	"pemira-rpl/internal/handler"
	"pemira-rpl/internal/middleware"
)

func SetupRoutes(voterHandler *handler.VoterHandler, authHandler *handler.AuthHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/health", handler.HealthCheck)

	mux.HandleFunc("/api/v1/auth/google/login", authHandler.GoogleLogin)
	mux.HandleFunc("/api/v1/auth/google/callback", authHandler.GoogleCallback)

	mux.Handle("/api/v1/voter/bind", middleware.AuthMiddleware(http.HandlerFunc(voterHandler.BindNIM)))

	return mux
}
