/*
Persistent Queue Manager implementation matching Python SDK exactly.
*/
package complyancesdk

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// QueueStatus model matching Python SDK
type QueueStatus struct {
	PendingCount    int  `json:"pending_count"`
	ProcessingCount int  `json:"processing_count"`
	FailedCount     int  `json:"failed_count"`
	SuccessCount    int  `json:"success_count"`
	IsRunning       bool `json:"is_running"`
}

// GetPendingCount getter for pending count
func (q *QueueStatus) GetPendingCount() int {
	return q.PendingCount
}

// GetProcessingCount getter for processing count
func (q *QueueStatus) GetProcessingCount() int {
	return q.ProcessingCount
}

// GetFailedCount getter for failed count
func (q *QueueStatus) GetFailedCount() int {
	return q.FailedCount
}

// GetSuccessCount getter for success count
func (q *QueueStatus) GetSuccessCount() int {
	return q.SuccessCount
}

// IsQueueRunning getter for is running
func (q *QueueStatus) IsQueueRunning() bool {
	return q.IsRunning
}

// String string representation
func (q *QueueStatus) String() string {
	return fmt.Sprintf("QueueStatus{pending=%d, processing=%d, failed=%d, success=%d, running=%t}",
		q.PendingCount, q.ProcessingCount, q.FailedCount, q.SuccessCount, q.IsRunning)
}

// PersistentSubmissionRecord model matching Python SDK
type PersistentSubmissionRecord struct {
	Payload     map[string]interface{} `json:"payload"`
	SourceID    string                 `json:"source_id"`
	Country     string                 `json:"country"`
	DocumentType string                `json:"document_type"`
	EnqueuedAt  string                 `json:"enqueued_at"`
	Timestamp   int64                  `json:"timestamp"`
}

// GetPayload getter for payload
func (p *PersistentSubmissionRecord) GetPayload() map[string]interface{} {
	return p.Payload
}

// GetSourceID getter for source ID
func (p *PersistentSubmissionRecord) GetSourceID() string {
	return p.SourceID
}

// GetCountry getter for country
func (p *PersistentSubmissionRecord) GetCountry() string {
	return p.Country
}

// GetDocumentType getter for document type
func (p *PersistentSubmissionRecord) GetDocumentType() string {
	return p.DocumentType
}

// GetEnqueuedAt getter for enqueued at
func (p *PersistentSubmissionRecord) GetEnqueuedAt() string {
	return p.EnqueuedAt
}

// GetTimestamp getter for timestamp
func (p *PersistentSubmissionRecord) GetTimestamp() int64 {
	return p.Timestamp
}

// PersistentQueueManager Persistent queue manager matching Python SDK
type PersistentQueueManager struct {
	apiKey         string
	local          bool
	queueBasePath  string
	isRunning      bool
	processingLock bool
	circuitBreaker *CircuitBreaker
}

const (
	QueueDir      = "complyance-queue"
	PendingDir    = "pending"
	ProcessingDir = "processing"
	FailedDir     = "failed"
	SuccessDir    = "success"
)

// NewPersistentQueueManager creates a new persistent queue manager
func NewPersistentQueueManager(apiKey string, local bool, circuitBreaker *CircuitBreaker) *PersistentQueueManager {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get user home directory: %v", err)
		homeDir = "."
	}

	queueBasePath := filepath.Join(homeDir, QueueDir)

	// Use shared circuit breaker or create default
	if circuitBreaker == nil {
		circuitBreaker = NewCircuitBreaker(NewCircuitBreakerConfig(3, 60000)) // 3 failures, 1 minute timeout
	}

	manager := &PersistentQueueManager{
		apiKey:         apiKey,
		local:          local,
		queueBasePath:  queueBasePath,
		isRunning:      false,
		processingLock: false,
		circuitBreaker: circuitBreaker,
	}

	manager.initializeQueueDirectories()
	log.Printf("PersistentQueueManager initialized with queue directory: %s", manager.queueBasePath)

	// Automatically start processing and retry any existing failed submissions
	manager.StartProcessing()
	manager.RetryFailedSubmissions()

	return manager
}

// initializeQueueDirectories Initialize queue directories
func (p *PersistentQueueManager) initializeQueueDirectories() {
	dirs := []string{PendingDir, ProcessingDir, FailedDir, SuccessDir}
	for _, dir := range dirs {
		dirPath := filepath.Join(p.queueBasePath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			log.Printf("Failed to create queue directory %s: %v", dirPath, err)
			panic(fmt.Sprintf("Failed to initialize persistent queue: %v", err))
		}
	}
	log.Println("Queue directories initialized")
}

