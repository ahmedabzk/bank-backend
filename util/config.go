package util

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DbDriver      string `mapstructure:"Db_Driver"`
	DbSource      string `mapstructure:"Db_Source"`
	ServerAddress string `mapstructure:"Server_Address"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot load the env files", err)
		return
	}

	viper.Unmarshal(&config)
	return
}
