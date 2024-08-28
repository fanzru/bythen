package env

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// Load loads environment variables from the specified filenames.
func Load(filenames ...string) error {
	return godotenv.Load(filenames...)
}

// GetEnv returns environment variable value by given key, or default value if
// not found.
func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// FillStruct populates the fields of a struct with environment variables using envconfig.
func FillStruct(ctx context.Context, cfg any) error {
	return envconfig.Process(ctx, cfg)
}
