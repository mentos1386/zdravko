package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Port          string `validate:"required"`
	RootUrl       string `validate:"required,url"`
	DatabasePath  string `validate:"required"`
	SessionSecret string `validate:"required"`

	OAuth2 OAuth2 `validate:"required"`

	Temporal Temporal `validate:"required"`

	HealthChecks []Healthcheck
	CronJobs     []CronJob
}

type OAuth2 struct {
	ClientID            string   `validate:"required"`
	ClientSecret        string   `validate:"required"`
	Scopes              []string `validate:"required"`
	EndpointTokenURL    string   `validate:"required"`
	EndpointAuthURL     string   `validate:"required"`
	EndpointUserInfoURL string   `validate:"required"`
	EndpointLogoutURL   string   `validate:"required"`
}

type Temporal struct {
	DatabasePath  string `validate:"required"`
	ListenAddress string `validate:"required"`
	UIHost        string `validate:"required"`
	ServerHost    string `validate:"required"`
}

type HealthCheckHTTP struct {
	URL    string `validate:"required,url"`
	Method string `validate:"required,oneof=GET POST PUT"`
}

type HealthCheckTCP struct {
	Host string `validate:"required,hostname"`
	Port int    `validate:"required,gte=1,lte=65535"`
}

type Healthcheck struct {
	Name     string        `validate:"required"`
	Retries  int           `validate:"optional,gte=0"`
	Schedule string        `validate:"required,cron"`
	Timeout  time.Duration `validate:"required"`

	HTTP HealthCheckHTTP `validate:"required"`
	TCP  HealthCheckTCP  `validate:"required"`
}

type CronJob struct {
	Name     string        `validate:"required"`
	Schedule string        `validate:"required,cron"`
	Buffer   time.Duration `validate:"required"`
}

func NewConfig() *Config {
	viper.SetConfigName("zdravko")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/zdravko/")
	viper.AddConfigPath("$HOME/.zdravko")
	viper.AddConfigPath("$HOME/.config/zdravko")
	viper.AddConfigPath("$XDG_CONFIG_HOME/zdravko")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("port", "8000")
	viper.SetDefault("rooturl", "http://localhost:8000")
	viper.SetDefault("databasepath", "zdravko.db")
	viper.SetDefault("oauth2.scopes", "openid profile email")
	viper.SetDefault("temporal.databasepath", "temporal.db")
	viper.SetDefault("temporal.listenaddress", "0.0.0.0")
	viper.SetDefault("temporal.uihost", "127.0.0.1:8223")
	viper.SetDefault("temporal.serverhost", "127.0.0.1:7233")

	// Cant figure out the viper env, so lets just do it manually.
	viper.SetDefault("sessionsecret", os.Getenv("SESSION_SECRET"))
	viper.SetDefault("oauth2.clientid", os.Getenv("OAUTH2_CLIENT_ID"))
	viper.SetDefault("oauth2.clientsecret", os.Getenv("OAUTH2_CLIENT_SECRET"))
	viper.SetDefault("oauth2.scope", os.Getenv("OAUTH2_ENDPOINT_SCOPE"))
	viper.SetDefault("oauth2.endpointtokenurl", os.Getenv("OAUTH2_ENDPOINT_TOKEN_URL"))
	viper.SetDefault("oauth2.endpointauthurl", os.Getenv("OAUTH2_ENDPOINT_AUTH_URL"))
	viper.SetDefault("oauth2.endpointuserinfourl", os.Getenv("OAUTH2_ENDPOINT_USER_INFO_URL"))
	viper.SetDefault("oauth2.endpointlogouturl", os.Getenv("OAUTH2_ENDPOINT_LOGOUT_URL"))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	log.Println("Config file used: ", viper.ConfigFileUsed())

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}

	// OAuth2 scopes are space separated
	config.OAuth2.Scopes = strings.Split(viper.GetString("oauth2.scopes"), " ")

	// Validate config
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		log.Fatalf("Error validating config, %s", err)
	}

	fmt.Printf("Config: %+v\n", config)

	return config
}
