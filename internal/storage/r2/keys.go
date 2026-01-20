package r2

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type KeyGenerator struct{}

func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{}
}

func (g *KeyGenerator) Generate(filename string) string {
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	// sanitize filename
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")

	return fmt.Sprintf("%d_%s%s", timestamp, name, ext)
}
