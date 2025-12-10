/*
Main entry point for the GETS Unify Go SDK - Part 1.

This matches the Python SDK GETSUnifySDK class exactly.
*/
package complyancesdk

import (
	"fmt"
	"log"
	"strings"
)

// GETSUnifySDK Main entry point for the GETS Unify Go SDK
type GETSUnifySDK struct {
	config       *SDKConfig
	apiClient    *APIClient
	queueManager *PersistentQueueManager
}

var globalSDK *GETSUnifySDK

// Configure Configure the SDK with API key, environment, and sources
func Configure(sdkConfig *SDKConfig) error {
	if sdkConfig == nil {
		errorDetail := NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"SDKConfig is required",
		)
		errorDetail.Suggestion = &[]string{"Call GETSUnifySDK.Configure() with a valid SDKConfig."}[0]
		return NewSDKError(errorDetail)
	}

	globalSDK = &GETSUnifySDK{
		config: sdkConfig,
	}

	// Validate country restrictions for production environments
	validateEnvironmentCountryRestrictions(sdkConfig.Environment)

	globalSDK.apiClient = NewAPIClient(
		sdkConfig.APIKey,
		sdkConfig.Environment,
		sdkConfig.RetryConfig,
	)

	// Initialize PersistentQueueManager for handling failed submissions with shared circuit breaker
	globalSDK.queueManager = NewPersistentQueueManager(
		sdkConfig.APIKey,
		sdkConfig.Environment == EnvironmentLocal,
		globalSDK.apiClient.GetCircuitBreaker(),
	)
	
	return nil
}

// validateEnvironmentCountryRestrictions Validate country restrictions based on environment
func validateEnvironmentCountryRestrictions(environment Environment) {
	if environment == EnvironmentSandbox || environment == EnvironmentSimulation || environment == EnvironmentProduction {
		// For production environments, only SA and MY are allowed
		// This validation happens at configuration time, not at request time
		log.Printf("Production environment detected: %s. Only SA and MY countries will be allowed.", environment)
	} else {
		// For development environments, all countries are allowed
		log.Printf("Development environment detected: %s. All countries are allowed.", environment)
	}
}

// SubmitPayload Submit a payload to the GETS Unify API
func SubmitPayload(clientPayloadJSON string, sourceID string, country Country, documentType DocumentType) (*SubmissionResponseOld, error) {
	if globalSDK == nil || globalSDK.config == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"SDK not configured",
		))
	}

	if strings.TrimSpace(clientPayloadJSON) == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Payload is required",
		))
	}

	if strings.TrimSpace(sourceID) == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Source ID is required",
		))
	}

	if country == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Country is required",
		))
	}

	if documentType == "" {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"Document type is required",
		))
	}

	// Find source by ID
	var source *Source
	for _, s := range globalSDK.config.Sources {
		if s.GetID() == sourceID {
			source = s
			break
		}
	}

	if source == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeInvalidSource,
			"Source not found",
		))
	}

	// Validate country restrictions for current environment
	if err := validateCountryForEnvironment(country, globalSDK.config.Environment); err != nil {
		return nil, err
	}

	return globalSDK.apiClient.SendPayload(clientPayloadJSON, source, country, documentType)
}

// GetStatus Get the status of a submission by its ID
func GetStatus(submissionID string) string {
	// Stub: In a real implementation, this would query the API or local cache.
	return "QUEUED"
}

// GetQueueStatus Get queue status and statistics
func GetQueueStatus() string {
	if globalSDK != nil && globalSDK.queueManager != nil {
		status := globalSDK.queueManager.GetQueueStatus()
		return fmt.Sprintf("Persistent Queue Status: %s", status.String())
	}
	return "Queue Manager is not initialized"
}

// GetDetailedQueueStatus Get detailed queue status
func GetDetailedQueueStatus() *QueueStatus {
	if globalSDK != nil && globalSDK.queueManager != nil {
		return globalSDK.queueManager.GetQueueStatus()
	}
	// Return a QueueStatus object with zeros
	return &QueueStatus{
		PendingCount:    0,
		ProcessingCount: 0,
		FailedCount:     0,
		SuccessCount:    0,
		IsRunning:       false,
	}
}

// RetryFailedSubmissions Retry failed submissions
func RetryFailedSubmissions() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		globalSDK.queueManager.RetryFailedSubmissions()
	}
}

// CleanupOldSuccessFiles Clean up old success files
func CleanupOldSuccessFiles(daysToKeep int) {
	if globalSDK != nil && globalSDK.queueManager != nil {
		globalSDK.queueManager.CleanupOldSuccessFiles(daysToKeep)
	}
}

// ClearAllQueues Clear all files from the queue (emergency cleanup)
func ClearAllQueues() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		globalSDK.queueManager.ClearAllQueues()
	} else {
		log.Println("Queue Manager is not initialized")
	}
}

// CleanupDuplicateFiles Clean up duplicate files across queue directories
func CleanupDuplicateFiles() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		globalSDK.queueManager.CleanupDuplicateFiles()
	} else {
		log.Println("Queue Manager is not initialized")
	}
}

// ProcessPendingSubmissions Process pending submissions
func ProcessPendingSubmissions() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		globalSDK.queueManager.ProcessPendingSubmissionsNow()
	}
}

// ProcessQueuedSubmissionsFirst Process queued submissions before handling new requests
func ProcessQueuedSubmissionsFirst() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		// Processing queued submissions
		globalSDK.queueManager.ProcessPendingSubmissionsNow()
	}
}


// validateCountryForEnvironment Validate country restrictions based on current environment
func validateCountryForEnvironment(country Country, environment Environment) error {
	if environment == EnvironmentSandbox || environment == EnvironmentSimulation || environment == EnvironmentProduction {
		// SA is allowed in all production environments
		if country == CountrySA {
			return nil // SA is always allowed
		}

		// MY is only allowed in SANDBOX and PRODUCTION (not SIMULATION)
		if country == CountryMY {
			if environment == EnvironmentSimulation {
				return NewSDKError(NewErrorDetailWithCode(
					ErrorCodeInvalidArgument,
					"Country not allowed for simulation environment",
				))
			}
			return nil // MY is allowed in SANDBOX and PRODUCTION
		}

		// All other countries are blocked in production environments
		return NewSDKError(NewErrorDetailWithCode(
			ErrorCodeInvalidArgument,
			fmt.Sprintf("Country not allowed for production environment. Only SA and MY are allowed for %s. Use DEV/TEST/STAGE for other countries.", environment),
		))
	}

	// For DEV/TEST/STAGE/LOCAL, all countries are allowed
	return nil
}