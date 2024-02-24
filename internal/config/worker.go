package config

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	Token  string `validate:"required"`
	ApiUrl string `validate:"required"`
}

func NewWorkerConfig() *WorkerConfig {
	v := newViper()

	// Set defaults
	v.SetDefault("token", os.Getenv("WORKER_GROUP_TOKEN"))
	v.SetDefault("apiurl", GetEnvOrDefault("WORKER_API_URL", "http://localhost:8000"))

	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// ignore
		} else {
			log.Fatalf("Error reading config file, %s", err)
		}
	}
	log.Println("Config file used: ", v.ConfigFileUsed())

	config := &WorkerConfig{}
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
