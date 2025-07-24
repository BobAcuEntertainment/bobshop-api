package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// Server
	Server ServerConfig `mapstructure:"server"`

	// Database
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`

	// Web
	Cookie CookieConfig `mapstructure:"cookie"`

	// JWT
	JWT JWTConfig `mapstructure:"jwt"`
}

type ServerConfig struct {
	Host           string   `mapstructure:"host"`
	Port           string   `mapstructure:"port"`
	TrustedProxies []string `mapstructure:"trusted_proxies"`
}

type DatabaseConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
}

type RedisConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours string `mapstructure:"expiration_hours"`
}

type CookieConfig struct {
	HttpOnly bool   `mapstructure:"http_only"`
	Secure   bool   `mapstructure:"secure"`
	SameSite string `mapstructure:"same_site"`
	MaxAge   string `mapstructure:"max_age"`
}

func IsDevelopment() bool {
	return viper.GetString("APP_ENV") == "development"
}

func IsStaging() bool {
	return viper.GetString("APP_ENV") == "staging"
}

func IsProduction() bool {
	return viper.GetString("APP_ENV") == "production"
}

func LoadConfig(envPath, configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if IsDevelopment() {
		viper.SetConfigFile(envPath)
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading env file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
