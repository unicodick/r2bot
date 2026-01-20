package r2

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type KeyGenerator struct{}

var invalidRe = regexp.MustCompile(`[^a-zA-Z0-9_.-]`)

func NewKeyGenerator() *KeyGenerator {
	return &KeyGenerator{}
}

func (g *KeyGenerator) Generate(filename string) string {
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	// sanitize filename
	name = invalidRe.ReplaceAllString(name, "_")

	return fmt.Sprintf("%d_%s%s", timestamp, name, ext)
}
