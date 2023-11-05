package router

import (
	"chatie.com/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func InitRouter(userHandler *handlers.UserHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:     nil,
		AllowOriginFunc:    nil,
		AllowedMethods:     nil,
		AllowedHeaders:     nil,
		ExposedHeaders:     nil,
		AllowCredentials:   false,
		MaxAge:             0,
		OptionsPassthrough: false,
		Debug:              false,
	}))

	r.Post("/signup", userHandler.CreateUser)
	r.Post("/login", userHandler.Login)
	r.Get("/logout", userHandler.Logout)

	return r
}