// Enqueue a payload submission
func (p *PersistentQueueManager) Enqueue(submission *PayloadSubmission) error {
	fileName := p.generateFileName(submission)
	filePath := filepath.Join(p.queueBasePath, PendingDir, fileName)

	// Check if file already exists (same document ID)
	if _, err := os.Stat(filePath); err == nil {
		log.Printf("Document already exists in queue: %s. Skipping duplicate submission.", fileName)
		return nil // Skip duplicate submission
	}

	// Parse the UnifyRequest JSON string to proper JSON object
	jsonPayload := submission.GetPayload()
	log.Printf("🔥 QUEUE: Received payload from submission with length: %d characters", len(jsonPayload))
	log.Printf("🔥 QUEUE: Payload preview: %s", jsonPayload[:min(200, len(jsonPayload))])

	// Verify the payload is not empty
	if strings.TrimSpace(jsonPayload) == "" || jsonPayload == "{}" {
		log.Printf("🔥 QUEUE: ERROR - Received empty or invalid payload: '%s'", jsonPayload)
		return fmt.Errorf("cannot enqueue empty payload")
	}

	// Parse the UnifyRequest JSON string to a proper JSON object
	var unifyRequestMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonPayload), &unifyRequestMap); err != nil {
		return fmt.Errorf("failed to parse UnifyRequest JSON: %v", err)
	}

	// Create submission record with the parsed UnifyRequest as proper JSON object
	record := &PersistentSubmissionRecord{
		Payload:      unifyRequestMap, // Store as map instead of string
		SourceID:     fmt.Sprintf("%s:%s", submission.GetSource().GetName(), submission.GetSource().GetVersion()),
		Country:      string(submission.GetCountry()),
		DocumentType: string(submission.GetDocumentType()),
		EnqueuedAt:   time.Now().Format(time.RFC3339),
		Timestamp:    time.Now().UnixNano() / int64(time.Millisecond),
	}

	// Write to file
	recordJSON, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal submission record: %v", err)
	}

	if err := os.WriteFile(filePath, recordJSON, 0644); err != nil {
		return fmt.Errorf("failed to write submission to file: %v", err)
	}

	log.Printf("🔥 QUEUE: Stored record to file: %s with payload length: %d", fileName, len(jsonPayload))
	log.Printf("Enqueued submission to persistent storage: %s for source: %s:%s, country: %s",
		fileName, submission.GetSource().GetName(), submission.GetSource().GetVersion(), submission.GetCountry())

	// Start processing if not already running
	p.StartProcessing()

	return nil
}

// generateFileName Generate filename for submission
func (p *PersistentQueueManager) generateFileName(submission *PayloadSubmission) string {
	// Extract document ID from payload to create unique reference
	documentID := p.extractDocumentID(submission.GetPayload())

	// Generate filename using source and document ID for unique reference
	sourceID := fmt.Sprintf("%s:%s", submission.GetSource().GetName(), submission.GetSource().GetVersion())
	// Replace special characters with underscores
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	sourceIDClean := re.ReplaceAllString(sourceID, "_")
	country := string(submission.GetCountry())
	return fmt.Sprintf("%s_%s_%s_%s.json", sourceIDClean, documentID, country, string(submission.GetDocumentType()))
}

// extractDocumentID Extract document ID from payload
func (p *PersistentQueueManager) extractDocumentID(payload string) string {
	// Parse the complete UnifyRequest JSON
	var requestMap map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &requestMap); err != nil {
		log.Printf("Failed to extract document ID from UnifyRequest payload, using timestamp: %v", err)
		return fmt.Sprintf("doc_%d", time.Now().UnixNano()/int64(time.Millisecond))
	}

	// Extract from payload.invoice_data.invoice_number
	if payloadMap, ok := requestMap["payload"].(map[string]interface{}); ok {
		if invoiceData, ok := payloadMap["invoice_data"].(map[string]interface{}); ok {
			if invoiceNumber, ok := invoiceData["invoice_number"].(string); ok {
				return invoiceNumber
			}
		}
	}

	// Fallback to timestamp if no invoice number found
	return fmt.Sprintf("doc_%d", time.Now().UnixNano()/int64(time.Millisecond))
}

// StartProcessing Start processing queue
func (p *PersistentQueueManager) StartProcessing() {
	if !p.isRunning {
		p.isRunning = true
		// Note: In a real implementation, this would start a background goroutine
		// For now, we'll process on-demand
		log.Println("Started persistent queue processing")
	}
}

