// client.go
package ollamago

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

// Client represents an Ollama API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    http.Header
}

// Option is a function that configures the client
type Option func(*Client)

// NewClient creates a new Ollama client with the given options
func NewClient(options ...Option) *Client {
	c := &Client{
		baseURL: parseHost(os.Getenv("OLLAMA_HOST")),
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		headers: make(http.Header),
	}

	// Set default headers
	c.headers.Set("Content-Type", "application/json")
	c.headers.Set("Accept", "application/json")
	c.headers.Set("User-Agent", fmt.Sprintf("ollama-go/%s (%s %s) Go/%s",
		Version, runtime.GOOS, runtime.GOARCH, runtime.Version()))

	// Apply options
	for _, opt := range options {
		opt(c)
	}

	return c
}

// WithBaseURL sets a custom base URL for the client
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = parseHost(baseURL)
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithHeader adds a custom header to the client
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.headers.Set(key, value)
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// request makes an HTTP request to the Ollama API
func (c *Client) request(ctx context.Context, method, path string, body interface{}, response interface{}, stream bool) error {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	// Add headers
	for key, values := range c.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading error response: %w", err)
		}
		fmt.Println(string(bodyBytes))
		// Try to parse error response as JSON
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Error != "" {
			return &ResponseError{
				StatusCode: resp.StatusCode,
				Message:    errResp.Error,
			}
		}

		return &ResponseError{
			StatusCode: resp.StatusCode,
			Message:    string(bodyBytes),
		}
	}

	if response == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}

// requestStream makes a streaming HTTP request to the Ollama API
func (c *Client) requestStream(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Add headers
	for key, values := range c.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("reading error response: %w", err)
		}

		// Try to parse error response as JSON
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &errResp); err == nil && errResp.Error != "" {
			return nil, &ResponseError{
				StatusCode: resp.StatusCode,
				Message:    errResp.Error,
			}
		}

		return nil, &ResponseError{
			StatusCode: resp.StatusCode,
			Message:    string(bodyBytes),
		}
	}

	// Check if response is JSON or NDJSON stream
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "application/x-ndjson") {
		return nil, fmt.Errorf("unexpected content type: %s", contentType)
	}

	return resp, nil
}

// parseHost parses and validates the host URL
func parseHost(host string) string {
	if host == "" {
		return "http://127.0.0.1:11434"
	}

	// Add scheme if missing
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "http://" + host
	}

	u, err := url.Parse(host)
	if err != nil {
		return "http://127.0.0.1:11434"
	}

	// Add port if missing
	if u.Port() == "" {
		switch u.Scheme {
		case "https":
			host = fmt.Sprintf("%s:443", host)
		case "http":
			if !strings.Contains(host[7:], ":") {
				host = fmt.Sprintf("%s:11434", host)
			}
		}
	}

	// Ensure no trailing slash
	return strings.TrimSuffix(host, "/")
}
