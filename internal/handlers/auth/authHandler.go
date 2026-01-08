package handlers

import (
	"blog-service/internal/http-server/api/response"
	"blog-service/internal/models"
	"blog-service/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type AuthRouter struct {
	service AuthService
	log     *slog.Logger
}

func NewAuthRouter(service AuthService, log *slog.Logger) *AuthRouter {
	return &AuthRouter{service: service, log: log}
}

type AuthService interface {
	SignIn(requestDto models.SignInDto) (string, error)
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (router *AuthRouter) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.signIn"

		router.log = router.log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var request models.SignInDto
		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			router.log.Error("failed to decode body", logger.Err(err))
			render.JSON(w, r, response.Error("failed to decode body"))
			return
		}
		router.log.Info("sign in request decoded")

		token, err := router.service.SignIn(request)
		if err != nil {
			router.log.Error("failed to sign in", logger.Err(err))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error("failed to sign in"))
			return
		}

		router.log.Info("sign in request signed successfully")
		render.JSON(w, r, TokenResponse{Token: token})
	}
}
