package config

import (
	"log/slog"
	"net"
	"net/url"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32
}

func (c DatabaseConfig) DSN() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   c.DBName,
	}
	q := u.Query()
	q.Set("sslmode", c.SSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

func (c DatabaseConfig) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("host", c.Host),
		slog.String("port", c.Port),
		slog.String("db", c.DBName),
		slog.String("sslmode", c.SSLMode),
		slog.Int("max_conns", int(c.MaxConns)),
	)
}

func newDatabaseConfig() (DatabaseConfig, error) {
	user, err := getEnvReq("DB_USER")
	if err != nil {
		return DatabaseConfig{}, err
	}
	password, err := getEnvReq("DB_PASSWORD")
	if err != nil {
		return DatabaseConfig{}, err
	}
	dbName, err := getEnvReq("DB_NAME")
	if err != nil {
		return DatabaseConfig{}, err
	}

	maxConns, err := getEnvInt("DB_MAX_CONNS", 10)
	if err != nil {
		return DatabaseConfig{}, err
	}

	return DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		MaxConns: int32(maxConns),
	}, nil
}
