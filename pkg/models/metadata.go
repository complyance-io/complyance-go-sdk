package models

import (
	"fmt"
	"time"
)

// RequestMetadata contains additional request information
type RequestMetadata struct {
	// APIKey is the authentication key (automatically populated)
	APIKey string `json:"api_key,omitempty"`

	// RequestID is a unique identifier for the request
	RequestID string `json:"request_id,omitempty"`

	// Timestamp is the request creation time
	Timestamp string `json:"timestamp,omitempty"`

	// Environment is the target environment
	Environment string `json:"environment,omitempty"`

	// ClientInfo contains information about the client making the request
	ClientInfo *ClientInfo `json:"client_info,omitempty"`
}

// ResponseMetadata contains additional response information
type ResponseMetadata struct {
	// RequestID is the unique identifier for the request
	RequestID string `json:"request_id,omitempty"`

	// ProcessingTime is the time taken to process the request in milliseconds
	ProcessingTime int64 `json:"processing_time,omitempty"`

	// TraceID is the unique identifier for tracing the request through the system
	TraceID string `json:"trace_id,omitempty"`

	// Timestamp is the response creation time
	Timestamp string `json:"timestamp,omitempty"`

	// ServerInfo contains information about the server processing the request
	ServerInfo *ServerInfo `json:"server_info,omitempty"`
}

// ClientInfo contains information about the client making the request
type ClientInfo struct {
	// SDKVersion is the version of the SDK making the request
	SDKVersion string `json:"sdk_version,omitempty"`

	// SDKLanguage is the programming language of the SDK
	SDKLanguage string `json:"sdk_language,omitempty"`

	// OSName is the operating system name
	OSName string `json:"os_name,omitempty"`

	// OSVersion is the operating system version
	OSVersion string `json:"os_version,omitempty"`
}

// ServerInfo contains information about the server processing the request
type ServerInfo struct {
	// Version is the server version
	Version string `json:"version,omitempty"`

	// Region is the server region
	Region string `json:"region,omitempty"`

	// NodeID is the server node ID
	NodeID string `json:"node_id,omitempty"`
}

// NewRequestMetadata creates a new RequestMetadata with default values
func NewRequestMetadata() *RequestMetadata {
	return &RequestMetadata{
		RequestID:  generateRequestID(),
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		ClientInfo: NewClientInfo(),
	}
}

// NewResponseMetadata creates a new ResponseMetadata with default values
func NewResponseMetadata() *ResponseMetadata {
	return &ResponseMetadata{
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		ServerInfo: &ServerInfo{},
	}
}

// NewClientInfo creates a new ClientInfo with default values
func NewClientInfo() *ClientInfo {
	return &ClientInfo{
		SDKVersion:  "1.0.0",
		SDKLanguage: "go",
	}
}

// WithAPIKey sets the API key
func (m *RequestMetadata) WithAPIKey(apiKey string) *RequestMetadata {
	m.APIKey = apiKey
	return m
}

// WithEnvironment sets the environment
func (m *RequestMetadata) WithEnvironment(environment string) *RequestMetadata {
	m.Environment = environment
	return m
}

// WithClientInfo sets the client info
func (m *RequestMetadata) WithClientInfo(clientInfo *ClientInfo) *RequestMetadata {
	m.ClientInfo = clientInfo
	return m
}

// WithRequestID sets the request ID in the response metadata
func (m *ResponseMetadata) WithRequestID(requestID string) *ResponseMetadata {
	m.RequestID = requestID
	return m
}

// WithProcessingTime sets the processing time
func (m *ResponseMetadata) WithProcessingTime(processingTime int64) *ResponseMetadata {
	m.ProcessingTime = processingTime
	return m
}

// WithTraceID sets the trace ID
func (m *ResponseMetadata) WithTraceID(traceID string) *ResponseMetadata {
	m.TraceID = traceID
	return m
}

// WithServerInfo sets the server info
func (m *ResponseMetadata) WithServerInfo(serverInfo *ServerInfo) *ResponseMetadata {
	m.ServerInfo = serverInfo
	return m
}

// WithSDKVersion sets the SDK version
func (c *ClientInfo) WithSDKVersion(version string) *ClientInfo {
	c.SDKVersion = version
	return c
}

// WithOSInfo sets the OS information
func (c *ClientInfo) WithOSInfo(name string, version string) *ClientInfo {
	c.OSName = name
	c.OSVersion = version
	return c
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	now := time.Now().UTC()
	return fmt.Sprintf("req_%d%d%d", now.Unix(), now.Nanosecond(), time.Now().UnixNano()%1000)
}