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

	voterRepo := repository.NewVoterRepository()
	voterSrv := service.NewVoterService(db, voterRepo)
	voterHandler := handler.NewVoterHandler(voterSrv)

	authSrv := service.NewAuthService(oauthConfig)
	authHandler := handler.NewAuthHandler(authSrv)

	router := routes.SetupRoutes(voterHandler, authHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("server running on http://localhost:%s\n", port)

	log.Fatal(http.ListenAndServe(":"+port, middleware.Logger(router)))

}
