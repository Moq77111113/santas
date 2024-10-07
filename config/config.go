package config

import (
	"strings"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Http     HTTP
		Database Database
		App      App
	}

	HTTP struct {
		Hostname string
		Port     uint16
	}

	App struct {
		Name          string
		ShowBanner    bool
		EncryptionKey string
	}
	Database struct {
		Driver     string
		Connection string
		Test       string
	}
)

func GetConfig() (Config, error) {
	var c Config

	// Load file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetEnvPrefix("chmoly")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return c, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}
