package internal

import (
	"github.com/spf13/viper"
)

type Config struct {
    DBUser     string `mapstructure:"db_user"`
    DBPassword string `mapstructure:"db_password"`
    DBName     string `mapstructure:"db_name"`
}

func LoadConfig(filename string) (Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)

	viper.SetEnvPrefix("RCB")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	if err != nil {
		return Config{}, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
