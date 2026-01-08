package handlers

import (
	"blog-service/internal/http-server/api/response"
	"blog-service/internal/models"
	"blog-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"strconv"
)

type UsersRouter struct {
	userRepo UserRepository
	log      *slog.Logger
}

func NewUsersRouter(userRepo UserRepository, log *slog.Logger) *UsersRouter {
	return &UsersRouter{userRepo: userRepo, log: log}
}

type UserRepository interface {
	Create(user *models.User) (models.User, error)
	GetById(id int) (models.User, error)
	GetByCredentials(login string, hashedPassword string) (models.User, error)
}

func (router *UsersRouter) SaveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.SaveUser"

		router.log = router.log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req models.User
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			router.log.Error("failed to decode body", logger.Err(err))
			render.JSON(w, r, response.Error("failed to decode body"))
			return
		}

		router.log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validationErrors := err.(validator.ValidationErrors)

			router.log.Error("failed to validate body", logger.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validationErrors))
			return
		}

		user, err := router.userRepo.Create(&req)
		if err != nil {
			router.log.Error("failed to create user", logger.Err(err))
			render.JSON(w, r, response.Error("failed to create user"))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		router.log.Info("successfully created user with", slog.Any("user", user))
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, response.OkWithData(user))
	}
}

func (router *UsersRouter) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.user.GetById"

		router.log = router.log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userId, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			router.log.Error("failed to decode user id", logger.Err(err))
			render.JSON(w, r, response.Error("failed to decode user id"))
			return
		}
		user, err := router.userRepo.GetById(userId)
		if err != nil {
			router.log.Error("failed to get user", logger.Err(err))
			render.JSON(w, r, response.Error("failed to get user"))
			return
		}

		router.log.Info("successfully fetched user", slog.Any("user", user))
		render.JSON(w, r, response.OkWithData(user))
		return
	}
}
