package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	Auth     AuthConfig
}

type AuthConfig struct {
	JWTSecret       string        `mapstructure:"jwt_secret"`
	AccessDuration  time.Duration `mapstructure:"access_duration"`
	RefreshDuration time.Duration `mapstructure:"refresh_duration"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	URL      string `mapstructure:"url"` // Full DATABASE_URL
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

func Load() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("log.level", "info")
	viper.SetDefault("auth.jwt_secret", "your-secret-key")
	viper.SetDefault("auth.access_duration", "15m")
	viper.SetDefault("auth.refresh_duration", "168h") // 7 days

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Parse DATABASE_URL if provided
	dbURL := viper.GetString("DATABASE_URL")
	if dbURL != "" {
		parsed, err := url.Parse(dbURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DATABASE_URL: %w", err)
		}

		cfg.Database.URL = dbURL
		cfg.Database.Host = parsed.Hostname()
		cfg.Database.Port = parsed.Port()

		if parsed.User != nil {
			cfg.Database.User = parsed.User.Username()
			cfg.Database.Password, _ = parsed.User.Password()
		}

		cfg.Database.Name = strings.TrimPrefix(parsed.Path, "/")

		// Parse query parameters for sslmode
		query := parsed.Query()
		if sslmode := query.Get("sslmode"); sslmode != "" {
			cfg.Database.SSLMode = sslmode
		}
	}

	// Explicitly map JWT_SECRET if present in env (priority over structured)
	if jwtSecret := viper.GetString("JWT_SECRET"); jwtSecret != "" {
		cfg.Auth.JWTSecret = jwtSecret
	}
	if jwtExpiry := viper.GetString("JWT_EXPIRY"); jwtExpiry != "" {
		if d, err := time.ParseDuration(jwtExpiry); err == nil {
			cfg.Auth.AccessDuration = d
		}
	}

	return &cfg, nil
}
