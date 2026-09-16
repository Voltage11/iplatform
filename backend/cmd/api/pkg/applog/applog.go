package applog

import (
	"fmt"
	"log/slog"
	"os"
)

// Level — уровень логирования
type Level int

const (
	LevelDebug Level = -4
	LevelInfo  Level = 0
	LevelWarn  Level = 4
	LevelError Level = 8
)

type Logger struct {
	logger *slog.Logger
	level  slog.Level
}

// New создаёт логгер с указанным минимальным уровнем
// Паникует, если уровень неизвестен
func New(level Level) *Logger {
	slogLevel, err := loggerLevelFromLevel(level)
	if err != nil {
		panic(fmt.Sprintf("applog: %v", err))
	}

	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slogLevel})

	return &Logger{
		logger: slog.New(h),
		level:  slogLevel,
	}
}

// SetAsDefault устанавливает логгер как глобальный
func (l *Logger) SetAsDefault() {
	slog.SetDefault(l.logger)
}

// Debug логирует сообщение уровня Debug
func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info логирует сообщение уровня Info
func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn логирует сообщение уровня Warn
func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error логирует сообщение уровня Error
func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// loggerLevelFromLevel переводит кастомный уровень в slog.Level
func loggerLevelFromLevel(level Level) (slog.Level, error) {
	switch level {
	case LevelDebug:
		return slog.LevelDebug, nil
	case LevelInfo:
		return slog.LevelInfo, nil
	case LevelWarn:
		return slog.LevelWarn, nil
	case LevelError:
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unknown level: %d", level)
	}
}
