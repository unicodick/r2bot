package config

import "os"

type Config struct {
	BotToken   string
	AllowedIDs map[int64]struct{}

	R2AccountID string
	R2AccessKey string
	R2SecretKey string
	R2Bucket    string
	R2PublicURL string
}

func Load() (*Config, error) {
	validator := NewValidator()

	allowedIDs, err := validator.LoadAllowedIDs()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		AllowedIDs:  allowedIDs,
		R2AccountID: os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKey: os.Getenv("R2_ACCESS_KEY"),
		R2SecretKey: os.Getenv("R2_SECRET_KEY"),
		R2Bucket:    os.Getenv("R2_BUCKET"),
		R2PublicURL: os.Getenv("R2_PUBLIC_URL"),
	}

	if err := validator.Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) IsUserAllowed(userID int64) bool {
	_, ok := c.AllowedIDs[userID]
	return ok
}
