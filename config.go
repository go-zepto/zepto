package zepto

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/go-zepto/zepto/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig    `json:"app" mapstructure:"app"`
	Server ServerConfig `json:"server" mapstructure:"server"`
	Logger LoggerConfig `json:"logger" mapstructure:"logger"`
}

// App Config is the configuration for the Zepto App
type AppConfig struct {
	// App Name (e.g. my-project) [required]
	Name string `json:"name" mapstructure:"name"`

	// App Version (e.g. "1.0.0")
	Version string `json:"version" mapstructure:"version"`

	// App Session
	Session SessionConfig `json:"session" mapstructure:"session"`

	// App Webpack Enabled [default: true]
	WebpackEnabled bool `json:"webpack_enabled" mapstructure:"webpack_enabled"`
}

// SessionConfig is the configuration for the default session.
//
// You can change the session provider, but this configuration will be ignored
type SessionConfig struct {
	// Session Name [default="zsid"]
	Name string `json:"name" mapstructure:"name"`

	// Unique protected string used to hash session [required]
	Secret string `json:"secret" mapstructure:"secret"`
}

// ServerConfig is the configuration for the default server.
//
// You can change the server instance (*http.Server), but this configuration will be ignored
type ServerConfig struct {
	// Server Host (e.g. localhost)
	Host string `json:"host" mapstructure:"host"`

	// Server Port (0 to 65535)
	Port int `json:"port" mapstructure:"port"`

	// Server read timeout (ms)
	ReadTimeout int `json:"read_timeout" mapstructure:"read_timeout"`

	// Server write timeout (ms)
	WriteTimeout int `json:"write_timeout" mapstructure:"write_timeout"`
}

// LoggerConfig is the configuration for the default logger.
//
// You can change the logger instance, but this configuration will be ignored
type LoggerConfig struct {
	// Level (e.g. "trace", "debug", "info", "warning", "error", "fatal", "panic")
	//
	// default: "debug"
	Level string `json:"level" mapstructure:"level"`

	// Enable colored log
	//
	// default: true
	Colors bool `json:"colors" mapstructure:"colors"`

	// Show timestamp in log
	//
	// default: true
	Timestamp bool `json:"timestamp" mapstructure:"timestamp"`
}

func NewDefaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Name:    "zepto",
			Version: "1.0.0",
			Session: SessionConfig{
				Name: "zsid",
			},
			WebpackEnabled: true,
		},
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8000,
			ReadTimeout:  15000,
			WriteTimeout: 15000,
		},
		Logger: LoggerConfig{
			Level:     "debug",
			Colors:    true,
			Timestamp: true,
		},
	}
}

func NewConfigFromFile(config *Config) (*Config, error) {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	env := utils.GetEnv("ZEPTO_ENV", "development")
	configFile := fmt.Sprintf("config/%s.yml", env)
	viper.SetConfigName(configFile)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("ZEPTO")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			bold := color.New(color.Bold)
			l.Warnf("Configuration file not found in %s. Using default zepto config...\n", bold.Sprint(configFile))
			return config, nil
		}
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
