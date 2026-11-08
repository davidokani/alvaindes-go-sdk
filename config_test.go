package alvaindes

import (
	"net/http"
	"testing"
	"time"
)

func TestResolveConfigUsesDefaultsAndNormalizesValues(t *testing.T) {
	config, err := resolveConfig(Config{
		APIKey:  " alva_test_example ",
		BaseURL: "https://api.alvaindes.com/",
	})
	if err != nil {
		t.Fatalf("resolveConfig() returned an unexpected error: %v", err)
	}

	if config.apiKey != "alva_test_example" {
		t.Errorf("expected trimmed API key, got %q", config.apiKey)
	}

	if config.baseURL != "https://api.alvaindes.com" {
		t.Errorf("expected normalized base URL, got %q", config.baseURL)
	}

	if config.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, config.timeout)
	}

	if config.httpClient == nil {
		t.Error("expected a default HTTP client")
	}
}

func TestResolveConfigSupportsCustomTimeoutAndHTTPClient(t *testing.T) {
	httpClient := &http.Client{}
	timeout := 5 * time.Second

	config, err := resolveConfig(Config{
		APIKey:     "alva_test_example",
		BaseURL:    "https://api.alvaindes.com",
		Timeout:    timeout,
		HTTPClient: httpClient,
	})
	if err != nil {
		t.Fatalf("resolveConfig() returned an unexpected error: %v", err)
	}

	if config.timeout != timeout {
		t.Errorf("expected timeout %v, got %v", timeout, config.timeout)
	}

	if config.httpClient != httpClient {
		t.Error("expected the injected HTTP client to be preserved")
	}
}

func TestResolveConfigRejectsInvalidConfigurations(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "empty API key",
			config: Config{
				BaseURL: "https://api.alvaindes.com",
			},
		},
		{
			name: "whitespace API key",
			config: Config{
				APIKey:  " ",
				BaseURL: "https://api.alvaindes.com",
			},
		},
		{
			name: "empty base URL",
			config: Config{
				APIKey: "alva_test_example",
			},
		},
		{
			name: "malformed base URL",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "not-a-url",
			},
		},
		{
			name: "missing host",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https:///api",
			},
		},
		{
			name: "unsupported protocol",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "ftp://api.alvaindes.com",
			},
		},
		{
			name: "credentials in URL",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://user:password@api.alvaindes.com",
			},
		},
		{
			name: "query string",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://api.alvaindes.com?version=1",
			},
		},
		{
			name: "empty query string",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://api.alvaindes.com?",
			},
		},
		{
			name: "fragment",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://api.alvaindes.com#section",
			},
		},
		{
			name: "path prefix",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://api.alvaindes.com/v1",
			},
		},
		{
			name: "encoded path",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://api.alvaindes.com/%2F",
			},
		},
		{
			name: "negative timeout",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://api.alvaindes.com",
				Timeout: -time.Millisecond,
			},
		},
		{
			name: "sub-millisecond timeout",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://api.alvaindes.com",
				Timeout: 500 * time.Microsecond,
			},
		},
		{
			name: "fractional-millisecond timeout",
			config: Config{
				APIKey:  "alva_test_example",
				BaseURL: "https://api.alvaindes.com",
				Timeout: 1500 * time.Microsecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resolveConfig(tt.config)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
	}
}
