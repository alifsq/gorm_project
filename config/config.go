package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBCharset   string
	DBParseTime string
	DBLoc       string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		DBHost:      viper.GetString("DB_HOST"),
		DBPort:      viper.GetString("DB_PORT"),
		DBUser:      viper.GetString("DB_USER"),
		DBPassword:  viper.GetString("DB_PASSWORD"),
		DBName:      viper.GetString("DB_NAME"),
		DBCharset:   viper.GetString("DB_CHARSET"),
		DBParseTime: viper.GetString("DB_PARSE_TIME"),
		DBLoc:       viper.GetString("DB_LOC"),
	}, nil
}
