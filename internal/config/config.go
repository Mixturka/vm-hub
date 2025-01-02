package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ListenAddr     string
	SessionOptions *SessionOptions
	RedisUri       string
}

type SessionOptions struct {
	MaxAge          int
	SessionName     string
	SessionDomain   string
	SessionSecure   bool
	SessionHttpOnly bool
	SessionFolder   string
	SessionSecret   string
	CookiesSecret   string
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

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New(fmt.Sprintf("Error loading .env file: %s", err))
	}

	sessionLifeTimeStr := os.Getenv("SESSION_LIFETIME")
	sessionLifeTime, err := parseDuration(sessionLifeTimeStr)
	if err != nil {
		return nil, errors.New("Invalid SESSION_LIFETIME value")
	}

	sessionSecure, err := strconv.ParseBool(os.Getenv("SESSION_SECURE"))
	if err != nil {
		return nil, errors.New("Invalid SESSION_SECURE value")
	}

	sessionHttpOnly, err := strconv.ParseBool(os.Getenv("SESSION_HTTP_ONLY"))
	if err != nil {
		return nil, errors.New("Invalid SESSION_HTTP_ONLY value")
	}

	sessionOptions := &SessionOptions{
		MaxAge:          sessionLifeTime,
		SessionName:     os.Getenv("SESSION_NAME"),
		SessionDomain:   os.Getenv("SESSION_DOMAIN"),
		SessionSecure:   sessionSecure,
		SessionHttpOnly: sessionHttpOnly,
		SessionFolder:   os.Getenv("SESSION_FOLDER"),
		SessionSecret:   os.Getenv("SESSION_SECRET"),
		CookiesSecret:   os.Getenv("COOKIES_SECRET"),
	}

	return &Config{
		ListenAddr:     os.Getenv("LISTEN_ADDR"),
		SessionOptions: sessionOptions,
		RedisUri:       os.Getenv("REDIS_URI"),
	}, nil
}
