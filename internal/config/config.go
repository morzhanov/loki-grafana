package config

import "github.com/spf13/viper"

type Config struct {
	LokiClientURL string `mapstructure:"LOKI_CLIENT_URL"`
	Version       string `mapstructure:"VERSION"`
}

func NewConfig() (config *Config, err error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
