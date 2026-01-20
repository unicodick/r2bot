package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken   string
	AllowedIDs []int64

	R2AccountID string
	R2AccessKey string
	R2SecretKey string
	R2Bucket    string
	R2PublicURL string
}

func Load() (*Config, error) {
	cfg := &Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		R2AccountID: os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKey: os.Getenv("R2_ACCESS_KEY"),
		R2SecretKey: os.Getenv("R2_SECRET_KEY"),
		R2Bucket:    os.Getenv("R2_BUCKET"),
		R2PublicURL: os.Getenv("R2_PUBLIC_URL"),
	}

	if cfg.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is required")
	}

	allowedIDsStr := os.Getenv("ALLOWED_IDS")
	if allowedIDsStr == "" {
		return nil, fmt.Errorf("ALLOWED_IDS is required")
	}

	ids := strings.Split(allowedIDsStr, ",")
	for _, idStr := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid user id: %s", idStr)
		}
		cfg.AllowedIDs = append(cfg.AllowedIDs, id)
	}

	if cfg.R2AccountID == "" || cfg.R2AccessKey == "" || cfg.R2SecretKey == "" || cfg.R2Bucket == "" || cfg.R2PublicURL == "" {
		return nil, fmt.Errorf("all r2 configuration parameters are required")
	}

	return cfg, nil
}

func (c *Config) IsUserAllowed(userID int64) bool {
	for _, id := range c.AllowedIDs {
		if id == userID {
			return true
		}
	}
	return false
}
