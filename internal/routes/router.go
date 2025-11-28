package routes

import (
	"database/sql"
	"main/internal/controllers"
	"main/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userRepo := &models.UserRepository{DB: db}
	userController := &controllers.UserController{Repo: userRepo}

	r.Route("/v1", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", userController.Register)
			r.Post("/login", userController.Login)
		})
	})
	return r
}
