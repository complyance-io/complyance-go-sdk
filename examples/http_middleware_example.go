package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

// SimpleLogger implements the Logger interface
type SimpleLogger struct{}

func (l *SimpleLogger) Debug(msg string, fields map[string]interface{}) {
	log.Printf("[DEBUG] %s %v\n", msg, fields)
}

func (l *SimpleLogger) Info(msg string, fields map[string]interface{}) {
	log.Printf("[INFO] %s %v\n", msg, fields)
}

func (l *SimpleLogger) Error(msg string, fields map[string]interface{}) {
	log.Printf("[ERROR] %s %v\n", msg, fields)
}

func main() {
	// Get API key from environment
	apiKey := os.Getenv("COMPLYANCE_API_KEY")
	if apiKey == "" {
		log.Fatal("COMPLYANCE_API_KEY environment variable is required")
	}

	// Create SDK configuration
	cfg := config.New().
		WithAPIKey(apiKey).
		WithEnvironment(config.EnvironmentSandbox).
		WithTimeout(30 * time.Second)

	// Add a source
	source := &models.Source{
		ID:   "my-erp-system",
		Type: models.SourceTypeFirstParty,
		Name: "My ERP System",
	}
	cfg.AddSource(source)

	// Configure the SDK
	if err := pkg.Configure(cfg); err != nil {
		log.Fatalf("Failed to configure SDK: %v", err)
	}

	// Create a new server middleware
	middleware, err := pkg.NewServerMiddleware()
	if err != nil {
		log.Fatalf("Failed to create middleware: %v", err)
	}

	// Add logger to middleware
	middleware.WithLogger(&SimpleLogger{})

	// Create HTTP server with routes
	mux := http.NewServeMux()

	// Route for processing invoices
	mux.Handle("/api/invoices", middleware.ProcessInvoice(http.HandlerFunc(handleInvoice)))

	// Route for health check
	mux.HandleFunc("/health", handleHealth)

	// Start HTTP server
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// handleInvoice processes an invoice after it's been processed by the middleware
func handleInvoice(w http.ResponseWriter, r *http.Request) {
	// Get the Complyance response from the context
	response := pkg.GetResponse(r)
	if response == nil {
		http.Error(w, "No response from Complyance middleware", http.StatusInternalServerError)
		return
	}

	// Check if the response indicates an error
	if response.Status == "error" {
		log.Printf("Complyance API error: %s", response.Message)
		http.Error(w, response.Message, http.StatusBadRequest)
		return
	}

	// Process the successful response
	log.Printf("Successfully processed invoice: %s", response.Message)

	// Return the response to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Invoice processed successfully",
		"data":    response.Data,
	})
}

// handleHealth is a simple health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Check SDK status
	status := "healthy"
	version := pkg.Version()

	// Try to make a simple API call to check connectivity
	_, err := pkg.GetStatus(ctx, "test")
	if err != nil {
		// Don't fail the health check for API errors, just log them
		log.Printf("API connectivity check failed: %v", err)
		status = "degraded"
	}

	// Return health status
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  status,
		"version": version,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

// Example of using the middleware with a custom handler chain
func setupCustomHandlerChain() http.Handler {
	// Create middleware
	middleware, _ := pkg.NewServerMiddleware()
	
	// Create a handler chain
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		
		// Process the request through the middleware
		middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This handler will be called after the middleware processes the request
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"success"}`))
		})).ServeHTTP(w, r)
	})
}

// Example of batch processing multiple requests concurrently
func batchProcessExample() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create source
	source := &models.Source{
		ID:   "my-erp-system",
		Type: models.SourceTypeFirstParty,
		Name: "My ERP System",
	}

	// Create multiple requests
	requests := []*models.UnifyRequest{
		models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "US").
			WithPayload(map[string]interface{}{"invoice_number": "INV-001"}),
		models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "GB").
			WithPayload(map[string]interface{}{"invoice_number": "INV-002"}),
		models.NewUnifyRequest(source, models.DocumentTypeTaxInvoice, "FR").
			WithPayload(map[string]interface{}{"invoice_number": "INV-003"}),
	}

	// Process requests concurrently
	responses, errors := pkg.BatchProcess(ctx, requests)

	// Process results
	for i, resp := range responses {
		if errors[i] != nil {
			fmt.Printf("Request %d failed: %v\n", i, errors[i])
			continue
		}
		fmt.Printf("Request %d succeeded: %s\n", i, resp.Status)
	}
}