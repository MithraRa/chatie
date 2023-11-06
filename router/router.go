package router

import (
	"chatie.com/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func InitRouter(userHandler *handlers.UserHandler, hubHandler *handlers.HubHandler) *chi.Mux {
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

	r.Post("/ws/createRoom", hubHandler.CreateRoom)
	r.Get("/ws/joinRoom/:roomId", hubHandler.JoinRoom)
	r.Get("/ws/getRooms", hubHandler.GetRooms)
	r.Get("/ws/getClients/:roomId", hubHandler.GetClients)

	return r
}
