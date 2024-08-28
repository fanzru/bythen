package env_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/fanzru/bythen/pkg/env"
)

// TestLoad tests the Load function manually
func TestLoad(t *testing.T) {

	test := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "success",
			path:    filepath.Join("testdata", "test.env"),
			wantErr: false,
		},
		{
			name:    "error",
			path:    filepath.Join("testdata", "inexist.env"),
			wantErr: true,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := env.Load(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

// TestFillStruct tests the FillStruct function
func TestFillStruct(t *testing.T) {
	t.Parallel()
	type Config struct {
		AppName    string `env:"APP_NAME_TEST"`
		Port       int    `env:"PORT_TEST"`
		Debug      bool   `env:"DEBUG_TEST"`
		MaxRetries uint   `env:"MAX_RETRIES_TEST"`
		InvalidKey string
		BlankKey   string            `env:"BLANK_KEY"`
		AnyKey     map[string]string `env:"ANY_KEY"`
	}

	tests := []struct {
		name  string
		key   string
		value string
		want  Config
	}{
		{
			name:  "map[string]string",
			key:   "ANY_KEY",
			value: `key:value`,
			want:  Config{AnyKey: map[string]string{"key": "value"}},
		},
		{
			name:  "blank string",
			key:   "BLANK_KEY",
			value: "",
			want:  Config{},
		},
		{
			name:  "nothing",
			key:   "InvalidKey",
			value: "",
			want:  Config{},
		},
		{
			name:  "string",
			key:   "APP_NAME_TEST",
			value: "MyApp",
			want:  Config{AppName: "MyApp"},
		},
		{
			name:  "int",
			key:   "PORT_TEST",
			value: "8080",
			want:  Config{Port: 8080},
		},
		{
			name:  "bool",
			key:   "DEBUG_TEST",
			value: "true",
			want:  Config{Debug: true},
		},
		{
			name:  "uint",
			key:   "MAX_RETRIES_TEST",
			value: "5",
			want:  Config{MaxRetries: 5},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			os.Setenv(tt.key, tt.value)
			config := Config{}
			err := env.FillStruct(context.Background(), &config)
			if err != nil {
				t.Errorf("FillStruct() error = %v", err)
			}
		})

	}

}

// TestGetEnv tests the GetEnv function manually
func TestGetEnv(t *testing.T) {
	os.Setenv("EXISTING_KEY", "EXISTING_VALUE")

	tests := []struct {
		name         string
		key          string
		defaultValue string
		expected     string
	}{
		{"existing key", "EXISTING_KEY", "DEFAULT_VALUE", "EXISTING_VALUE"},
		{"nin existing key", "NON_EXISTING_KEY", "DEFAULT_VALUE", "DEFAULT_VALUE"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			result := env.GetEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("For key '%s', expected '%s', got '%s'", tt.key, tt.expected, result)
			}
		})
	}
}
