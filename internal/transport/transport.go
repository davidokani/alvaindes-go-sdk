package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Config contains the configuration required by the HTTP transport.
type Config struct {
	APIKey     string
	BaseURL    string
	Timeout    time.Duration
	HTTPClient *http.Client
}

// RequestOptions contains optional settings for an API request.
type RequestOptions struct {
	Headers http.Header
}

type AlvaindesError struct {
	Status   int
	Message  string
	Response any
	Code     string
}

func (e *AlvaindesError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

type Transport struct {
	config Config
	client *http.Client
}

func New(config Config) *Transport {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Transport{
		config: config,
		client: httpClient,
	}
}

// request executes an HTTP request against the Alvaindes API.
func (t *Transport) Request(
	ctx context.Context,
	method string,
	path string,
	body any,
	result any,
) error {
	return t.RequestWithOptions(ctx, method, path, body, result, RequestOptions{})
}

// RequestWithOptions executes an API request with optional custom headers.
func (t *Transport) RequestWithOptions(
	ctx context.Context,
	method string,
	path string,
	body any,
	result any,
	options RequestOptions,
) error {
	if !strings.HasPrefix(path, "/") || strings.HasPrefix(path, "//") {
		return fmt.Errorf(`API request paths must start with exactly one "/"`)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	allowedMethods := map[string]bool{
		http.MethodGet:    true,
		http.MethodPost:   true,
		http.MethodPatch:  true,
		http.MethodPut:    true,
		http.MethodDelete: true,
	}

	if !allowedMethods[method] {
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	var requestBody io.Reader

	if body != nil {
		encodedBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("unable to encode request body: %w", err)
		}

		requestBody = bytes.NewReader(encodedBody)
	}

	requestURL := t.config.BaseURL + path

	requestCtx, cancel := context.WithTimeout(ctx, t.config.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		requestCtx,
		method,
		requestURL,
		requestBody,
	)
	if err != nil {
		return fmt.Errorf("unable to create API request: %w", err)
	}

	// Authentication is controlled by the SDK.
	// Apply caller headers first. SDK-managed headers are set afterwards.
	for name, values := range options.Headers {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", t.config.APIKey)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response, err := t.client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return &AlvaindesError{
				Status:  0,
				Code:    "REQUEST_ABORTED",
				Message: "The request was aborted.",
			}
		}

		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) {
			return &AlvaindesError{
				Status: 0,
				Code:   "TIMEOUT",
				Message: fmt.Sprintf(
					"The request timed out after %dms.",
					t.config.Timeout.Milliseconds(),
				),
			}
		}

		return &AlvaindesError{
			Status:  0,
			Code:    "NETWORK_ERROR",
			Message: "Unable to complete the request to the Alvaindes API.",
		}
	}
	defer response.Body.Close()

	// Read the complete response while the request timeout is active.
	responseBytes, err := io.ReadAll(response.Body)
	if err != nil {
		if ctx.Err() != nil {
			return &AlvaindesError{
				Status:  0,
				Code:    "REQUEST_ABORTED",
				Message: "The request was aborted.",
			}
		}

		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) {
			return &AlvaindesError{
				Status: 0,
				Code:   "TIMEOUT",
				Message: fmt.Sprintf(
					"The request timed out after %dms.",
					t.config.Timeout.Milliseconds(),
				),
			}
		}

		return &AlvaindesError{
			Status:  0,
			Code:    "NETWORK_ERROR",
			Message: "Unable to complete the request to the Alvaindes API.",
		}
	}

	responseText := string(responseBytes)
	parsedResponse := parseHTTPResponse(responseText)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return createAPIError(response.StatusCode, parsedResponse)
	}

	if len(responseBytes) == 0 || result == nil {
		return nil
	}

	// Preserve non-JSON successful responses when the caller expects text.
	if textResult, ok := result.(*string); ok {
		if parsedText, isString := parsedResponse.(string); isString {
			*textResult = parsedText
			return nil
		}
	}

	if err := json.Unmarshal(responseBytes, result); err != nil {
		return fmt.Errorf("unable to decode Alvaindes API response: %w", err)
	}

	return nil
}

// parseHTTPResponse attempts to decode JSON.
// If the response is not JSON, the original text is preserved.
func parseHTTPResponse(responseText string) any {
	if responseText == "" {
		return nil
	}

	var parsed any

	if err := json.Unmarshal([]byte(responseText), &parsed); err != nil {
		return responseText
	}

	return parsed
}

// createAPIError converts an unsuccessful HTTP response into AlvaindesError.
func createAPIError(status int, response any) *AlvaindesError {
	message := fmt.Sprintf(
		"Alvaindes API request failed with status %d.",
		status,
	)

	var code string

	body, ok := response.(map[string]any)
	if ok {
		switch responseMessage := body["message"].(type) {
		case string:
			message = responseMessage

		case []any:
			messages := make([]string, 0, len(responseMessage))

			allStrings := true

			for _, item := range responseMessage {
				value, isString := item.(string)
				if !isString {
					allStrings = false
					break
				}

				messages = append(messages, value)
			}

			if allStrings {
				message = strings.Join(messages, "; ")
			}
		}

		if responseCode, ok := body["code"].(string); ok {
			code = responseCode
		}
	}

	return &AlvaindesError{
		Status:   status,
		Message:  message,
		Response: response,
		Code:     code,
	}
}
