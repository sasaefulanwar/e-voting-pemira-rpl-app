package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"pemira-rpl/internal/config"
	"pemira-rpl/internal/handler"
	"pemira-rpl/internal/middleware"
	"pemira-rpl/internal/repository"
	"pemira-rpl/internal/routes"
	"pemira-rpl/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := config.InitDB()
	defer db.Close()

	oauthConfig := config.InitOAuthConfig()

	voterRepo := repository.NewVoterRepository(db)
	voterSrv := service.NewVoterService(db, voterRepo)
	voterHandler := handler.NewVoterHandler(voterSrv)

	authSrv := service.NewAuthService(oauthConfig, voterRepo)
	authHandler := handler.NewAuthHandler(authSrv)

	voteRepo := repository.NewVoteRepository()
	electionRepo := repository.NewElectionRepository(db)
	electionService := service.NewElectionService(db, electionRepo, voterRepo)

	adminHandler :=
		handler.NewAdminHandler(
			electionService,
		)
	auditRepo := repository.NewAuditRepository()
	voteSrv := service.NewVoteService(
		db,
		voterRepo,
		voteRepo,
		electionRepo,
		auditRepo,
	)
	voteHandler := handler.NewVoteHandler(voteSrv)

	candidateRepo := repository.NewCandidateRepository(db)
	candidateService := service.NewCandidateService(candidateRepo, electionRepo)
	candidateHandler := handler.NewCandidateHandler(candidateService)
	electionHandler := handler.NewElectionHandler(electionService)

	disputeRepo :=
		repository.NewDisputeRepository(
			db,
		)

	disputeService :=
		service.NewDisputeService(
			disputeRepo,
			voterRepo,
		)

	disputeHandler :=
		handler.NewDisputeHandler(
			disputeService,
		)

	router := routes.SetupRoutes(
		voterHandler,
		authHandler,
		voteHandler,
		candidateHandler,
		adminHandler,
		electionHandler,
		disputeHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("server running on http://localhost:%s\n", port)
	handlerWithCors :=
		middleware.CORS(
			middleware.Logger(
				router,
			),
		)

	log.Fatal(
		http.ListenAndServe(
			":"+port,
			handlerWithCors,
		),
	)

}
