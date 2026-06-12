package routes

import (
	"net/http"
	"pemira-rpl/internal/handler"
	"pemira-rpl/internal/middleware"
)

func SetupRoutes(voterHandler *handler.VoterHandler, authHandler *handler.AuthHandler, voteHandler *handler.VoteHandler, candidateHandler *handler.CandidateHandler, adminHandler *handler.AdminHandler, electionHandler *handler.ElectionHandler) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(
		http.Dir("./images"),
	)

	mux.Handle(
		"/images/",
		http.StripPrefix(
			"/images/",
			fs,
		),
	)

	mux.HandleFunc("/api/v1/health", handler.HealthCheck)

	mux.HandleFunc("/api/v1/auth/google/login", authHandler.GoogleLogin)

	mux.HandleFunc("/api/v1/auth/google/callback", authHandler.GoogleCallback)

	mux.Handle("/api/v1/voter/bind", middleware.AuthMiddleware(http.HandlerFunc(voterHandler.BindNIM)))

	mux.Handle(
		"/api/v1/vote",
		middleware.AuthMiddleware(
			http.HandlerFunc(
				voteHandler.CastVote,
			),
		),
	)

	mux.Handle(
		"/api/v1/candidates",
		middleware.AuthMiddleware(
			http.HandlerFunc(
				candidateHandler.GetAll,
			),
		),
	)

	mux.Handle(
		"/api/v1/results",
		middleware.AuthMiddleware(
			http.HandlerFunc(candidateHandler.GetResults),
		),
	)

	mux.HandleFunc("/api/v1/election/status", electionHandler.GetStatus)
	mux.HandleFunc("/api/v1/admin/election/open", electionHandler.OpenElection)
	mux.HandleFunc("/api/v1/admin/election/close", electionHandler.CloseElection)

	mux.HandleFunc(
		"/api/v1/admin/statistics",
		adminHandler.GetStatistics,
	)

	return mux
}