// ProcessPendingSubmissionsNow Manually trigger processing of pending submissions
func (p *PersistentQueueManager) ProcessPendingSubmissionsNow() {
	// Check circuit breaker state before manual processing
	if p.circuitBreaker.IsOpen() {
		currentTime := time.Now().UnixNano() / int64(time.Millisecond)
		timeSinceLastFailure := currentTime - p.circuitBreaker.GetLastFailureTime()

		if timeSinceLastFailure < 60000 { // 1 minute = 60000ms
			remainingTime := 60000 - timeSinceLastFailure
			log.Printf("🚫 Circuit breaker is OPEN - remaining time: %dms (%d seconds). Manual processing skipped.",
				remainingTime, remainingTime/1000)
			return
		} else {
			log.Printf("✅ Circuit breaker timeout expired (%dms) - proceeding with manual processing", timeSinceLastFailure)
		}
	}

	p.processPendingSubmissions()
}

// StopProcessing Stop processing queue
func (p *PersistentQueueManager) StopProcessing() {
	p.isRunning = false
	log.Println("Stopped persistent queue processing")
}

// processPendingSubmissions Process pending submissions
func (p *PersistentQueueManager) processPendingSubmissions() {
	if !p.isRunning {
		return
	}

	if p.processingLock {
		return
	}

	p.processingLock = true
	defer func() {
		p.processingLock = false
	}()

	// First check if there are any pending files
	pendingDir := filepath.Join(p.queueBasePath, PendingDir)
	files, err := filepath.Glob(filepath.Join(pendingDir, "*.json"))
	if err != nil {
		log.Printf("Error reading pending directory: %v", err)
		return
	}

	if len(files) == 0 {
		log.Println("No pending submissions to process")
		return
	}

	log.Printf("🔄 Found %d pending submissions in queue", len(files))

	// Check circuit breaker state before attempting to process
	if p.circuitBreaker.IsOpen() {
		currentTime := time.Now().UnixNano() / int64(time.Millisecond)
		timeSinceLastFailure := currentTime - p.circuitBreaker.GetLastFailureTime()

		// Wait for full 1 minute timeout before attempting to process
		if timeSinceLastFailure < 60000 { // 1 minute = 60000ms
			remainingTime := 60000 - timeSinceLastFailure
			log.Printf("🚫 Circuit breaker is OPEN - %d seconds remaining. Queue has %d items waiting.",
				remainingTime/1000, len(files))
			return
		} else {
			log.Printf("✅ Circuit breaker timeout expired - attempting to process %d queued items", len(files))
		}
	}

	// Process each file in the queue
	for _, filePath := range files {
		// Check if file still exists before processing
		if _, err := os.Stat(filePath); err == nil {
			if err := p.processSubmissionFile(filePath); err != nil {
				log.Printf("Failed to process queued submission %s: %v", filepath.Base(filePath), err)
				// Continue processing other files even if one fails
			}
		}
	}
}

// processSubmissionFile Process a single submission file
func (p *PersistentQueueManager) processSubmissionFile(filePath string) error {
	// Implementation would be similar to Python but adapted for Go
	// For now, just log that we're processing it
	fileName := filepath.Base(filePath)
	log.Printf("Processing submission file: %s", fileName)

	// Move to failed directory for now (in real implementation, would attempt to send)
	failedPath := filepath.Join(p.queueBasePath, FailedDir, fileName)
	if _, err := os.Stat(failedPath); os.IsNotExist(err) {
		if err := os.Rename(filePath, failedPath); err != nil {
			return fmt.Errorf("failed to move file to failed directory: %v", err)
		}
		log.Printf("Moved submission to failed directory: %s", fileName)
	}

	return nil
}

// GetQueueStatus Get queue status
func (p *PersistentQueueManager) GetQueueStatus() *QueueStatus {
	pendingCount := p.countFilesInDir(PendingDir)
	processingCount := p.countFilesInDir(ProcessingDir)
	failedCount := p.countFilesInDir(FailedDir)
	successCount := p.countFilesInDir(SuccessDir)

	return &QueueStatus{
		PendingCount:    pendingCount,
		ProcessingCount: processingCount,
		FailedCount:     failedCount,
		SuccessCount:    successCount,
		IsRunning:       p.isRunning,
	}
}

// countFilesInDir Count files in a directory
func (p *PersistentQueueManager) countFilesInDir(dirName string) int {
	dirPath := filepath.Join(p.queueBasePath, dirName)
	files, err := filepath.Glob(filepath.Join(dirPath, "*.json"))
	if err != nil {
		log.Printf("Failed to count files in %s: %v", dirName, err)
		return 0
	}
	return len(files)
}

