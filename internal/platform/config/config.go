package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// App
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`

	// Database
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`

	// Web
	Cookie CookieConfig `mapstructure:"cookie"`

	// JWT
	JWT JWTConfig `mapstructure:"jwt"`
}

type AppConfig struct {
	Env string `mapstructure:"env"`
}

type ServerConfig struct {
	Host           string   `mapstructure:"host"`
	Port           string   `mapstructure:"port"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
	Mode           string   `mapstructure:"mode"`
}

type DatabaseConfig struct {
	URI  string `mapstructure:"uri"`
	Name string `mapstructure:"name"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
}

type CookieConfig struct {
	HttpOnly bool   `mapstructure:"http_only"`
	Secure   bool   `mapstructure:"secure"`
	SameSite string `mapstructure:"same_site"`
	MaxAge   int    `mapstructure:"max_age"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg *Config
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return cfg, nil
}
