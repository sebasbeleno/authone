package env

import (
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func GetBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value == "true"
}

func GetInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	valAsInt, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return valAsInt
}
