// Package complyance provides a Go SDK for the Complyance Unified e-invoicing platform.
//
// The Complyance SDK enables developers to integrate e-invoicing compliance
// across multiple countries and document types with a simple, idiomatic Go API.
//
// Basic usage:
//
//	import "github.com/complyance-io/complyance-go-sdk/v3"
//
//	// Configure the SDK
//	config := &complyance.SDKConfig{
//		APIKey:      "ak_your_api_key",
//		Environment: complyance.EnvironmentSandbox,
//	}
//	
//	sdk, err := complyance.NewSDK(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Create and submit a request
//	request := &complyance.UnifyRequest{
//		Source:       source,
//		DocumentType: complyance.DocumentTypeTaxInvoice,
//		Country:      "US",
//		Operation:    complyance.OperationSingle,
//		Mode:         complyance.ModeDocuments,
//		Purpose:      complyance.PurposeInvoicing,
//		Payload:      invoiceData,
//	}
//
//	response, err := sdk.PushToUnify(context.Background(), request)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// For more information, visit https://docs.complyance.io/sdk/go
package complyance