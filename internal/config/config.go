package config

import (
	"os"
	"strings"

	"mkluxe-backend/internal/utils"
)

type Config struct {
	GinMode          string
	Port             string
	JWTAccessSecret  string
	JWTRefreshSecret string
	FrontendURLs     []string // 👈 List of allowed frontend origins
	IsProduction     bool
	R2AccountID      string
	R2AccessKey      string
	R2SecretKey      string
	R2BucketName     string
	R2PublicBaseURL  string
}

// CookieConfig holds individual configuration flags for SetCookie
type CookieConfig struct {
	Name     string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HttpOnly bool
}

// Load reads values from the environment into the Config struct
func Load() *Config {
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = "debug"
	}

	// Read comma-separated frontend URLs
	urlsStr := os.Getenv("FRONTEND_URLS")
	if urlsStr == "" {
		urlsStr = "http://localhost:3000,http://localhost:3001" // Defaults for local dev
	}

	// Split and trim the URLs into a slice
	var frontendURLs []string
	for _, u := range strings.Split(urlsStr, ",") {
		frontendURLs = append(frontendURLs, strings.TrimSpace(u))
	}

	return &Config{
		GinMode:          mode,
		Port:             utils.GetEnv("PORT", "8080"),
		JWTAccessSecret:  utils.GetEnv("JWT_ACCESS_SECRET", "default_access_secret_change_me"),
		JWTRefreshSecret: utils.GetEnv("JWT_REFRESH_SECRET", "default_refresh_secret_change_me"),
		FrontendURLs:     frontendURLs,
		IsProduction:     mode == "release",
		R2AccountID:      utils.GetEnv("R2_ACCOUNT_ID", ""),
		R2AccessKey:      utils.GetEnv("R2_ACCESS_KEY", ""),
		R2SecretKey:      utils.GetEnv("R2_SECRET_KEY", ""),
		R2BucketName:     utils.GetEnv("R2_BUCKET_NAME", ""),
		R2PublicBaseURL:  utils.GetEnv("R2_PUBLIC_BASE_URL", ""),
	}
}

// GetAccessCookieConfig returns standard secure configurations for short-lived access tokens
func (c *Config) GetAccessCookieConfig() *CookieConfig {
	domain := ""
	if c.IsProduction {
		domain = ".mk-luxe-divine.in"
	}
	return &CookieConfig{
		Name:     "access_token",
		MaxAge:   900,
		Path:     "/",
		Domain:   domain,
		Secure:   c.IsProduction,
		HttpOnly: true,
	}
}

// GetRefreshCookieConfig returns standard secure configurations for long-lived refresh tokens
func (c *Config) GetRefreshCookieConfig() *CookieConfig {
	domain := ""
	if c.IsProduction {
		domain = ".mk-luxe-divine.in"
	}
	return &CookieConfig{
		Name:     "refresh_token",
		MaxAge:   604800,
		Path:     "/",
		Domain:   domain,
		Secure:   c.IsProduction,
		HttpOnly: true,
	}
}
