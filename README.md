# Complyance Go SDK

The official Go SDK for the Complyance Unified e-invoicing platform.

## Installation

```bash
go get github.com/complyance-io/sdk-go/v1
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/complyance-io/sdk-go/v1/pkg/config"
	"github.com/complyance-io/sdk-go/v1/pkg/models"
	sdk "github.com/complyance-io/sdk-go/v1/pkg"
)

func main() {
	// Configure the SDK
	cfg := config.New(
		config.WithAPIKey("your-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	// Initialize the SDK client
	client, err := sdk.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize SDK client: %v", err)
	}

	// Create a request
	req := &models.UnifyRequest{
		Source: &models.Source{
			ID:   "source-id",
			Type: models.SourceTypeFirstParty,
			Name: "My ERP System",
		},
		DocumentType: models.DocumentTypeTaxInvoice,
		Country:      "SA",
		Operation:    models.OperationSingle,
		Mode:         models.ModeDocuments,
		Purpose:      models.PurposeInvoicing,
		Payload:      map[string]interface{}{
			"invoice_number": "INV-001",
			// Add your invoice data here
		},
	}

	// Send the request
	ctx := context.Background()
	resp, err := client.PushToUnify(ctx, req)
	if err != nil {
		log.Fatalf("Failed to push to unify: %v", err)
	}

	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response message: %s\n", resp.Message)
}
```

## Features

- Idiomatic Go API with context support
- Goroutine-safe operations
- Automatic retries with exponential backoff
- Circuit breaker for resilience
- Comprehensive error handling
- HTTP middleware for easy integration

## Documentation

For detailed documentation, visit [https://docs.complyance.io/sdk/go](https://docs.complyance.io/sdk/go)

## License

MIT
