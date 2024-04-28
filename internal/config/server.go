package config

import (
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port                 string `validate:"required"`
	RootUrl              string `validate:"required,url"`
	SqliteDatabasePath   string `validate:"required"`
	KeyValueDatabasePath string `validate:"required"`
	SessionSecret        string `validate:"required"`

	Jwt    ServerJwt    `validate:"required"`
	OAuth2 ServerOAuth2 `validate:"required"`

	Temporal ServerTemporal `validate:"required"`
}

type ServerJwt struct {
	PrivateKey string `validate:"required"`
	PublicKey  string `validate:"required"`
}

type ServerOAuth2 struct {
	ClientID            string   `validate:"required"`
	ClientSecret        string   `validate:"required"`
	Scopes              []string `validate:"required"`
	EndpointTokenURL    string   `validate:"required"`
	EndpointAuthURL     string   `validate:"required"`
	EndpointUserInfoURL string   `validate:"required"`
	EndpointLogoutURL   string   // Optional as not all SSO support this.
}

type ServerTemporal struct {
	UIHost     string `validate:"required"`
	ServerHost string `validate:"required"`
}

func NewServerConfig() *ServerConfig {
	v := newViper()

	// Set defaults
	v.SetDefault("port", GetEnvOrDefault("PORT", "8000"))
	v.SetDefault("rooturl", GetEnvOrDefault("ROOT_URL", "http://localhost:8000"))
	v.SetDefault("sqlitedatabasepath", GetEnvOrDefault("SQLITE_DATABASE_PATH", "zdravko.db"))
	v.SetDefault("keyvaluedatabasepath", GetEnvOrDefault("KEYVALUE_DATABASE_PATH", "zdravko_kv.db"))
	v.SetDefault("sessionsecret", os.Getenv("SESSION_SECRET"))
	v.SetDefault("temporal.uihost", GetEnvOrDefault("TEMPORAL_UI_HOST", "127.0.0.1:8223"))
	v.SetDefault("temporal.serverhost", GetEnvOrDefault("TEMPORAL_SERVER_HOST", "127.0.0.1:7233"))
	v.SetDefault("jwt.privatekey", os.Getenv("JWT_PRIVATE_KEY"))
	v.SetDefault("jwt.publickey", os.Getenv("JWT_PUBLIC_KEY"))
	v.SetDefault("oauth2.clientid", os.Getenv("OAUTH2_CLIENT_ID"))
	v.SetDefault("oauth2.clientsecret", os.Getenv("OAUTH2_CLIENT_SECRET"))
	v.SetDefault("oauth2.scopes", GetEnvOrDefault("OAUTH2_ENDPOINT_SCOPES", "openid profile email"))
	v.SetDefault("oauth2.endpointtokenurl", os.Getenv("OAUTH2_ENDPOINT_TOKEN_URL"))
	v.SetDefault("oauth2.endpointauthurl", os.Getenv("OAUTH2_ENDPOINT_AUTH_URL"))
	v.SetDefault("oauth2.endpointuserinfourl", os.Getenv("OAUTH2_ENDPOINT_USER_INFO_URL"))
	v.SetDefault("oauth2.endpointlogouturl", GetEnvOrDefault("OAUTH2_ENDPOINT_LOGOUT_URL", ""))

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ignore
		} else {
			log.Fatalf("Error reading config file, %s", err)
		}
	}
	log.Println("Config file used: ", v.ConfigFileUsed())

	config := &ServerConfig{}
	err = v.Unmarshal(config)
	if err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}

	// OAuth2 scopes are space separated
	config.OAuth2.Scopes = strings.Split(v.GetString("oauth2.scopes"), " ")

	// Validate config
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		log.Fatalf("Error validating config, %s", err)
	}

	return config
}
