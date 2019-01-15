package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func SetDefaults() {
	viper.SetDefault("spanner.enabled", false)
	viper.SetDefault("spanner.shards", 10)
	viper.SetDefault("pubsub.enabled", false)
	viper.SetDefault("rabbitmq.enabled", false)
	viper.SetDefault("cockroachdb.enabled", false)
	viper.SetDefault("pubsub.topic", "sheep")
	viper.SetDefault("pubsub.subscription", "sheep")
	viper.SetDefault("master", true)
	viper.SetDefault("worker", true)
	viper.SetDefault("direct", false)
	viper.SetDefault("service.port", 5309)
	viper.SetDefault("service.rest", 8080)
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
	SetDefaults()

	err := viper.ReadInConfig()
	if err != nil {
		log.Warn().Err(err).Msg("unable to read config, using defaults")
	}
}
