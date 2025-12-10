package tests

import (
	"context"
	"testing"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSDKConfiguration(t *testing.T) {
	// Test valid configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	err := pkg.Configure(cfg)
	assert.NoError(t, err)

	// Test nil configuration
	err = pkg.Configure(nil)
	assert.Error(t, err)
	assert.True(t, errors.IsConfigError(err))

	// Test invalid configuration
	invalidCfg := config.New()
	err = pkg.Configure(invalidCfg)
	assert.Error(t, err)
	assert.True(t, errors.IsConfigError(err))
}

func TestNewClient(t *testing.T) {
	// Test valid configuration
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	client, err := pkg.NewClient(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	// Test nil configuration
	client, err = pkg.NewClient(nil)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.True(t, errors.IsConfigError(err))

	// Test invalid configuration
	invalidCfg := config.New()
	client, err = pkg.NewClient(invalidCfg)
	assert.Error(t, err)
	assert.Nil(t, client)
	assert.True(t, errors.IsConfigError(err))
}

func TestPushToUnify(t *testing.T) {
	// Configure the SDK
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	client, err := pkg.NewClient(cfg)
	assert.NoError(t, err)

	// Create a valid source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Create a valid request
	request := models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "SA")
	request.WithPayload(map[string]interface{}{
		"invoice_number": "INV-001",
	})

	// Test valid request
	ctx := context.Background()
	resp, err := client.PushToUnify(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "success", resp.Status)

	// Test nil request
	resp, err = client.PushToUnify(ctx, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test invalid request
	invalidRequest := models.NewUnifyRequest(nil, models.DocumentTypeTaxInvoice, "SA")
	resp, err = client.PushToUnify(ctx, invalidRequest)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))
}

func TestSubmitInvoice(t *testing.T) {
	// Configure the SDK
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	client, err := pkg.NewClient(cfg)
	assert.NoError(t, err)

	// Create a valid source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Create a valid payload
	payload := map[string]interface{}{
		"invoice_number": "INV-001",
		"issue_date":     "2023-01-01",
	}

	// Test valid submission
	ctx := context.Background()
	resp, err := client.SubmitInvoice(ctx, source, "SA", payload)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "success", resp.Status)

	// Test nil source
	resp, err = client.SubmitInvoice(ctx, nil, "SA", payload)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test invalid source
	invalidSource := models.NewSource("", models.SourceTypeFirstParty, "Test Source")
	resp, err = client.SubmitInvoice(ctx, invalidSource, "SA", payload)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test empty country
	resp, err = client.SubmitInvoice(ctx, source, "", payload)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test invalid country format
	resp, err = client.SubmitInvoice(ctx, source, "USA", payload)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test nil payload
	resp, err = client.SubmitInvoice(ctx, source, "SA", nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))
}

func TestCreateMapping(t *testing.T) {
	// Configure the SDK
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	client, err := pkg.NewClient(cfg)
	assert.NoError(t, err)

	// Create a valid source
	source := models.NewSource("test-source", models.SourceTypeFirstParty, "Test Source")

	// Create a valid payload
	payload := map[string]interface{}{
		"invoice_number": "INV-001",
		"issue_date":     "2023-01-01",
	}

	// Test valid mapping
	ctx := context.Background()
	resp, err := client.CreateMapping(ctx, source, "SA", payload)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "success", resp.Status)

	// Test nil source
	resp, err = client.CreateMapping(ctx, nil, "SA", payload)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test invalid source
	invalidSource := models.NewSource("", models.SourceTypeFirstParty, "Test Source")
	resp, err = client.CreateMapping(ctx, invalidSource, "SA", payload)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test empty country
	resp, err = client.CreateMapping(ctx, source, "", payload)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test invalid country format
	resp, err = client.CreateMapping(ctx, source, "USA", payload)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))

	// Test nil payload
	resp, err = client.CreateMapping(ctx, source, "SA", nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))
}

func TestGetStatus(t *testing.T) {
	// Configure the SDK
	cfg := config.New(
		config.WithAPIKey("test-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	client, err := pkg.NewClient(cfg)
	assert.NoError(t, err)

	// Test valid submission ID
	ctx := context.Background()
	resp, err := client.GetStatus(ctx, "sub_123456")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, "sub_123456", resp.Data["submission_id"])
	assert.Equal(t, "PROCESSED", resp.Data["status"])

	// Test empty submission ID
	resp, err = client.GetStatus(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, errors.IsValidationError(err))
}

func TestVersion(t *testing.T) {
	version := pkg.Version()
	assert.NotEmpty(t, version)
}