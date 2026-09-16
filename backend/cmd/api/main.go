package main

import (
	"log"

	"github.com/Voltage11/iplatform/internal/config"
	"github.com/Voltage11/iplatform/pkg/applog"
)

func main() {
	// Загружаем конфиг приложения, если есть ошибка, то падаем с ошибкой
	config, err := config.New()
	if err != nil {
		log.Fatalf("error get config: %v", err)
	}
	log.Println("Конфигурация загружена")

	// Логгер
	logger := applog.New(applog.GetLoggerLevelFromInt(config.Log.Level))
}
