package config

import (
	"github.com/spf13/viper"
)

func GetEnv[T comparable](file string, env *T) error {
	viper.SetConfigName(file)
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.AutomaticEnv()

	if err = viper.Unmarshal(env); err != nil {
		return err
	}

	return nil
}
