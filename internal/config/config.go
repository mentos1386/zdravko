package config

import (
	"os"

	"github.com/spf13/viper"
)

func GetEnvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

func newViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("zdravko")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/zdravko/")
	v.AddConfigPath("$HOME/.zdravko")
	v.AddConfigPath("$HOME/.config/zdravko")
	v.AddConfigPath("$XDG_CONFIG_HOME/zdravko")
	v.AddConfigPath(".")
	return v
}
