package config

import (
	"fmt"
	"os"
	"strconv"
)

// ////////////////// Работа с переменными окружения, чтение
func getEnvReq(keyName string) (string, error) {
	val, ok := os.LookupEnv(keyName)
	if !ok {
		return "", fmt.Errorf("key not found in .env: %s", keyName)
	}

	if val == "" {
		return "", fmt.Errorf("key is empty in .env: %s", keyName)
	}

	return val, nil
}

func getEnv(keyName, defaultValue string) string {
	val := os.Getenv(keyName)

	if val == "" {
		return defaultValue
	}

	return val
}

func getEnvInt(keyName string, defaultValue int) int {
	val := os.Getenv(keyName)

	valInt, err := strconv.Atoi(val)

	if err != nil {
		return defaultValue
	}

	return valInt
}

func getEnvIntReq(keyName string) (int, error) {
	val := os.Getenv(keyName)

	valInt, err := strconv.Atoi(val)

	if err != nil {
		return 0, fmt.Errorf("key from .env not int: %s, value: %s", keyName, val)
	}

	return valInt, nil
}
