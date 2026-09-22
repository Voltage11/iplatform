package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Voltage11/iplatform/internal/config"
	"github.com/Voltage11/iplatform/internal/db"
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

	// 7. Сервисы
	userService := service.NewUserService(userRepo, database, cfg.HashPreffix)

	logger.Info("Запуск сервера на порту", "port", cfg.Server.Port)

	// 6. Запуск HTTP-сервера и ожидание ctx.Done()

	return nil
}
