package tests

import (
	"os"
	"testing"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestConfigValidation(t *testing.T) {
	// Test valid configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	err := cfg.Validate()
	assert.NoError(t, err)

	// Test missing API key
	invalidCfg := config.New(
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	err = invalidCfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API key is required")

	// Test missing environment
	invalidCfg = config.New(
		config.WithAPIKey("test-api-key"),
	)
	invalidCfg.Environment = ""

	err = invalidCfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "environment is required")

	// Test invalid timeout
	invalidCfg = config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithTimeout(-1 * time.Second),
	)

	err = invalidCfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timeout must be greater than 0")
}

func TestConfigFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("COMPLYANCE_API_KEY", "env-api-key")
	os.Setenv("COMPLYANCE_ENVIRONMENT", "production")
	os.Setenv("COMPLYANCE_BASE_URL", "https://custom-api.complyance.io/v1")
	os.Setenv("COMPLYANCE_MAX_RETRIES", "5")
	os.Setenv("COMPLYANCE_TIMEOUT", "60")

	// Create config from environment
	cfg := config.FromEnv()

	// Verify values
	assert.Equal(t, "env-api-key", cfg.APIKey)
	assert.Equal(t, models.EnvironmentProduction, cfg.Environment)
	assert.Equal(t, "https://custom-api.complyance.io/v1", cfg.BaseURL)
	assert.Equal(t, 5, cfg.RetryConfig.MaxRetries)
	assert.Equal(t, 60*time.Second, cfg.Timeout)

	// Clean up
	os.Unsetenv("COMPLYANCE_API_KEY")
	os.Unsetenv("COMPLYANCE_ENVIRONMENT")
	os.Unsetenv("COMPLYANCE_BASE_URL")
	os.Unsetenv("COMPLYANCE_MAX_RETRIES")
	os.Unsetenv("COMPLYANCE_TIMEOUT")
}

func TestGetBaseURL(t *testing.T) {
	// Test custom base URL
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithBaseURL("https://custom-api.complyance.io/v1"),
	)

	assert.Equal(t, "https://custom-api.complyance.io/v1", cfg.GetBaseURL())

	// Test sandbox environment
	cfg = config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	assert.Equal(t, "https://api.sandbox.complyance.io/v1", cfg.GetBaseURL())

	// Test production environment
	cfg = config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentProduction),
	)

	assert.Equal(t, "https://api.complyance.io/v1", cfg.GetBaseURL())

	// Test local environment
	cfg = config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentLocal),
	)

	assert.Equal(t, "http://localhost:8080/v1", cfg.GetBaseURL())
}

func TestRetryConfigs(t *testing.T) {
	// Test aggressive retry config
	aggressive := config.AggressiveRetryConfig()
	assert.Equal(t, 7, aggressive.MaxRetries)
	assert.Equal(t, 200*time.Millisecond, aggressive.BaseDelay)
	assert.Equal(t, 2*time.Second, aggressive.MaxDelay)
	assert.True(t, aggressive.CircuitBreakerEnabled)

	// Test conservative retry config
	conservative := config.ConservativeRetryConfig()
	assert.Equal(t, 3, conservative.MaxRetries)
	assert.Equal(t, 1*time.Second, conservative.BaseDelay)
	assert.Equal(t, 10*time.Second, conservative.MaxDelay)
	assert.True(t, conservative.CircuitBreakerEnabled)

	// Test no retry config
	noRetry := config.NoRetryConfig()
	assert.Equal(t, 0, noRetry.MaxRetries)
	assert.False(t, noRetry.CircuitBreakerEnabled)
}