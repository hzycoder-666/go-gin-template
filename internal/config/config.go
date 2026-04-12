package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Port int
	}

	Database struct {
		DSN string
	}

	Log struct {
		Level string
	}

	Jwt struct {
		Secret string
		Expire int
	}

	Ai struct {
		Token string
	}
}

var Global Config

func Init() error {
	viper.SetConfigFile("config/config.yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&Global)
	return err
}
