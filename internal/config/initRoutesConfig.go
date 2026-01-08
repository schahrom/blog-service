package config

import (
	authHanlers "blog-service/internal/handlers/auth"
	noteHandlers "blog-service/internal/handlers/note"
	userHandlers "blog-service/internal/handlers/user"
	mvLogger "blog-service/internal/http-server/middleware"
	"github.com/go-chi/chi/v5"
)

func InitRoutes(router *chi.Mux, notesRouter *noteHandlers.NotesRouter, userRouter *userHandlers.UsersRouter, authRouter *authHanlers.AuthRouter, middleware *mvLogger.AuthMiddleware) {
	router.Route("/users/{id}", func(r chi.Router) {
		r.Use(middleware.CheckAuthorizationToken)

		r.Post("/notes", notesRouter.SaveNote())
		r.Get("/notes", notesRouter.GetNotes())
		r.Get("/notes/{note_id}", notesRouter.GetById())
		r.Put("/notes/{note_id}", notesRouter.UpdateNote())
		r.Delete("/notes/{note_id}", notesRouter.DeleteNote())
	})
	router.Get("/users/{id}", userRouter.GetById())
	router.Post("/users", userRouter.SaveUser())
	router.Post("/login", authRouter.SignIn())
}
