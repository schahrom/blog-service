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

type NotesRepo interface {
	Save(note models.NoteDto, userId int) (int, error)
	UpdateNote(note models.NoteDto) error
	GetById(id int) (models.NoteDto, error)
	GetAllNotes(userId int, page models.PageDto) ([]models.NoteDto, error)
	DeleteById(id int) error
}

type NotesRouter struct {
	repo NotesRepo
	log  *slog.Logger
}

func NewNotesRoutes(log slog.Logger, repo NotesRepo) *NotesRouter {
	return &NotesRouter{repo: repo, log: &log}
}

func (router *NotesRouter) SaveNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.SaveNote"

		requestPathUserId, done := compareTokenAndPathUserId(w, r, router, op)
		if done {
			return
		}

		router.log = router.log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req models.NoteDto
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
		intUserId, _ := strconv.Atoi(requestPathUserId)

		_, err = router.repo.Save(req, intUserId)

		if err != nil {
			router.log.Error("failed to save note", logger.Err(err))
			render.JSON(w, r, response.Error("failed to save note"))
			return
		}

		router.log.Info("note saved", slog.Any("note", req))
		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, response.OK())
	}
}

func (router *NotesRouter) UpdateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.UpdateNote"

		router.log = router.log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		_, isNotSimilar := compareTokenAndPathUserId(w, r, router, op)
		if isNotSimilar {
			return
		}
		var req models.NoteDto
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

			render.JSON(w, r, response.ValidationError(validationErrors))

			return
		}

		err = router.repo.UpdateNote(req)

		if err != nil {
			router.log.Error("failed to update note by id: "+strconv.Itoa(req.NoteId), logger.Err(err))
		}

		router.log.Info("note updated", slog.Any("note", req))
		render.JSON(w, r, response.OK())
	}
}

func (router *NotesRouter) DeleteNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.UpdateNote"

		_, isNotSimilar := compareTokenAndPathUserId(w, r, router, op)
		if isNotSimilar {
			return
		}

		router.log = router.log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var noteIdParam int
		noteIdParam, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.JSON(w, r, response.Error("failed to decode note id"))
		}

		err = router.repo.DeleteById(noteIdParam)

		if err != nil {
			router.log.Error("failed to delete note", logger.Err(err))
			render.JSON(w, r, response.Error("failed to delete note"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		router.log.Info("note deleted", slog.Any("note", noteIdParam))
		render.JSON(w, r, response.OK())
	}
}

func (router *NotesRouter) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.GetNotes"
		router.log = router.log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		_, isNotSimilar := compareTokenAndPathUserId(w, r, router, op)
		if isNotSimilar {
			return
		}

		var noteIdParam int
		noteIdParam, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.JSON(w, r, response.Error("failed to decode note id"))
		}

		note, err := router.repo.GetById(noteIdParam)
		if err != nil {
			router.log.Error("failed to fetch note", logger.Err(err))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error("failed to fetch note"))
			return
		}

		router.log.Info("note retrieved", slog.Any("note", note))
		render.JSON(w, r, response.OkWithData(note))
	}
}

func (router *NotesRouter) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.note.GetNotes"

		router.log = router.log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		_, isNotSimilar := compareTokenAndPathUserId(w, r, router, op)
		if isNotSimilar {
			return
		}

		var userIdParam int
		userIdParam, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			render.JSON(w, r, response.Error("failed to decode note id"))
			return
		}
		query := r.URL.Query()
		queryParamDto, err := models.NewPageDto(query.Get("limit"), query.Get("offset"), query.Get("sort"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to extract query param"))
			return
		}

		notes, err := router.repo.GetAllNotes(userIdParam, queryParamDto)
		if err != nil {
			router.log.Error("failed to fetch notes", logger.Err(err))
			render.JSON(w, r, response.Error("failed to fetch notes"))
			return
		}

		router.log.Info("notes retrieved", slog.Any("notes", notes))
		render.JSON(w, r, response.OkWithData(notes))
	}
}

func compareTokenAndPathUserId(w http.ResponseWriter, r *http.Request, router *NotesRouter, op string) (string, bool) {
	tokenUserId := r.Header.Get("user_id")
	requestPathUserId := chi.URLParam(r, "id")
	if tokenUserId != requestPathUserId {
		w.WriteHeader(http.StatusForbidden)
		router.log.Error(op, "invalid user", "user_id", requestPathUserId)
		render.JSON(w, r, response.Error("invalid user"))
		return "", true
	}
	return requestPathUserId, false
}
