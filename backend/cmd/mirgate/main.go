package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/Voltage11/iplatform/internal/config"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func main() {
	// Флаг для команды (up или down)
	var command string
	flag.StringVar(&command, "cmd", "up", "migration command: up or down")
	var step int
	flag.IntVar(&step, "step", 0, "number of migrations to apply (for down)")
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	if err := runMigrations(cfg.Database.DSN(), command, step); err != nil {
		log.Fatalf("Миграция не удалась: %v", err)
	}
	log.Println("Миграции применены успешно")
}

func runMigrations(databaseURL string, command string, step int) error {
	m, err := migrate.New("file://migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("не удалось инициализировать мигратор: %w", err)
	}
	defer m.Close()

	switch command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("ошибка при выполнении миграций UP: %w", err)
		}
	case "down":
		if step > 0 {
			if err := m.Steps(-step); err != nil {
				return fmt.Errorf("ошибка при откате миграций на %d шагов: %w", step, err)
			}
		} else {
			if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("ошибка при выполнении миграций DOWN: %w", err)
			}
		}
	default:
		return fmt.Errorf("неизвестная команда: %s", command)
	}

	return nil
}
