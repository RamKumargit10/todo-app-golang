package routes

import (
	taskHandlers "todo-app/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetUpRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(taskHandlers.AuthMiddleware)

	// API routes
	r.Get("/tasks", taskHandlers.GetTasks)
	r.Post("/tasks", taskHandlers.AddTask)
	r.Delete("/tasks/{id}", taskHandlers.DeleteTask)
	r.Put("/tasks/{id}", taskHandlers.UpdateTask)

	return r
}
