package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var v *viper.Viper

func Init(configPath string) error {
	v = viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return nil
}

func GetString(key string) string {
	value := v.GetString(key)
	if value == "" {
		panic(fmt.Errorf("%s is a required config", key))
	}
	return strings.TrimSpace(value)
}

func GetStringSlice(key string) []string {
	value := v.GetStringSlice(key)
	if len(value) == 0 {
		panic(fmt.Errorf("%s is a required config and must have at least one value", key))
	}
	return value
}

func GetDuration(key string) time.Duration {
	value := v.GetDuration(key)
	if value == 0 {
		panic(fmt.Errorf("%s is a required config", key))
	}
	return value
}

func GetInt(key string) int {
	if !v.IsSet(key) {
		panic(fmt.Errorf("%s is a required config", key))
	}
	return v.GetInt(key)
}

func GetFloat64(key string) float64 {
	if !v.IsSet(key) {
		panic(fmt.Errorf("%s is a required config", key))
	}
	return v.GetFloat64(key)
}

func GetBool(key string) bool {
	if !v.IsSet(key) {
		panic(fmt.Errorf("%s is a required config", key))
	}
	return v.GetBool(key)
}

// GetStringOptional returns string value or empty string if not set
func GetStringOptional(key string) string {
	return strings.TrimSpace(v.GetString(key))
}

// GetIntOptional returns int value or 0 if not set
func GetIntOptional(key string) int {
	return v.GetInt(key)
}

// GetDurationOptional returns duration or 0 if not set
func GetDurationOptional(key string) time.Duration {
	return v.GetDuration(key)
}