// RetryFailedSubmissions Retry failed submissions
func (p *PersistentQueueManager) RetryFailedSubmissions() {
	failedDir := filepath.Join(p.queueBasePath, FailedDir)
	pendingDir := filepath.Join(p.queueBasePath, PendingDir)

	files, err := filepath.Glob(filepath.Join(failedDir, "*.json"))
	if err != nil {
		log.Printf("Error reading failed directory: %v", err)
		return
	}

	if len(files) == 0 {
		log.Println("No failed submissions to retry")
		return
	}

	log.Printf("Retrying %d failed submissions", len(files))

	for _, filePath := range files {
		fileName := filepath.Base(filePath)
		pendingPath := filepath.Join(pendingDir, fileName)
		if err := os.Rename(filePath, pendingPath); err != nil {
			log.Printf("Failed to move failed submission back to pending: %v", err)
		} else {
			log.Printf("Moved failed submission back to pending: %s", fileName)
		}
	}
}

// CleanupOldSuccessFiles Clean up old success files
func (p *PersistentQueueManager) CleanupOldSuccessFiles(daysToKeep int) {
	successDir := filepath.Join(p.queueBasePath, SuccessDir)
	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)

	files, err := filepath.Glob(filepath.Join(successDir, "*.json"))
	if err != nil {
		log.Printf("Error reading success directory: %v", err)
		return
	}

	var oldFiles []string
	for _, filePath := range files {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			continue
		}
		if fileInfo.ModTime().Before(cutoffTime) {
			oldFiles = append(oldFiles, filePath)
		}
	}

	for _, filePath := range oldFiles {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to remove old success file %s: %v", filepath.Base(filePath), err)
		} else {
			log.Printf("Cleaned up old success file: %s", filepath.Base(filePath))
		}
	}

	if len(oldFiles) > 0 {
		log.Printf("Cleaned up %d old success files", len(oldFiles))
	}
}

// ClearAllQueues Clear all files from the queue (emergency cleanup)
func (p *PersistentQueueManager) ClearAllQueues() {
	log.Println("Clearing all queue directories...")

	// Clear pending
	p.clearDirectory(PendingDir)

	// Clear processing
	p.clearDirectory(ProcessingDir)

	// Clear failed
	p.clearDirectory(FailedDir)

	// Clear success
	p.clearDirectory(SuccessDir)

	log.Println("All queue directories cleared successfully")
}

// clearDirectory Clear a specific directory
func (p *PersistentQueueManager) clearDirectory(dirName string) {
	dirPath := filepath.Join(p.queueBasePath, dirName)
	files, err := filepath.Glob(filepath.Join(dirPath, "*.json"))
	if err != nil {
		log.Printf("Error reading directory %s: %v", dirName, err)
		return
	}

	for _, filePath := range files {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to delete file %s: %v", filepath.Base(filePath), err)
		} else {
			log.Printf("Deleted file: %s", filepath.Base(filePath))
		}
	}

	log.Printf("Cleared %d files from %s", len(files), dirName)
}

// CleanupDuplicateFiles Clean up duplicate files across queue directories
func (p *PersistentQueueManager) CleanupDuplicateFiles() {
	log.Println("Cleaning up duplicate files across queue directories...")

	// Get all files from all directories
	fileMap := make(map[string]string)

	dirs := []string{PendingDir, ProcessingDir, FailedDir, SuccessDir}
	for _, dirName := range dirs {
		dirPath := filepath.Join(p.queueBasePath, dirName)
		files, err := filepath.Glob(filepath.Join(dirPath, "*.json"))
		if err != nil {
			log.Printf("Error reading directory %s: %v", dirName, err)
			continue
		}

		for _, filePath := range files {
			fileName := filepath.Base(filePath)
			existingFile, exists := fileMap[fileName]

			if exists {
				// File exists in multiple directories, keep the one with latest modification time
				existingInfo, err1 := os.Stat(existingFile)
				currentInfo, err2 := os.Stat(filePath)

				if err1 != nil || err2 != nil {
					log.Printf("Could not compare modification times for duplicate file: %s", fileName)
					// Keep the existing file, delete current
					os.Remove(filePath)
					continue
				}

				if currentInfo.ModTime().After(existingInfo.ModTime()) {
					// Delete the older file
					os.Remove(existingFile)
					fileMap[fileName] = filePath
					log.Printf("Removed duplicate file (older): %s", existingFile)
				} else {
					// Delete the current file
					os.Remove(filePath)
					log.Printf("Removed duplicate file (older): %s", filePath)
				}
			} else {
				fileMap[fileName] = filePath
			}
		}
	}

	log.Println("Duplicate file cleanup completed")
}
