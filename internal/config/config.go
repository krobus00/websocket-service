package config

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

var (
	serviceName    = ""
	serviceVersion = ""
)

func ServiceName() string {
	return serviceName
}

func ServiceVersion() string {
	return serviceVersion
}

func ServerPort() string {
	str := viper.GetString("server.port")
	if str != "" {
		return str
	}
	return DefaultServerPort
}

func LogLevel() string {
	str := viper.GetString("log_level")
	if str != "" {
		return str
	}
	return DefaultLogLevel
}

func Environtment() string {
	str := viper.GetString("environtment")
	if str != "" {
		return str
	}
	return DefaultEnvirontment
}

func WebsocketDeadline() time.Duration {
	return viper.GetDuration("websocket.deadline")
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("config not found")
		}
		return err
	}
	return nil
}
