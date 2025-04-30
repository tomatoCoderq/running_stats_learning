package config

import (
	"log"
	"github.com/spf13/viper"
)

func InitConfig(filename string) *viper.Viper{
	config := viper.New()
	config.SetConfigName(filename)
	config.AddConfigPath(".") // current folder
	config.AddConfigPath("$HOME") // directory of hosting system
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal("Error while parsing config file")
	}
	return config
}