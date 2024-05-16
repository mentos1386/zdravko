package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type TemporalConfig struct {
	DatabasePath  string `validate:"required"`
	ListenAddress string `validate:"required"`

	Jwt TemporalJwt `validate:"required"`
}

type TemporalJwt struct {
	PublicKey string `validate:"required"`
}

func NewTemporalConfig() *TemporalConfig {
	v := newViper()

	// Set defaults
	v.SetDefault("databasepath", GetEnvOrDefault("TEMPORAL_DATABASE_PATH", "store/temporal.db"))
	v.SetDefault("listenaddress", GetEnvOrDefault("TEMPORAL_LISTEN_ADDRESS", "0.0.0.0"))
	v.SetDefault("jwt.publickey", os.Getenv("JWT_PUBLIC_KEY"))

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ignore
		} else {
			log.Fatalf("Error reading config file, %s", err)
		}
	}
	log.Println("Config file used: ", v.ConfigFileUsed())

	config := &TemporalConfig{}
	err = v.Unmarshal(config)
	if err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}

	// Validate config
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(config)
	if err != nil {
		log.Fatalf("Error validating config, %s", err)
	}

	return config
}
