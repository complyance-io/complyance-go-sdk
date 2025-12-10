package main

import (
	"context"
	"fmt"
	"log"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

func main() {
	// Configure the SDK
	cfg := config.New(
		config.WithAPIKey("your-api-key"),
		config.WithEnvironment(models.EnvironmentSandbox),
	)

	// Initialize the SDK
	if err := pkg.Configure(cfg); err != nil {
		log.Fatalf("Failed to configure SDK: %v", err)
	}

	// Create a source
	source := models.NewSource("erp-system", models.SourceTypeFirstParty, "My ERP System")
	source.WithVersion("1.0.0")

	// Create a context
	ctx := context.Background()

	// Example 1: Submit an invoice using the static method
	invoicePayload := map[string]interface{}{
		"invoice_number": "INV-001",
		"issue_date":     "2023-01-01",
		"customer_name":  "Test Customer",
		"total_amount":   100.50,
		"currency":       "USD",
	}

	invoiceResp, err := pkg.SubmitInvoice(ctx, source, "SA", invoicePayload)
	if err != nil {
		log.Fatalf("Failed to submit invoice: %v", err)
	}

	fmt.Printf("Invoice submission response: %s\n", invoiceResp.Status)
	fmt.Printf("Submission ID: %s\n", invoiceResp.Data["submission_id"])

	// Example 2: Create a mapping using the client instance
	client, err := pkg.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	mappingPayload := map[string]interface{}{
		"invoice_number": "INV-001",
		"issue_date":     "2023-01-01",
		"customer_info": map[string]interface{}{
			"name":    "Test Customer",
			"address": "123 Test St",
		},
	}

	mappingResp, err := client.CreateMapping(ctx, source, "SA", mappingPayload)
	if err != nil {
		log.Fatalf("Failed to create mapping: %v", err)
	}

	fmt.Printf("Mapping response: %s\n", mappingResp.Status)

	// Example 3: Check submission status
	statusResp, err := pkg.GetStatus(ctx, "sub_123456")
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	fmt.Printf("Status: %s\n", statusResp.Data["status"])

	// Example 4: Create and send a custom request
	request := models.NewUnifyRequest(source, models.DocumentTypeCreditNote, "SA")
	request.WithOperation(models.OperationSingle)
	request.WithMode(models.ModeDocuments)
	request.WithPurpose(models.PurposeInvoicing)
	request.WithPayload(map[string]interface{}{
		"credit_note_number": "CN-001",
		"reference_invoice":  "INV-001",
		"issue_date":         "2023-01-15",
		"total_amount":       50.25,
	})

	customResp, err := client.PushToUnify(ctx, request)
	if err != nil {
		log.Fatalf("Failed to push to unify: %v", err)
	}

	fmt.Printf("Custom request response: %s\n", customResp.Status)
}