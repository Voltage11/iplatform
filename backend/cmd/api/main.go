package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Voltage11/iplatform/internal/appmiddleware"
	"github.com/Voltage11/iplatform/internal/config"
	"github.com/Voltage11/iplatform/internal/db"
	"github.com/Voltage11/iplatform/internal/handlers"
	"github.com/Voltage11/iplatform/internal/repo"
	"github.com/Voltage11/iplatform/internal/service"
	"github.com/Voltage11/iplatform/pkg/applog"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
}

func run() error {
	// 1. Конфиг
	cfg, err := config.New()
	if err != nil {
		return err
	}

	// 2. Уровень логгера
	lvl, err := applog.ParseLevel(cfg.Log.Level)
	if err != nil {
		return err
	}

	// 3. Логгер
	logger := applog.New(lvl)
	logger.SetAsDefault()

	// 4. Сигнальный контекст для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 5. БД
	database, err := db.New(ctx, &cfg.Database)
	if err != nil {
		return err
	}
	logger.Info("Соединение с бд успешно")
	defer database.Close()

	// 6. Репозитории
	userRepo := repo.NewUserRepo(database.Pool())
	sessionRepo := repo.NewSessionRepo(database.Pool())

	// 7. Сервисы
	jwtService := service.NewJWTService(service.ConfigJWT{
		SecretKey:  cfg.Jwt.Secret,
		AccessTTL:  cfg.Jwt.AccessTTL,
		RefreshTTL: cfg.Jwt.RefreshTTL,
	})
	userService := service.NewUserService(userRepo, database, cfg.Pepper)
	authService := service.NewAuthService(userService, sessionRepo, jwtService)

	logger.Info("Запуск сервера на порту", "port", cfg.Server.Port)

	// 8. Запуск HTTP-сервера и ожидание ctx.Done()
	r := chi.NewRouter()
	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Стандартные middleware chi
	r.Use(middleware.Timeout(cfg.Server.RequestTimeout))
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(cfg.Server.ReadTimeout))

	// Извлечение пользователя (глобально для всех маршрутов)
	authMW := appmiddleware.NewAuthMiddleware(userService, jwtService)
	r.Use(authMW.ExtractUser)

	// Handlers регистрация
	handlers.NewAuthHandler(authService).Register(r, authMW)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	ln, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	logger.Info("Сервер запущен", "port", cfg.Server.Port)

	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("serve", "err", err)
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
