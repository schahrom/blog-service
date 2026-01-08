package main

import (
	"blog-service/internal/config"
	authHanlers "blog-service/internal/handlers/auth"
	noteHandlers "blog-service/internal/handlers/note"
	userHandlers "blog-service/internal/handlers/user"
	mvLogger "blog-service/internal/http-server/middleware"
	"blog-service/internal/lib"
	"blog-service/internal/lib/hasher"
	"blog-service/internal/repository"
	"blog-service/internal/service"
	"blog-service/internal/storage/postgresql"
	"blog-service/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetUpLogger(cfg.Env)

	db, err := postgresql.New(cfg.StorageUrl)
	if err != nil {
		log.Error("failed to connect to database", logger.Err(err))
		os.Exit(1)
	}

	_ = db

	router := chi.NewRouter()

	shaHasher := hasher.NewSHA256Hasher(cfg.Security.HasherSalt)
	notesRepo := repository.NewNotesRepoImpl(db)
	usersRepo := repository.NewUsersRepositoryImpl(db, shaHasher)

	ttl := time.Millisecond * time.Duration(cfg.TokenTTL)
	tokenManager := lib.NewTokenManagerImpl(cfg.Security.SigningKey, ttl)
	authService := service.NewAuthServiceImpl(usersRepo, tokenManager, shaHasher)

	//middleware

	authMiddleware := mvLogger.NewAuthMiddleware(tokenManager)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(mvLogger.New(log))
	router.Use(middleware.Recoverer)

	notesRouter := noteHandlers.NewNotesRoutes(*log, notesRepo)
	userRouter := userHandlers.NewUsersRouter(usersRepo, log)
	authRouter := authHanlers.NewAuthRouter(authService, log)

	//config routes
	config.InitRoutes(router, notesRouter, userRouter, authRouter, authMiddleware)

	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", logger.Err(err))
	}

	log.Info("starting blog-service ", slog.String("env", cfg.Env))
	log.Debug("debug logging enabled")
}
