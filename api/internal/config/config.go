package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
}

func Load() (*Config, error) {
	var missing []string

	cfg := &Config{
		DBHost:     require("DB_HOST", &missing),
		DBPort:     require("DB_PORT", &missing),
		DBUser:     require("DB_USER", &missing),
		DBPassword: require("DB_PASSWORD", &missing),
		DBName:     require("DB_NAME", &missing),
		Port:       require("PORT", &missing),
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC&charset=utf8mb4",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func require(key string, missing *[]string) string {
	value := os.Getenv(key)
	if value == "" {
		*missing = append(*missing, key)
	}
	return value
}
