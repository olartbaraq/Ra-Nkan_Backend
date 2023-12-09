package utils

import "github.com/spf13/viper"

type Config struct {
	DBdriver     string `mapstructure:"DB_DRIVER"`
	DBsource     string `mapstructure:"DB_SOURCE"`
	DBsourceLive string `mapstructure:"DB_SOURCE_LIVE"`
	SigningKey   string `mapstructure:"SIGNING_KEY"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
