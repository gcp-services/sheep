package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("spanner.enabled", false)
	viper.SetDefault("pubsub.enabled", false)
	viper.SetDefault("rabbitmq.enabled", false)
	viper.SetDefault("cockroachdb.enabled", false)
}

func Setup() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("sheep")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath("/etc/sheep/")
	setDefaults()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read config")
	}
}
