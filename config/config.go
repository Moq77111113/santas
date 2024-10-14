package config

import (
	"strings"

	"github.com/labstack/gommon/log"
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
		Name                string
		ShowBanner          bool
		EncryptionKey       string
		ExclusionPercentage float32
		RoundUp             bool
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

	check(&c)

	return c, nil
}

func check(c *Config) {

	if c.App.ExclusionPercentage < 0 || c.App.ExclusionPercentage > 1 {
		log.Warnf("ExclusionPercentage is out of range, setting to default")
		c.App.ExclusionPercentage = 0.15
	}
}
