package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabaseURL   string
	HTTPAddr      string
	GRPCAddr      string
	PublicLink    string
	InternalToken string
	EasySMSURL    string
	EasySMSKey    string
	EasySMSMock   bool
	RevealOTP     bool
}

func Load() Config {
	return Config{
		DatabaseURL:   env("DATABASE_URL", "postgresql://connect:connect@localhost:5434/connect?sslmode=disable"),
		HTTPAddr:      env("HTTP_ADDR", ":8080"),
		GRPCAddr:      env("GRPC_ADDR", ":50051"),
		PublicLink:    strings.TrimRight(env("PUBLIC_LINK_BASE", "http://localhost:8080"), "/"),
		InternalToken: env("INTERNAL_TOKEN", "change-me-internal-token"),
		EasySMSURL:    strings.TrimRight(env("EASYSMS_BASE_URL", "https://easysms.co"), "/"),
		EasySMSKey:    strings.TrimSpace(os.Getenv("EASYSMS_API_KEY")),
		EasySMSMock:   envBool("EASYSMS_MOCK", true),
		RevealOTP:     envBool("MOCK_REVEAL_OTP", true),
	}
}

func env(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return parsed
}
