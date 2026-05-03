/*
Main entry point for the GETS Unify Go SDK - Part 1.

This matches the Python SDK GETSUnifySDK class exactly.
*/
package complyancesdk

import (
	"fmt"
	"log"
	"strings"
	"time"
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
		// For production environments, only SA, MY, and AE (UAE) are allowed
		// This validation happens at configuration time, not at request time
		log.Printf("Production environment detected: %s. Only SA (Saudi Arabia), MY (Malaysia), and AE (UAE) countries will be allowed.", environment)
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

// GetDocumentStatus gets retrieval status by documentId.
func GetDocumentStatus(documentID string) (map[string]interface{}, error) {
	if globalSDK == nil || globalSDK.apiClient == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"SDK not configured",
		).WithSuggestion("Call Configure() first."))
	}

	return globalSDK.apiClient.GetDocumentStatus(documentID)
}

// GetSubmissionStatus is deprecated and intentionally blocked.
func GetSubmissionStatus(submissionID string) (map[string]interface{}, error) {
	if globalSDK == nil || globalSDK.apiClient == nil {
		return nil, NewSDKError(NewErrorDetailWithCode(
			ErrorCodeMissingField,
			"SDK not configured",
		).WithSuggestion("Call Configure() first."))
	}

	return globalSDK.apiClient.GetSubmissionStatus(submissionID)
}

// GetStatus is deprecated and forwards to the deprecated submissionId endpoint behavior.
func GetStatus(submissionID string) (map[string]interface{}, error) {
	return GetSubmissionStatus(submissionID)
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

func GetQueueStatusDetailed() *QueueStatusDetailed {
	if globalSDK != nil && globalSDK.queueManager != nil {
		return globalSDK.queueManager.GetQueueStatusDetailed()
	}
	return &QueueStatusDetailed{
		PendingCount:    0,
		ProcessingCount: 0,
		FailedCount:     0,
		SuccessCount:    0,
		TotalCount:      0,
		IsRunning:       false,
		IsPaused:        false,
		QueueDir:        "",
	}
}

// RetryFailedSubmissions Retry failed submissions
func RetryFailedSubmissions() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		globalSDK.queueManager.RetryFailedSubmissions()
	}
}

func RetryFailed(queueItemID string) bool {
	if globalSDK != nil && globalSDK.queueManager != nil {
		return globalSDK.queueManager.RetryFailed(queueItemID)
	}
	return false
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

func PauseQueueProcessing() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		globalSDK.queueManager.PauseProcessing()
	}
}

func ResumeQueueProcessing() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		globalSDK.queueManager.ResumeProcessing()
	}
}

func DrainQueue(timeout time.Duration) bool {
	if globalSDK != nil && globalSDK.queueManager != nil {
		return globalSDK.queueManager.DrainQueue(timeout)
	}
	return true
}

// ProcessQueuedSubmissionsFirst Process queued submissions before handling new requests
func ProcessQueuedSubmissionsFirst() {
	if globalSDK != nil && globalSDK.queueManager != nil {
		// Processing queued submissions
		globalSDK.queueManager.ProcessPendingSubmissionsNow()
	}
}

// validateCountryForEnvironment Validate country restrictions based on current environment
// Implements the three-tier country access control:
// - SA: Allowed in all production environments (SANDBOX, SIMULATION, PRODUCTION)
// - MY: Allowed in SANDBOX and PRODUCTION only (blocked in SIMULATION)
// - AE: Allowed in SANDBOX and PRODUCTION only (blocked in SIMULATION)
// - Others: Blocked in all production environments
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
					"Country not allowed for simulation environment. MY (Malaysia) is not allowed in SIMULATION environment. Use SANDBOX or PRODUCTION.",
				))
			}
			return nil // MY is allowed in SANDBOX and PRODUCTION
		}

		// AE (UAE) is only allowed in SANDBOX and PRODUCTION (not SIMULATION)
		if country == CountryAE {
			if environment == EnvironmentSimulation {
				return NewSDKError(NewErrorDetailWithCode(
					ErrorCodeInvalidArgument,
					"Country not allowed for simulation environment. AE (UAE) is not allowed in SIMULATION environment. Use SANDBOX or PRODUCTION.",
				))
			}
			return nil // AE is allowed in SANDBOX and PRODUCTION
		}

		// All other countries are blocked in production environments
		return NewSDKError(NewErrorDetailWithCode(
			ErrorCodeInvalidArgument,
			fmt.Sprintf("Country not allowed for production environment. Only SA, MY, and AE are allowed for %s. Use DEV/TEST/STAGE for other countries.", environment),
		))
	}

	// For DEV/TEST/STAGE/LOCAL, all countries are allowed
	return nil
}
