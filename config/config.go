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
	viper.SetDefault("pubsub.topic", "sheep")
	viper.SetDefault("pubsub.subscription", "sheep")
	viper.SetDefault("master", true)
	viper.SetDefault("worker", true)
}

func Setup(path string) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("sheep")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath("/etc/sheep/")
	if path != "" {
		viper.AddConfigPath(path)
	}
	setDefaults()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to read config")
	}
}
