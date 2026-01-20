package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(cfg *Config) error {
	if cfg.BotToken == "" {
		return errors.New("BOT_TOKEN is required")
	}

	if cfg.R2AccountID == "" {
		return errors.New("R2_ACCOUNT_ID is required")
	}

	if cfg.R2AccessKey == "" {
		return errors.New("R2_ACCESS_KEY is required")
	}

	if cfg.R2SecretKey == "" {
		return errors.New("R2_SECRET_KEY is required")
	}

	if cfg.R2Bucket == "" {
		return errors.New("R2_BUCKET is required")
	}

	if cfg.R2PublicURL == "" {
		return errors.New("R2_PUBLIC_URL is required")
	}

	return nil
}

func (v *Validator) LoadAllowedIDs() (map[int64]struct{}, error) {
	allowedIDsStr := os.Getenv("ALLOWED_IDS")
	if allowedIDsStr == "" {
		return make(map[int64]struct{}), nil
	}

	ids := strings.Split(allowedIDsStr, ",")
	allowedIDs := make(map[int64]struct{}, len(ids))

	for _, idStr := range ids {
		id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			return nil, err
		}
		allowedIDs[id] = struct{}{}
	}

	return allowedIDs, nil
}
