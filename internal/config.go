package internal

import (
	"os"
	"strings"
)

type Config struct {
	PORT     string
	ROOT_URL string // Needed for oauth2 redirect

	SQLITE_DB_PATH string

	SESSION_SECRET string

	OAUTH2_CLIENT_ID              string
	OAUTH2_CLIENT_SECRET          string
	OAUTH2_SCOPES                 []string
	OAUTH2_ENDPOINT_TOKEN_URL     string
	OAUTH2_ENDPOINT_AUTH_URL      string
	OAUTH2_ENDPOINT_USER_INFO_URL string
	OAUTH2_ENDPOINT_LOGOUT_URL    string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvRequired(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Environment variable " + key + " is required")
}

func NewConfig() *Config {
	return &Config{
		PORT:     getEnv("PORT", "8000"),
		ROOT_URL: getEnvRequired("ROOT_URL"),

		SQLITE_DB_PATH: getEnv("SQLITE_DB_PATH", "zdravko.db"),
		SESSION_SECRET: getEnvRequired("SESSION_SECRET"),

		OAUTH2_CLIENT_ID:              getEnvRequired("OAUTH2_CLIENT_ID"),
		OAUTH2_CLIENT_SECRET:          getEnvRequired("OAUTH2_CLIENT_SECRET"),
		OAUTH2_SCOPES:                 strings.Split(getEnvRequired("OAUTH2_SCOPES"), ","),
		OAUTH2_ENDPOINT_TOKEN_URL:     getEnvRequired("OAUTH2_ENDPOINT_TOKEN_URL"),
		OAUTH2_ENDPOINT_AUTH_URL:      getEnvRequired("OAUTH2_ENDPOINT_AUTH_URL"),
		OAUTH2_ENDPOINT_USER_INFO_URL: getEnvRequired("OAUTH2_ENDPOINT_USER_INFO_URL"),
		OAUTH2_ENDPOINT_LOGOUT_URL:    getEnvRequired("OAUTH2_ENDPOINT_LOGOUT_URL"),
	}
}
