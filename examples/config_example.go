package main

import (
	"fmt"
	"log"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

func main() {
	// Create a configuration using the functional options pattern
	cfg := config.New(
		config.WithAPIKey("your-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithTimeout(30 * time.Second),
	)

	// Add a source
	source := models.NewSource("erp-system", models.SourceTypeFirstParty, "My ERP System")
	source.WithVersion("1.0.0")
	source.AddMetadata("tenant_id", "tenant-123")

	cfg = config.New(
		config.WithAPIKey("your-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
		config.WithSource(source),
		config.WithRetryConfig(config.AggressiveRetryConfig()),
	)

	// Validate the configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Print the configuration
	fmt.Printf("API Key: %s\n", cfg.APIKey)
	fmt.Printf("Environment: %s\n", cfg.Environment)
	fmt.Printf("Base URL: %s\n", cfg.GetBaseURL())
	fmt.Printf("Timeout: %s\n", cfg.Timeout)
	fmt.Printf("Max Retries: %d\n", cfg.RetryConfig.MaxRetries)

	// Create a configuration from environment variables
	envCfg := config.FromEnv()
	fmt.Printf("Environment config API Key: %s\n", envCfg.APIKey)
	fmt.Printf("Environment config Environment: %s\n", envCfg.Environment)
}