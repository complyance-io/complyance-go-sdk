package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
)

// MiddlewareFunc is a function that processes a request before it's sent
type MiddlewareFunc func(*Request) (*Request, error)

// AuthMiddleware adds authentication headers to requests
type AuthMiddleware struct {
	apiKey string
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(apiKey string) MiddlewareFunc {
	middleware := &AuthMiddleware{
		apiKey: apiKey,
	}
	return middleware.Process
}

// Process adds authentication headers to the request
func (m *AuthMiddleware) Process(req *Request) (*Request, error) {
	if m.apiKey == "" {
		return nil, errors.NewAuthError("API key is required", nil).
			WithSuggestion("Configure the SDK with a valid API key")
	}

	// Add API key header
	req.Headers["X-API-Key"] = m.apiKey

	// Add timestamp for request signing
	timestamp := time.Now().UTC().Format(time.RFC3339)
	req.Headers["X-Request-Timestamp"] = timestamp

	// Generate request signature
	signature, err := m.generateSignature(req, timestamp)
	if err != nil {
		return nil, errors.NewAuthError("failed to generate request signature", err)
	}
	req.Headers["X-Request-Signature"] = signature

	return req, nil
}

// generateSignature creates an HMAC signature for the request
func (m *AuthMiddleware) generateSignature(req *Request, timestamp string) (string, error) {
	// Create signature string from request details
	// Format: METHOD:PATH:TIMESTAMP
	signatureString := fmt.Sprintf("%s:%s:%s", req.Method, req.URL, timestamp)

	// Create HMAC-SHA256 signature using API key as secret
	h := hmac.New(sha256.New, []byte(m.apiKey))
	_, err := h.Write([]byte(signatureString))
	if err != nil {
		return "", err
	}

	// Return hex-encoded signature
	return hex.EncodeToString(h.Sum(nil)), nil
}

// LoggingMiddleware logs request and response details
type LoggingMiddleware struct {
	logger Logger
}

// Logger is an interface for logging
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger Logger) MiddlewareFunc {
	middleware := &LoggingMiddleware{
		logger: logger,
	}
	return middleware.Process
}

// Process logs request details
func (m *LoggingMiddleware) Process(req *Request) (*Request, error) {
	if m.logger == nil {
		return req, nil
	}

	// Log request details
	m.logger.Debug("Sending HTTP request", map[string]interface{}{
		"method": req.Method,
		"url":    req.URL,
	})

	return req, nil
}

// HeaderMiddleware adds custom headers to requests
type HeaderMiddleware struct {
	headers map[string]string
}

// NewHeaderMiddleware creates a new header middleware
func NewHeaderMiddleware(headers map[string]string) MiddlewareFunc {
	middleware := &HeaderMiddleware{
		headers: headers,
	}
	return middleware.Process
}

// Process adds headers to the request
func (m *HeaderMiddleware) Process(req *Request) (*Request, error) {
	for key, value := range m.headers {
		req.Headers[key] = value
	}
	return req, nil
}

// UserAgentMiddleware adds a User-Agent header to requests
func UserAgentMiddleware(userAgent string) MiddlewareFunc {
	return func(req *Request) (*Request, error) {
		req.Headers["User-Agent"] = userAgent
		return req, nil
	}
}