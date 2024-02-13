package internal

import (
	"os"
	"strings"
)

type Config struct {
	PORT     string
	ROOT_URL string // Needed for oauth2 redirect

	SESSION_SECRET string

	OAUTH2_CLIENT_ID              string
	OAUTH2_CLIENT_SECRET          string
	OAUTH2_SCOPES                 []string
	OAUTH2_ENDPOINT_TOKEN_URL     string
	OAUTH2_ENDPOINT_AUTH_URL      string
	OAUTH2_ENDPOINT_USER_INFO_URL string
	OAUTH2_ENDPOINT_LOGOUT_URL    string

	ZDRAVKO_DATABASE_PATH string

	TEMPORAL_DATABASE_PATH  string
	TEMPORAL_LISTEN_ADDRESS string
	TEMPORAL_UI_HOST        string
	TEMPORAL_SERVER_HOST    string
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

		SESSION_SECRET: getEnvRequired("SESSION_SECRET"),

		OAUTH2_CLIENT_ID:              getEnvRequired("OAUTH2_CLIENT_ID"),
		OAUTH2_CLIENT_SECRET:          getEnvRequired("OAUTH2_CLIENT_SECRET"),
		OAUTH2_SCOPES:                 strings.Split(getEnv("OAUTH2_SCOPES", "openid,profile,email"), ","),
		OAUTH2_ENDPOINT_TOKEN_URL:     getEnvRequired("OAUTH2_ENDPOINT_TOKEN_URL"),
		OAUTH2_ENDPOINT_AUTH_URL:      getEnvRequired("OAUTH2_ENDPOINT_AUTH_URL"),
		OAUTH2_ENDPOINT_USER_INFO_URL: getEnvRequired("OAUTH2_ENDPOINT_USER_INFO_URL"),
		OAUTH2_ENDPOINT_LOGOUT_URL:    getEnvRequired("OAUTH2_ENDPOINT_LOGOUT_URL"),

		ZDRAVKO_DATABASE_PATH: getEnv("ZDRAVKO_DATABASE_PATH", "zdravko.db"),

		TEMPORAL_DATABASE_PATH:  getEnv("TEMPORAL_DATABASE_PATH", "temporal.db"),
		TEMPORAL_LISTEN_ADDRESS: getEnv("TEMPORAL_LISTEN_ADDRESS", "0.0.0.0"),
		TEMPORAL_UI_HOST:        getEnv("TEMPORAL_UI_HOST", "127.0.0.1:8223"),
		TEMPORAL_SERVER_HOST:    getEnv("TEMPORAL_SERVER_HOST", "127.0.0.1:7233"),
	}
}
