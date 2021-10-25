package config

import "github.com/spf13/viper"

type Config struct {
	ESuri   string `mapstructure:"ES_URI"`
	ESindex string `mapstructure:"ES_INDEX"`
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
