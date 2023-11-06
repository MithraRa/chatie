package main

import (
	"net/http"

	"chatie.com/internal/domain"
	"chatie.com/internal/handlers"
	repository "chatie.com/internal/repository/postgres"
	"chatie.com/internal/services"
	pkgpostgres "chatie.com/pkg/postgres"
	"chatie.com/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	logger := log.New()
	logger.SetLevel(log.DebugLevel)

	dsn := "postgres://admin:admin@localhost:5442/chatie" +
		"?sslmode=disable"

	pool, err := pkgpostgres.NewPool(dsn, logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer pool.Close()

	repo := repository.NewRepository(pool, logger)
	userService := services.NewUserService(logger, repo)
	userHandler := handlers.NewUserHandler(logger, userService)

	hubService := services.NewHubService(&domain.Hub{})
	hubHandler := handlers.NewHubHandler(hubService)

	r := router.InitRouter(userHandler, hubHandler)

	http.ListenAndServe(":3000", r)
}
