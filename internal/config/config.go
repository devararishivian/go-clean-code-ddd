package config

import (
	"github.com/spf13/viper"
)

type database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type server struct {
	Address string
}

var (
	Database  database
	Server    server
	JWTSecret string
)

func LoadConfig(path string) error {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("database", &Database); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("server", &Server); err != nil {
		return err
	}

	JWTSecret = viper.GetString("jwtSecret")

	return nil
}
