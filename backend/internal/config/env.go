package config

import (
	"fmt"
	"os"
	"strconv"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvReq(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return "", fmt.Errorf("env %s обязательный ключ и не может быть пустым", key)
	}
	return v, nil
}

func getEnvInt(key string, def int) (int, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("env %s: тип int, равен %q", key, v)
	}
	return n, nil
}
