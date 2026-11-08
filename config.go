package alvaindes

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"net/http"
)

const DefaultTimeout = 30 * time.Second

// Config contains the configuration required to create an Alvaindes client.
type Config struct {
	APIKey     string
	BaseURL    string
	Timeout    time.Duration
	HTTPClient *http.Client
}

type resolvedConfig struct {
	apiKey     string
	baseURL    string
	timeout    time.Duration
	httpClient *http.Client
}

func resolveConfig(config Config) (resolvedConfig, error) {
	apiKey := strings.TrimSpace(config.APIKey)
	if apiKey == "" {
		return resolvedConfig{}, fmt.Errorf("alvaindes API key is required")
	}

	baseURL := strings.TrimSpace(config.BaseURL)
	if baseURL == "" {
		return resolvedConfig{}, fmt.Errorf("alvaindes base URL is required")
	}

	parsedURL, err := url.Parse(baseURL)
	if err != nil || parsedURL.Host == "" {
		return resolvedConfig{}, fmt.Errorf("alvaindes base URL must be a valid URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return resolvedConfig{}, fmt.Errorf("alvaindes base URL must use HTTP or HTTPS")
	}

	if parsedURL.User != nil {
		return resolvedConfig{}, fmt.Errorf("alvaindes base URL must not contain credentials")
	}

	if parsedURL.RawQuery != "" || parsedURL.ForceQuery || parsedURL.Fragment != "" {
		return resolvedConfig{}, fmt.Errorf("alvaindes base URL must not contain a query string or fragment")
	}

	if parsedURL.Path != "" && parsedURL.Path != "/" {
		return resolvedConfig{}, fmt.Errorf("alvaindes base URL must be an API origin without a path prefix")
	}

	if parsedURL.RawPath != "" {
		return resolvedConfig{}, fmt.Errorf("alvaindes base URL must not contain an encoded path")
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	if timeout < time.Millisecond || timeout%time.Millisecond != 0 {
		return resolvedConfig{}, fmt.Errorf("alvaindes timeout must be a positive whole number of milliseconds")
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	origin := parsedURL.Scheme + "://" + parsedURL.Host

	return resolvedConfig{
		apiKey:     apiKey,
		baseURL:    origin,
		timeout:    timeout,
		httpClient: httpClient,
	}, nil
}
