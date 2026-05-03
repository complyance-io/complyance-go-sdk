package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
)

// Client is the interface for making HTTP requests
type Client interface {
	// Do executes an HTTP request with the provided context
	Do(ctx context.Context, req *Request) (*Response, error)

	// Get performs an HTTP GET request
	Get(ctx context.Context, path string, headers map[string]string) (*Response, error)

	// Post performs an HTTP POST request with a JSON body
	Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error)

	// Put performs an HTTP PUT request with a JSON body
	Put(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error)

	// Delete performs an HTTP DELETE request
	Delete(ctx context.Context, path string, headers map[string]string) (*Response, error)
}

// DefaultClient is the default implementation of the Client interface
type DefaultClient struct {
	// httpClient is the underlying HTTP client
	httpClient *http.Client

	// baseURL is the base URL for all requests
	baseURL string

	// config is the SDK configuration
	config *config.Config

	// middleware is a list of middleware functions to apply to requests
	middleware []MiddlewareFunc
}

// NewClient creates a new HTTP client with the provided configuration
func NewClient(cfg *config.Config) Client {
	if cfg == nil {
		cfg = config.New()
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: cfg.Timeout,
	}

	// Create default client
	client := &DefaultClient{
		httpClient: httpClient,
		baseURL:    cfg.GetBaseURL(),
		config:     cfg,
		middleware: []MiddlewareFunc{},
	}

	// Add authentication middleware
	client.Use(NewAuthMiddleware(cfg.APIKey))
	
	// Add retry middleware if enabled
	if cfg.RetryConfig != nil && cfg.RetryConfig.MaxRetries > 0 {
		client.Use(NewRetryMiddleware(cfg.RetryConfig).Process)
	}

	return client
}

// Use adds middleware to the client
func (c *DefaultClient) Use(middleware MiddlewareFunc) {
	c.middleware = append(c.middleware, middleware)
}

// Do executes an HTTP request with the provided context
func (c *DefaultClient) Do(ctx context.Context, req *Request) (*Response, error) {
	// Check if context is already canceled
	if ctx.Err() != nil {
		return nil, errors.NewNetworkError("context canceled", ctx.Err())
	}

	// Apply middleware to the request
	for _, middleware := range c.middleware {
		var err error
		req, err = middleware(req)
		if err != nil {
			return nil, err
		}
	}

	// If retry middleware is attached to the request, use it
	if req.retryMiddleware != nil {
		return req.retryMiddleware.DoWithRetry(ctx, c, req)
	}

	// Otherwise, perform a regular request
	return c.doRequest(ctx, req)
}

// doRequest performs the actual HTTP request without retry logic
func (c *DefaultClient) doRequest(ctx context.Context, req *Request) (*Response, error) {
	// Build the full URL
	url := req.URL
	if !isAbsoluteURL(url) {
		url = fmt.Sprintf("%s%s", c.baseURL, url)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, req.Body)
	if err != nil {
		return nil, errors.NewNetworkError("failed to create request", err)
	}

	// Add headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	startTime := time.Now()
	httpResp, err := c.httpClient.Do(httpReq)
	duration := time.Since(startTime)

	// Handle network errors
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, errors.NewNetworkError("request timed out", err).
				WithSuggestion("Consider increasing the timeout in the SDK configuration")
		}
		if ctx.Err() == context.Canceled {
			return nil, errors.NewNetworkError("request canceled", err)
		}
		return nil, errors.NewNetworkError("request failed", err)
	}

	// Create response
	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Header,
		Duration:   duration,
	}

	// Read response body
	defer httpResp.Body.Close()
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, errors.NewNetworkError("failed to read response body", err)
	}
	resp.Body = body

	// Handle error responses
	if resp.StatusCode >= 400 {
		return resp, c.handleErrorResponse(resp)
	}

	return resp, nil
}

// Get performs an HTTP GET request
func (c *DefaultClient) Get(ctx context.Context, path string, headers map[string]string) (*Response, error) {
	req := NewRequest(http.MethodGet, path, nil, headers)
	return c.Do(ctx, req)
}

