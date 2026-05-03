package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

// Default configuration values
const (
	DefaultTimeout      = 30 * time.Second
	DefaultMaxRetries   = 3
	DefaultBaseDelay    = 500 * time.Millisecond
	DefaultMaxDelay     = 5 * time.Second
	DefaultJitterFactor = 0.1
)

// Environment variable names
const (
	EnvAPIKey      = "COMPLYANCE_API_KEY"
	EnvEnvironment = "COMPLYANCE_ENVIRONMENT"
	EnvBaseURL     = "COMPLYANCE_BASE_URL"
	EnvMaxRetries  = "COMPLYANCE_MAX_RETRIES"
	EnvTimeout     = "COMPLYANCE_TIMEOUT"
)

// Config holds the SDK configuration
type Config struct {
	// APIKey is the authentication key for the Complyance API
	APIKey string

	// Environment specifies the target environment (sandbox, production)
	Environment models.Environment

	// BaseURL overrides the default API URL for the specified environment
	BaseURL string

	// Timeout specifies the HTTP request timeout
	Timeout time.Duration

	// Sources is a list of source systems for document processing
	Sources []*models.Source

	// RetryConfig holds the retry and circuit breaker configuration
	RetryConfig *RetryConfig
}

// RetryConfig holds retry and circuit breaker settings
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// BaseDelay is the initial delay between retries
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration

	// JitterFactor adds randomness to the retry delay to prevent thundering herd
	JitterFactor float64

	// CircuitBreakerEnabled enables the circuit breaker pattern
	CircuitBreakerEnabled bool

	// FailureThreshold is the number of failures before opening the circuit
	FailureThreshold int

	// CircuitBreakerTimeout is the duration to keep the circuit open
	CircuitBreakerTimeout time.Duration

	// RetryableHTTPCodes is a list of HTTP status codes that should trigger a retry
	RetryableHTTPCodes []int
}

// Option is a function that configures the Config
type Option func(*Config)

// New creates a new Config with the provided options
func New(options ...Option) *Config {
	cfg := &Config{
		Environment: models.EnvironmentSandbox,
		Timeout:     DefaultTimeout,
		RetryConfig: &RetryConfig{
			MaxRetries:           DefaultMaxRetries,
			BaseDelay:            DefaultBaseDelay,
			MaxDelay:             DefaultMaxDelay,
			JitterFactor:         DefaultJitterFactor,
			CircuitBreakerEnabled: true,
			FailureThreshold:     5,
			CircuitBreakerTimeout: 60 * time.Second,
			RetryableHTTPCodes:   []int{408, 429, 500, 502, 503, 504},
		},
	}

	// Apply options
	for _, option := range options {
		option(cfg)
	}

	return cfg
}

// FromEnv creates a new Config from environment variables
func FromEnv() *Config {
	cfg := New()

	if apiKey := os.Getenv(EnvAPIKey); apiKey != "" {
		cfg.APIKey = apiKey
	}

	if env := os.Getenv(EnvEnvironment); env != "" {
		switch strings.ToLower(env) {
		case "sandbox":
			cfg.Environment = models.EnvironmentSandbox
		case "production":
			cfg.Environment = models.EnvironmentProduction
		case "local":
			cfg.Environment = models.EnvironmentLocal
		}
	}

	if baseURL := os.Getenv(EnvBaseURL); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	if maxRetries := os.Getenv(EnvMaxRetries); maxRetries != "" {
		var mr int
		if _, err := fmt.Sscanf(maxRetries, "%d", &mr); err == nil && mr >= 0 {
			cfg.RetryConfig.MaxRetries = mr
		}
	}

	if timeout := os.Getenv(EnvTimeout); timeout != "" {
		var t int
		if _, err := fmt.Sscanf(timeout, "%d", &t); err == nil && t > 0 {
			cfg.Timeout = time.Duration(t) * time.Second
		}
	}

	return cfg
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return errors.New("API key is required")
	}

	if c.Environment == "" {
		return errors.New("environment is required")
	}

	if c.Timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	if c.RetryConfig == nil {
		return errors.New("retry configuration is required")
	}

	if c.RetryConfig.MaxRetries < 0 {
		return errors.New("max retries must be greater than or equal to 0")
	}

	if c.RetryConfig.BaseDelay <= 0 {
		return errors.New("base delay must be greater than 0")
	}

	if c.RetryConfig.MaxDelay < c.RetryConfig.BaseDelay {
		return errors.New("max delay must be greater than or equal to base delay")
	}

	if c.RetryConfig.JitterFactor < 0 || c.RetryConfig.JitterFactor > 1 {
		return errors.New("jitter factor must be between 0 and 1")
	}

	return nil
}

// GetBaseURL returns the appropriate base URL for the configured environment
func (c *Config) GetBaseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}

	switch c.Environment {
	case models.EnvironmentProduction:
		return "https://api.complyance.io/v1"
	case models.EnvironmentSandbox:
		return "https://api.sandbox.complyance.io/v1"
	case models.EnvironmentLocal:
		return "http://localhost:8080/v1"
	default:
		return "https://api.sandbox.complyance.io/v1"
	}
}

// WithAPIKey sets the API key
func WithAPIKey(apiKey string) Option {
	return func(c *Config) {
		c.APIKey = apiKey
	}
}
//}
// WithEnvironment sets the environment
func WithEnvironment(env models.Environment) Option {
	return func(c *Config) {
		c.Environment = env
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) Option {
	return func(c *Config) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets the HTTP request timeout
func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithSource adds a source to the configuration
func WithSource(source *models.Source) Option {
	return func(c *Config) {
		c.Sources = append(c.Sources, source)
	}
}

// WithSources sets the sources list
func WithSources(sources []*models.Source) Option {
	return func(c *Config) {
		c.Sources = sources
	}
}

// WithRetryConfig sets the retry configuration
func WithRetryConfig(retryConfig *RetryConfig) Option {
	return func(c *Config) {
		c.RetryConfig = retryConfig
	}
}

// AggressiveRetryConfig returns a retry configuration optimized for high availability
func AggressiveRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:           7,
		BaseDelay:            200 * time.Millisecond,
		MaxDelay:             2 * time.Second,
		JitterFactor:         0.1,
		CircuitBreakerEnabled: true,
		FailureThreshold:     10,
		CircuitBreakerTimeout: 30 * time.Second,
		RetryableHTTPCodes:   []int{408, 429, 500, 502, 503, 504},
	}
}

// ConservativeRetryConfig returns a retry configuration optimized for production safety
func ConservativeRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:           3,
		BaseDelay:            1 * time.Second,
		MaxDelay:             10 * time.Second,
		JitterFactor:         0.1,
		CircuitBreakerEnabled: true,
		FailureThreshold:     5,
		CircuitBreakerTimeout: 60 * time.Second,
		RetryableHTTPCodes:   []int{408, 429, 500, 502, 503, 504},
	}
}

// NoRetryConfig returns a retry configuration with retries disabled
func NoRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:           0,
		BaseDelay:            0,
		MaxDelay:             0,
		JitterFactor:         0,
		CircuitBreakerEnabled: false,
		FailureThreshold:     0,
		CircuitBreakerTimeout: 0,
		RetryableHTTPCodes:   []int{},
	}
}