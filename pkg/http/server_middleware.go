package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

// ServerMiddleware provides HTTP middleware for Go web servers
type ServerMiddleware struct {
	config *config.Config
	client Client
	logger Logger
}

// NewServerMiddleware creates a new HTTP middleware for Go web servers
func NewServerMiddleware(cfg *config.Config) *ServerMiddleware {
	if cfg == nil {
		cfg = config.New()
	}

	return &ServerMiddleware{
		config: cfg,
		client: NewClient(cfg),
		logger: nil,
	}
}

// WithLogger sets the logger for the middleware
func (m *ServerMiddleware) WithLogger(logger Logger) *ServerMiddleware {
	m.logger = logger
	return m
}

// WithClient sets the HTTP client for the middleware
func (m *ServerMiddleware) WithClient(client Client) *ServerMiddleware {
	m.client = client
	return m
}

// Handler returns an http.Handler that processes requests through the Complyance API
func (m *ServerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), m.config.Timeout)
		defer cancel()

		// Process the request through the middleware
		if err := m.processRequest(ctx, w, r); err != nil {
			m.handleError(w, err)
			return
		}

		// Call the next handler if the request wasn't handled
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// HandlerFunc returns an http.HandlerFunc that processes requests through the Complyance API
func (m *ServerMiddleware) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), m.config.Timeout)
		defer cancel()

		// Process the request through the middleware
		if err := m.processRequest(ctx, w, r); err != nil {
			m.handleError(w, err)
			return
		}

		// Call the next handler if the request wasn't handled
		if next != nil {
			next(w, r)
		}
	}
}

// ProcessInvoice is a middleware that processes invoice data from the request
func (m *ServerMiddleware) ProcessInvoice(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), m.config.Timeout)
		defer cancel()

		// Get country from request (query param, header, or path)
		country := r.URL.Query().Get("country")
		if country == "" {
			country = r.Header.Get("X-Country")
		}

		// Validate country
		if country == "" {
			m.handleError(w, errors.NewValidationError("country is required", nil).
				WithSuggestion("Provide country as query parameter or X-Country header"))
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			m.handleError(w, errors.NewValidationError("failed to read request body", err))
			return
		}
		defer r.Body.Close()

		// Parse request body as JSON
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			m.handleError(w, errors.NewValidationError("invalid JSON payload", err))
			return
		}

		// Get source from config
		source := m.getDefaultSource()
		if source == nil {
			m.handleError(w, errors.NewConfigError("no default source configured", nil).
				WithSuggestion("Configure the SDK with at least one source"))
			return
		}

		// Create request
		request := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, country)
		request.WithOperation(models.OperationSingle)
		request.WithMode(models.ModeDocuments)
		request.WithPurpose(models.PurposeInvoicing)
		request.WithPayload(payload)

		// Process request
		service := NewService(m.config).WithClient(m.client)
		response, err := service.PushToUnify(ctx, request)
		if err != nil {
			m.handleError(w, err)
			return
		}

		// Store response in context
		ctx = context.WithValue(r.Context(), contextKeyResponse, response)
		r = r.WithContext(ctx)

		// Call the next handler
		if next != nil {
			next.ServeHTTP(w, r)
		} else {
			// If no next handler, write response directly
			m.writeResponse(w, response)
		}
	})
}

// contextKey is a type for context keys
type contextKey string

// Context keys
const (
	contextKeyResponse contextKey = "complyance_response"
)

// GetResponse retrieves the Complyance response from the request context
func GetResponse(r *http.Request) *models.UnifyResponse {
	if value := r.Context().Value(contextKeyResponse); value != nil {
		if response, ok := value.(*models.UnifyResponse); ok {
			return response
		}
	}
	return nil
}

// processRequest processes the request through the middleware
func (m *ServerMiddleware) processRequest(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Check if this is a request we should handle
	if !m.shouldProcess(r) {
		return nil
	}

	// Log request if logger is available
	if m.logger != nil {
		m.logger.Info("Processing request through Complyance middleware", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
		})
	}

	// TODO: Implement request processing logic
	// This would depend on specific requirements for how the middleware should
	// interact with the Complyance API

	return nil
}

// shouldProcess determines if the request should be processed by this middleware
func (m *ServerMiddleware) shouldProcess(r *http.Request) bool {
	// This is a placeholder implementation
	// In a real implementation, this would check if the request matches
	// certain criteria (e.g., path prefix, content type, etc.)
	return false
}

// handleError writes an error response
func (m *ServerMiddleware) handleError(w http.ResponseWriter, err error) {
	// Log error if logger is available
	if m.logger != nil {
		m.logger.Error("Complyance middleware error", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Determine status code and error details
	statusCode := http.StatusInternalServerError
	errorCode := "INTERNAL_ERROR"
	message := err.Error()
	suggestion := ""

	// Extract error details based on error type
	switch e := err.(type) {
	case *errors.SDKError:
		message = e.Error()
		suggestion = e.Suggestion
		switch e.Code {
		case models.ErrorCodeValidationError:
			statusCode = http.StatusBadRequest
			errorCode = "VALIDATION_ERROR"
		case models.ErrorCodeAuthenticationError:
			statusCode = http.StatusUnauthorized
			errorCode = "AUTH_ERROR"
		case models.ErrorCodeRateLimitError:
			statusCode = http.StatusTooManyRequests
			errorCode = "RATE_LIMIT_ERROR"
		case models.ErrorCodeServerError:
			statusCode = http.StatusServiceUnavailable
			errorCode = "SERVER_ERROR"
		case models.ErrorCodeNetworkError:
			statusCode = http.StatusBadGateway
			errorCode = "NETWORK_ERROR"
		case models.ErrorCodeAPIError:
			statusCode = http.StatusInternalServerError
			errorCode = "API_ERROR"
		default:
			statusCode = http.StatusInternalServerError
			errorCode = "INTERNAL_ERROR"
		}
	default:
		// Handle other error types
		statusCode = http.StatusInternalServerError
		errorCode = "INTERNAL_ERROR"
		message = err.Error()
	}

	// Create error response
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":       errorCode,
			"message":    message,
			"suggestion": suggestion,
			"timestamp":  time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Write response
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

// writeResponse writes a success response
func (m *ServerMiddleware) writeResponse(w http.ResponseWriter, response *models.UnifyResponse) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Determine status code
	statusCode := http.StatusOK
	if response.Status == "error" {
		statusCode = http.StatusBadRequest
	}

	// Write response
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// getDefaultSource returns the default source from the configuration
func (m *ServerMiddleware) getDefaultSource() *models.Source {
	if len(m.config.Sources) > 0 {
		return m.config.Sources[0]
	}
	return nil
}