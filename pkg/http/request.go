package http

import (
	"io"
	"net/http"
)

// Request represents an HTTP request
type Request struct {
	// Method is the HTTP method (GET, POST, etc.)
	Method string

	// URL is the request URL
	URL string

	// Headers are the HTTP headers
	Headers map[string]string

	// Body is the request body
	Body io.Reader

	// Context data for middleware
	ContextData map[string]interface{}
	
	// retryMiddleware is the retry middleware for this request
	retryMiddleware *RetryMiddleware
}

// NewRequest creates a new HTTP request
func NewRequest(method, url string, body io.Reader, headers map[string]string) *Request {
	if headers == nil {
		headers = make(map[string]string)
	}

	return &Request{
		Method:      method,
		URL:         url,
		Headers:     headers,
		Body:        body,
		ContextData: make(map[string]interface{}),
	}
}

// WithHeader adds a header to the request
func (r *Request) WithHeader(key, value string) *Request {
	r.Headers[key] = value
	return r
}

// WithHeaders adds multiple headers to the request
func (r *Request) WithHeaders(headers map[string]string) *Request {
	for key, value := range headers {
		r.Headers[key] = value
	}
	return r
}

// WithContextData adds context data to the request
func (r *Request) WithContextData(key string, value interface{}) *Request {
	r.ContextData[key] = value
	return r
}

// GetContextData retrieves context data from the request
func (r *Request) GetContextData(key string) (interface{}, bool) {
	value, ok := r.ContextData[key]
	return value, ok
}

// Clone creates a copy of the request
func (r *Request) Clone() *Request {
	headers := make(map[string]string)
	for k, v := range r.Headers {
		headers[k] = v
	}

	contextData := make(map[string]interface{})
	for k, v := range r.ContextData {
		contextData[k] = v
	}

	return &Request{
		Method:      r.Method,
		URL:         r.URL,
		Headers:     headers,
		Body:        r.Body,
		ContextData: contextData,
	}
}

// IsGet returns true if the request method is GET
func (r *Request) IsGet() bool {
	return r.Method == http.MethodGet
}

// IsPost returns true if the request method is POST
func (r *Request) IsPost() bool {
	return r.Method == http.MethodPost
}

// IsPut returns true if the request method is PUT
func (r *Request) IsPut() bool {
	return r.Method == http.MethodPut
}

// IsDelete returns true if the request method is DELETE
func (r *Request) IsDelete() bool {
	return r.Method == http.MethodDelete
}