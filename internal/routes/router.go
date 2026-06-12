package routes

import (
	"net/http"
	"pemira-rpl/internal/handler"
	"pemira-rpl/internal/middleware"
)

func SetupRoutes(voterHandler *handler.VoterHandler, authHandler *handler.AuthHandler, voteHandler *handler.VoteHandler, candidateHandler *handler.CandidateHandler, adminHandler *handler.AdminHandler, electionHandler *handler.ElectionHandler, disputeHandler *handler.DisputeHandler) *http.ServeMux {
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

	mux.HandleFunc(
		"/api/v1/auth/logout",
		authHandler.Logout,
	)

	mux.Handle(
		"/api/v1/auth/me",
		middleware.AuthMiddleware(
			http.HandlerFunc(
				authHandler.Me,
			),
		),
	)

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

	mux.Handle(
		"/api/v1/admin/election/open",
		middleware.AdminOnly(
			middleware.AuthMiddleware(
				http.HandlerFunc(
					electionHandler.OpenElection,
				),
			),
		),
	)

	mux.Handle(
		"/api/v1/admin/election/close",
		middleware.AdminOnly(
			middleware.AuthMiddleware(
				http.HandlerFunc(
					electionHandler.CloseElection,
				),
			),
		),
	)

	mux.Handle(
		"/api/v1/disputes",
		middleware.AuthMiddleware(
			http.HandlerFunc(
				disputeHandler.SubmitDispute,
			),
		),
	)

	mux.Handle(
		"/api/v1/admin/disputes",
		middleware.AdminOnly(
			middleware.AuthMiddleware(
				http.HandlerFunc(
					disputeHandler.GetAllDisputes,
				),
			),
		),
	)

	mux.Handle(
		"/api/v1/admin/disputes/approve",
		middleware.AdminOnly(
			middleware.AuthMiddleware(
				http.HandlerFunc(
					disputeHandler.ApproveDispute,
				),
			),
		),
	)

	mux.Handle(
		"/api/v1/admin/disputes/reject",
		middleware.AdminOnly(
			middleware.AuthMiddleware(
				http.HandlerFunc(
					disputeHandler.RejectDispute,
				),
			),
		),
	)

	mux.Handle(
		"/api/v1/admin/statistics",
		middleware.AdminOnly(
			middleware.AuthMiddleware(
				http.HandlerFunc(
					adminHandler.GetStatistics,
				),
			),
		),
	)

	return mux
}
