package config

import "fmt"

// DatabaseConfig конфиг бд
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

// validate проверяет корректность конфигурации БД
func (c *DatabaseConfig) validate() error {
	if c.User == "" {
		return fmt.Errorf("DB_USER не может быть пустым")
	}
	if c.Password == "" {
		return fmt.Errorf("DB_PASSWORD не может быть пустым")
	}
	if c.DBName == "" {
		return fmt.Errorf("DB_NAME не может быть пустым")
	}
	return nil
}

func newDatabaseConfig() (*DatabaseConfig, error) {
	user, err := getEnvReq("DB_USER")
	if err != nil {
		return nil, err
	}

	password, err := getEnvReq("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnvReq("DB_NAME")
	if err != nil {
		return nil, err
	}

	sslModeInt := getEnvInt("DB_SSL_MODE", 0)
	sslMode := ""
	if sslModeInt == 0 {
		sslMode = "disable"
	} else {
		sslMode = "enable"
	}

	maxConns := getEnvInt("DB_MAX_CONNS", 10)

	dbConfig := &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
		MaxConns: int32(maxConns),
	}

	if err := dbConfig.validate(); err != nil {
		return nil, err
	}

	return dbConfig, nil
}