// Post performs an HTTP POST request with a JSON body
func (c *DefaultClient) Post(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error) {
	// Marshal body to JSON
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, errors.NewValidationError("failed to marshal request body", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	// Create request
	req := NewRequest(http.MethodPost, path, bodyReader, headers)
	if body != nil {
		req.Headers["Content-Type"] = "application/json"
	}

	return c.Do(ctx, req)
}

// Put performs an HTTP PUT request with a JSON body
func (c *DefaultClient) Put(ctx context.Context, path string, body interface{}, headers map[string]string) (*Response, error) {
	// Marshal body to JSON
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, errors.NewValidationError("failed to marshal request body", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	// Create request
	req := NewRequest(http.MethodPut, path, bodyReader, headers)
	if body != nil {
		req.Headers["Content-Type"] = "application/json"
	}

	return c.Do(ctx, req)
}

// Delete performs an HTTP DELETE request
func (c *DefaultClient) Delete(ctx context.Context, path string, headers map[string]string) (*Response, error) {
	req := NewRequest(http.MethodDelete, path, nil, headers)
	return c.Do(ctx, req)
}

// handleErrorResponse converts HTTP error responses to SDK errors
func (c *DefaultClient) handleErrorResponse(resp *Response) error {
	// Try to parse error response as JSON
	var errorResponse struct {
		Error struct {
			Code    string                 `json:"code"`
			Message string                 `json:"message"`
			Context map[string]interface{} `json:"context,omitempty"`
		} `json:"error"`
	}

	if err := json.Unmarshal(resp.Body, &errorResponse); err == nil && errorResponse.Error.Message != "" {
		// Create appropriate error based on status code
		switch {
		case resp.StatusCode == http.StatusUnauthorized:
			return errors.NewAuthError(errorResponse.Error.Message, nil).
				WithContext(errorResponse.Error.Context).
				WithSuggestion("Check your API key and ensure it has the necessary permissions")

		case resp.StatusCode == http.StatusForbidden:
			return errors.NewAuthError(errorResponse.Error.Message, nil).
				WithContext(errorResponse.Error.Context).
				WithSuggestion("Your API key does not have permission to access this resource")

		case resp.StatusCode == http.StatusNotFound:
			return errors.NewAPIError(errorResponse.Error.Message, nil).
				WithContext(errorResponse.Error.Context).
				WithSuggestion("Check the request URL and ensure the resource exists")

		case resp.StatusCode == http.StatusTooManyRequests:
			return errors.NewRateLimitError(errorResponse.Error.Message, nil).
				WithContext(errorResponse.Error.Context).
				WithSuggestion("Reduce request frequency or contact support to increase your rate limit")

		case resp.StatusCode >= 500:
			return errors.NewServerError(errorResponse.Error.Message, nil).
				WithContext(errorResponse.Error.Context).
				WithSuggestion("This is a server-side issue. Please try again later or contact support")

		default:
			return errors.NewAPIError(errorResponse.Error.Message, nil).
				WithContext(errorResponse.Error.Context)
		}
	}

	// Fallback for non-JSON error responses
	var message string
	if len(resp.Body) > 0 {
		message = string(resp.Body)
	} else {
		message = fmt.Sprintf("HTTP error %d", resp.StatusCode)
	}

	switch {
	case resp.StatusCode == http.StatusUnauthorized:
		return errors.NewAuthError(message, nil).
			WithSuggestion("Check your API key and ensure it has the necessary permissions")

	case resp.StatusCode == http.StatusForbidden:
		return errors.NewAuthError(message, nil).
			WithSuggestion("Your API key does not have permission to access this resource")

	case resp.StatusCode == http.StatusNotFound:
		return errors.NewAPIError(message, nil).
			WithSuggestion("Check the request URL and ensure the resource exists")

	case resp.StatusCode == http.StatusTooManyRequests:
		return errors.NewRateLimitError(message, nil).
			WithSuggestion("Reduce request frequency or contact support to increase your rate limit")

	case resp.StatusCode >= 500:
		return errors.NewServerError(message, nil).
			WithSuggestion("This is a server-side issue. Please try again later or contact support")

	default:
		return errors.NewAPIError(message, nil)
	}
}

// isAbsoluteURL checks if the URL is absolute
func isAbsoluteURL(url string) bool {
	return len(url) > 7 && (url[:7] == "http://" || url[:8] == "https://")
}