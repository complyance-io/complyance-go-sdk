package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
)

// Response represents an HTTP response
type Response struct {
	// StatusCode is the HTTP status code
	StatusCode int

	// Headers are the HTTP headers
	Headers http.Header

	// Body is the response body
	Body []byte

	// Duration is the request duration
	Duration time.Duration
}

// IsSuccess returns true if the response status code is in the 2xx range
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if the response status code is in the 4xx or 5xx range
func (r *Response) IsError() bool {
	return r.StatusCode >= 400
}

// IsClientError returns true if the response status code is in the 4xx range
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the response status code is in the 5xx range
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500
}

// JSON unmarshals the response body into the provided value
func (r *Response) JSON(v interface{}) error {
	if len(r.Body) == 0 {
		return errors.NewAPIError("empty response body", nil)
	}

	if err := json.Unmarshal(r.Body, v); err != nil {
		return errors.NewAPIError("failed to unmarshal response", err).
			AddContext("body", string(r.Body))
	}

	return nil
}

// String returns the response body as a string
func (r *Response) String() string {
	return string(r.Body)
}

// GetHeader returns the value of the specified header
func (r *Response) GetHeader(key string) string {
	return r.Headers.Get(key)
}

// GetContentType returns the Content-Type header
func (r *Response) GetContentType() string {
	return r.GetHeader("Content-Type")
}

// IsJSON returns true if the Content-Type header indicates JSON
func (r *Response) IsJSON() bool {
	contentType := r.GetContentType()
	return contentType == "application/json" || contentType == "application/json; charset=utf-8"
}