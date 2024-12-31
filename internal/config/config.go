package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr string

	CookiesSecret   string
	SessionSecret   string
	SessionName     string
	SessionDomain   string
	SessionLifeTime int
	SessionHttpOnly bool
	SessionSecure   bool
	SessionFolder   string

	RedisUri string
}

/*
Parses duration with unit e.g. "3d", "15h", "12m" and returns result duration in seconds
with possible error. If no unit provided parses as seconds.
*/
func parseDuration(duration string) (int, error) {
	multipliers := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": 24 * time.Hour,
	}

	if duration == "" {
		return 0, nil
	}

	for unit, multiplier := range multipliers {
		if strings.HasSuffix(duration, unit) {
			valueStr := strings.TrimSuffix(duration, unit)
			value, err := strconv.Atoi(valueStr)

			if err != nil {
				return 0, err
			}

			return int(time.Duration(value) * multiplier / time.Second), nil
		}
	}

	value, err := strconv.Atoi(duration)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	sessionLifeTimeStr := os.Getenv("SESSION_LIFETIME")
	sessionLifeTime, err := parseDuration(sessionLifeTimeStr)
	if err != nil {
		slog.Error("Invalid SESSION_LIFETIME value", "value", sessionLifeTimeStr, "error", err)
		os.Exit(1)
	}

	sessionSecure, err := strconv.ParseBool(os.Getenv("SESSION_SECURE"))
	if err != nil {
		slog.Error("Invalid SESSION_SECURE value", "value", os.Getenv("SESSION_SECURE"), "error", err)
		os.Exit(1)
	}

	sessionHttpOnly, err := strconv.ParseBool(os.Getenv("SESSION_HTTP_ONLY"))
	if err != nil {
		slog.Error("Invalid SESSION_HTTP_ONLY value", "value", os.Getenv("SESSION_HTTP_ONLY"), "error", err)
		os.Exit(1)
	}

	return &Config{
		ListenAddr:      os.Getenv("LISTEN_ADDR"),
		CookiesSecret:   os.Getenv("COOKIES_SECRET"),
		SessionSecret:   os.Getenv("SESSION_SECRET"),
		SessionName:     os.Getenv("SESSION_NAME"),
		SessionDomain:   os.Getenv("SESSION_DOMAIN"),
		SessionLifeTime: sessionLifeTime,
		SessionHttpOnly: sessionHttpOnly,
		SessionSecure:   sessionSecure,
		SessionFolder:   os.Getenv("SESSION_FOLDER"),
		RedisUri:        os.Getenv("REDIS_URI"),
	}
}
