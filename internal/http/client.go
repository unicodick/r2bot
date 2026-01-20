package http

import (
	"net/http"
	"sync"
	"time"
)

var (
	instance *http.Client
	once     sync.Once
)

// client returns singleton http client
func Client() *http.Client {
	once.Do(func() {
		instance = &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    90 * time.Second,
				DisableCompression: true,
			},
		}
	})
	return instance
}
