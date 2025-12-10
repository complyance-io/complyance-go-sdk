/*
Data models for the Complyance SDK.

This module contains all the data models that match the Python SDK exactly.
*/
package complyancesdk

import (
	"fmt"
)

// Environment enumeration matching Python SDK exactly
type Environment string

const (
	EnvironmentDev        Environment = "DEV"
	EnvironmentTest       Environment = "TEST"
	EnvironmentStage      Environment = "STAGE"
	EnvironmentLocal      Environment = "LOCAL"
	EnvironmentSandbox    Environment = "SANDBOX"
	EnvironmentSimulation Environment = "SIMULATION"
	EnvironmentProduction Environment = "PRODUCTION"
)

// GetBaseURL Get the base URL for this environment (matching Python SDK)
func (e Environment) GetBaseURL() string {
	switch e {
	case EnvironmentDev:
		return "https://prod.gets.complyance.io/unify"
	case EnvironmentTest:
		return "https://prod.gets.complyance.io/unify"
	case EnvironmentStage:
		return "https://prod.gets.complyance.io/unify"
	case EnvironmentLocal:
		return "http://127.0.0.1:4000/unify"
	case EnvironmentSandbox:
		return "https://prod.gets.complyance.io/unify"
	case EnvironmentSimulation:
		return "https://prod.gets.complyance.io/unify"
	case EnvironmentProduction:
		return "https://prod.gets.complyance.io/unify"
	default:
		return "https://prod.gets.complyance.io/unify" // Default to dev	
}

// Country enumeration matching Python SDK
type Country string

const (
	CountrySA Country = "SA" // Saudi Arabia
	CountryMY Country = "MY" // Malaysia
	CountryAE Country = "AE" // UAE
	CountrySG Country = "SG" // Singapore
)

// DocumentType enumeration matching Python SDK
type DocumentType string

const (
	DocumentTypeTaxInvoice                         DocumentType = "tax_invoice"
	DocumentTypeSimplifiedInvoice                  DocumentType = "simplified_invoice"
	DocumentTypeCreditNote                         DocumentType = "credit_note"
	DocumentTypeSimplifiedCreditNote               DocumentType = "simplified_credit_note"
	DocumentTypeDebitNote                          DocumentType = "debit_note"
	DocumentTypeSimplifiedDebitNote                DocumentType = "simplified_debit_note"
	DocumentTypePrepaymentInvoice                  DocumentType = "prepayment_invoice"
	DocumentTypeSimplifiedPrepaymentInvoice        DocumentType = "simplified_prepayment_invoice"
	DocumentTypePrepaymentAdjustedInvoice          DocumentType = "prepayment_adjusted_invoice"
	DocumentTypeSimplifiedPrepaymentAdjustedInvoice DocumentType = "simplified_prepayment_adjusted_invoice"
)

// FromString Convert string to DocumentType enum
func (d DocumentType) FromString(value string) DocumentType {
	switch value {
	case "tax_invoice":
		return DocumentTypeTaxInvoice
	case "simplified_invoice":
		return DocumentTypeSimplifiedInvoice
	case "credit_note":
		return DocumentTypeCreditNote
	case "simplified_credit_note":
		return DocumentTypeSimplifiedCreditNote
	case "debit_note":
		return DocumentTypeDebitNote
	case "simplified_debit_note":
		return DocumentTypeSimplifiedDebitNote
	case "prepayment_invoice":
		return DocumentTypePrepaymentInvoice
	case "simplified_prepayment_invoice":
		return DocumentTypeSimplifiedPrepaymentInvoice
	case "prepayment_adjusted_invoice":
		return DocumentTypePrepaymentAdjustedInvoice
	case "simplified_prepayment_adjusted_invoice":
		return DocumentTypeSimplifiedPrepaymentAdjustedInvoice
	default:
		// Return empty string for unknown types
		return ""
	}
}

func (d DocumentType) String() string {
	return string(d)
}

// LogicalDocType enumeration matching Python SDK
type LogicalDocType string

const (
	// Base types
	LogicalDocTypeInvoice    LogicalDocType = "INVOICE"
	LogicalDocTypeCreditNote LogicalDocType = "CREDIT_NOTE"
	LogicalDocTypeDebitNote  LogicalDocType = "DEBIT_NOTE"
	LogicalDocTypeReceipt    LogicalDocType = "RECEIPT"

	// B2B Tax Invoice types
	LogicalDocTypeTaxInvoice                      LogicalDocType = "TAX_INVOICE"
	LogicalDocTypeTaxInvoiceCreditNote            LogicalDocType = "TAX_INVOICE_CREDIT_NOTE"
	LogicalDocTypeTaxInvoiceDebitNote             LogicalDocType = "TAX_INVOICE_DEBIT_NOTE"
	LogicalDocTypeTaxInvoicePrepayment            LogicalDocType = "TAX_INVOICE_PREPAYMENT"
	LogicalDocTypeTaxInvoicePrepaymentAdjusted    LogicalDocType = "TAX_INVOICE_PREPAYMENT_ADJUSTED"
	LogicalDocTypeTaxInvoiceExportInvoice         LogicalDocType = "TAX_INVOICE_EXPORT_INVOICE"
	LogicalDocTypeTaxInvoiceExportCreditNote      LogicalDocType = "TAX_INVOICE_EXPORT_CREDIT_NOTE"
	LogicalDocTypeTaxInvoiceExportDebitNote       LogicalDocType = "TAX_INVOICE_EXPORT_DEBIT_NOTE"
	LogicalDocTypeTaxInvoiceThirdPartyInvoice     LogicalDocType = "TAX_INVOICE_THIRD_PARTY_INVOICE"
	LogicalDocTypeTaxInvoiceSelfBilledInvoice     LogicalDocType = "TAX_INVOICE_SELF_BILLED_INVOICE"
	LogicalDocTypeTaxInvoiceNominalSupplyInvoice  LogicalDocType = "TAX_INVOICE_NOMINAL_SUPPLY_INVOICE"
	LogicalDocTypeTaxInvoiceSummaryInvoice        LogicalDocType = "TAX_INVOICE_SUMMARY_INVOICE"

	// B2C Simplified Tax Invoice types
	LogicalDocTypeSimplifiedTaxInvoice                      LogicalDocType = "SIMPLIFIED_TAX_INVOICE"
	LogicalDocTypeSimplifiedTaxInvoiceCreditNote            LogicalDocType = "SIMPLIFIED_TAX_INVOICE_CREDIT_NOTE"
	LogicalDocTypeSimplifiedTaxInvoiceDebitNote             LogicalDocType = "SIMPLIFIED_TAX_INVOICE_DEBIT_NOTE"
	LogicalDocTypeSimplifiedTaxInvoicePrepayment            LogicalDocType = "SIMPLIFIED_TAX_INVOICE_PREPAYMENT"
	LogicalDocTypeSimplifiedTaxInvoicePrepaymentAdjusted    LogicalDocType = "SIMPLIFIED_TAX_INVOICE_PREPAYMENT_ADJUSTED"
	LogicalDocTypeSimplifiedTaxInvoiceExportInvoice         LogicalDocType = "SIMPLIFIED_TAX_INVOICE_EXPORT_INVOICE"
	LogicalDocTypeSimplifiedTaxInvoiceExportCreditNote      LogicalDocType = "SIMPLIFIED_TAX_INVOICE_EXPORT_CREDIT_NOTE"
	LogicalDocTypeSimplifiedTaxInvoiceExportDebitNote       LogicalDocType = "SIMPLIFIED_TAX_INVOICE_EXPORT_DEBIT_NOTE"
	LogicalDocTypeSimplifiedTaxInvoiceThirdPartyInvoice     LogicalDocType = "SIMPLIFIED_TAX_INVOICE_THIRD_PARTY_INVOICE"
	LogicalDocTypeSimplifiedTaxInvoiceSelfBilledInvoice     LogicalDocType = "SIMPLIFIED_TAX_INVOICE_SELF_BILLED_INVOICE"
	LogicalDocTypeSimplifiedTaxInvoiceNominalSupplyInvoice  LogicalDocType = "SIMPLIFIED_TAX_INVOICE_NOMINAL_SUPPLY_INVOICE"
	LogicalDocTypeSimplifiedTaxInvoiceSummaryInvoice        LogicalDocType = "SIMPLIFIED_TAX_INVOICE_SUMMARY_INVOICE"

	// Country-specific logical types
	LogicalDocTypeExportInvoice            LogicalDocType = "EXPORT_INVOICE"
	LogicalDocTypeExportCreditNote         LogicalDocType = "EXPORT_CREDIT_NOTE"
	LogicalDocTypeExportThirdPartyInvoice  LogicalDocType = "EXPORT_THIRD_PARTY_INVOICE"
	LogicalDocTypeThirdPartyInvoice        LogicalDocType = "THIRD_PARTY_INVOICE"
	LogicalDocTypeSelfBilledInvoice        LogicalDocType = "SELF_BILLED_INVOICE"
	LogicalDocTypeNominalSupplyInvoice     LogicalDocType = "NOMINAL_SUPPLY_INVOICE"
	LogicalDocTypeSummaryInvoice           LogicalDocType = "SUMMARY_INVOICE"
)

// Operation types matching Python SDK
type Operation string

const (
	OperationSingle Operation = "single"
	OperationBulk   Operation = "bulk"
)

// FromString Convert string to Operation enum
func (o Operation) FromString(value string) Operation {
	switch value {
	case "single":
		return OperationSingle
	case "bulk":
		return OperationBulk
	default:
		return ""
	}
}

func (o Operation) String() string {
	return string(o)
}

// Mode types matching Python SDK
type Mode string

const (
	ModeDocuments  Mode = "documents"
	ModeOnboarding Mode = "onboarding"
)

// FromString Convert string to Mode enum
func (m Mode) FromString(value string) Mode {
	switch value {
	case "documents":
		return ModeDocuments
	case "onboarding":
		return ModeOnboarding
	default:
		return ""
	}
}

func (m Mode) String() string {
	return string(m)
}

// Purpose types matching Python SDK
type Purpose string

const (
	PurposeMapping   Purpose = "mapping"
	PurposeInvoicing Purpose = "invoicing"
)

// FromString Convert string to Purpose enum
func (p Purpose) FromString(value string) Purpose {
	switch value {
	case "mapping":
		return PurposeMapping
	case "invoicing":
		return PurposeInvoicing
	default:
		return ""
	}
}

func (p Purpose) String() string {
	return string(p)
}

// SourceType enumeration matching Python SDK
type SourceType string

const (
	SourceTypeFirstParty  SourceType = "FIRST_PARTY"
	SourceTypeThirdParty  SourceType = "THIRD_PARTY"
	SourceTypeMarketplace SourceType = "MARKETPLACE"
)

// DestinationType enumeration matching Python SDK
type DestinationType string

const (
	DestinationTypeTaxAuthority DestinationType = "TAX_AUTHORITY"
	DestinationTypeEmail        DestinationType = "EMAIL"
	DestinationTypeArchive      DestinationType = "ARCHIVE"
	DestinationTypePeppol       DestinationType = "PEPPOL"
)

// ErrorCode enumeration matching Python SDK
type ErrorCode string

const (
	ErrorCodeMissingField                  ErrorCode = "MISSING_FIELD"
	ErrorCodeInvalidSource                 ErrorCode = "INVALID_SOURCE"
	ErrorCodeInvalidArgument               ErrorCode = "INVALID_ARGUMENT"
	ErrorCodeAuthenticationFailed          ErrorCode = "AUTHENTICATION_FAILED"
	ErrorCodeAuthorizationDenied           ErrorCode = "AUTHORIZATION_DENIED"
	ErrorCodeValidationFailed              ErrorCode = "VALIDATION_FAILED"
	ErrorCodeTemplateNotFound              ErrorCode = "TEMPLATE_NOT_FOUND"
	ErrorCodeConversionError               ErrorCode = "CONVERSION_ERROR"
	ErrorCodeDocumentError                 ErrorCode = "DOCUMENT_ERROR"
	ErrorCodeSubmissionError               ErrorCode = "SUBMISSION_ERROR"
	ErrorCodeProcessingError               ErrorCode = "PROCESSING_ERROR"
	ErrorCodeAPIError                      ErrorCode = "API_ERROR"
	ErrorCodeNetworkError                  ErrorCode = "NETWORK_ERROR"
	ErrorCodeTimeoutError                  ErrorCode = "TIMEOUT_ERROR"
	ErrorCodeRateLimitExceeded             ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeInternalServerError           ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeServiceUnavailable            ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeDatabaseError                 ErrorCode = "DATABASE_ERROR"
	ErrorCodeQueueError                    ErrorCode = "QUEUE_ERROR"
	ErrorCodeGovernmentSystemUnavailable   ErrorCode = "GOVERNMENT_SYSTEM_UNAVAILABLE"
	ErrorCodeSubmissionTimeout             ErrorCode = "SUBMISSION_TIMEOUT"
	ErrorCodeCircuitBreakerOpen            ErrorCode = "CIRCUIT_BREAKER_OPEN"
	ErrorCodeMaxRetriesExceeded            ErrorCode = "MAX_RETRIES_EXCEEDED"
	ErrorCodeEmptyPayload                  ErrorCode = "EMPTY_PAYLOAD"
	ErrorCodeMalformedJSON                 ErrorCode = "MALFORMED_JSON"
	ErrorCodeInvalidPayloadFormat          ErrorCode = "INVALID_PAYLOAD_FORMAT"
)

// SubmissionStatus enumeration matching Python SDK
type SubmissionStatus string

const (
	SubmissionStatusPending    SubmissionStatus = "PENDING"
	SubmissionStatusProcessing SubmissionStatus = "PROCESSING"
	SubmissionStatusSubmitted  SubmissionStatus = "SUBMITTED"
	SubmissionStatusAccepted   SubmissionStatus = "ACCEPTED"
	SubmissionStatusRejected   SubmissionStatus = "REJECTED"
	SubmissionStatusFailed     SubmissionStatus = "FAILED"
	SubmissionStatusQueued     SubmissionStatus = "QUEUED"
)

// Source model matching Python SDK
type Source struct {
	Name    string      `json:"name"`
	Version string      `json:"version"`
	Type    *SourceType `json:"type,omitempty"`
}

// NewSource creates a new Source
func NewSource(name, version string, sourceType *SourceType) *Source {
	source := &Source{
		Name:    name,
		Version: version,
		Type:    sourceType,
	}
	
	if source.Name == "" {
		source.Name = ""
	}
	if source.Version == "" {
		source.Version = ""
	}
	
	return source
}

// GetIdentity Get identity string for the source
func (s *Source) GetIdentity() string {
	return fmt.Sprintf("%s:%s", s.Name, s.Version)
}

// GetID Legacy getter for backward compatibility
func (s *Source) GetID() string {
	return s.GetIdentity()
}

// GetType Get type as string
func (s *Source) GetType() string {
	if s.Type != nil {
		return string(*s.Type)
	}
	return ""
}

// GetSourceTypeEnum Get the actual SourceType enum
func (s *Source) GetSourceTypeEnum() *SourceType {
	return s.Type
}

// GetName getter for name
func (s *Source) GetName() string {
	return s.Name
}

// GetVersion getter for version
func (s *Source) GetVersion() string {
	return s.Version
}

// DestinationDetails model matching Python SDK
type DestinationDetails struct {
	Country       *string   `json:"country,omitempty"`
	Authority     *string   `json:"authority,omitempty"`
	DocumentType  *string   `json:"document_type,omitempty"`
	Recipients    *[]string `json:"recipients,omitempty"`
	Subject       *string   `json:"subject,omitempty"`
	Body          *string   `json:"body,omitempty"`
	ParticipantID *string   `json:"participant_id,omitempty"`
	ProcessID     *string   `json:"process_id,omitempty"`
}

// SetCountry setter for country
func (d *DestinationDetails) SetCountry(country string) {
	d.Country = &country
}

// SetAuthority setter for authority
func (d *DestinationDetails) SetAuthority(authority string) {
	d.Authority = &authority
}

// SetDocumentType setter for document type
func (d *DestinationDetails) SetDocumentType(documentType string) {
	d.DocumentType = &documentType
}

// SetRecipients setter for recipients
func (d *DestinationDetails) SetRecipients(recipients []string) {
	d.Recipients = &recipients
}

// SetSubject setter for subject
func (d *DestinationDetails) SetSubject(subject string) {
	d.Subject = &subject
}

// SetBody setter for body
func (d *DestinationDetails) SetBody(body string) {
	d.Body = &body
}

// SetParticipantID setter for participant ID
func (d *DestinationDetails) SetParticipantID(participantID string) {
	d.ParticipantID = &participantID
}

// SetProcessID setter for process ID
func (d *DestinationDetails) SetProcessID(processID string) {
	d.ProcessID = &processID
}

// Destination model matching Python SDK
type Destination struct {
	Type    DestinationType     `json:"type"`
	Details *DestinationDetails `json:"details"`
}

// NewTaxAuthorityDestination Create tax authority destination
func NewTaxAuthorityDestination(country, authority, documentType string) *Destination {
	details := &DestinationDetails{}
	details.SetCountry(country)
	details.SetAuthority(authority)
	details.SetDocumentType(documentType)
	return &Destination{
		Type:    DestinationTypeTaxAuthority,
		Details: details,
	}
}

// NewEmailDestination Create email destination
func NewEmailDestination(recipients []string, subject, body string) *Destination {
	details := &DestinationDetails{}
	details.SetRecipients(recipients)
	details.SetSubject(subject)
	details.SetBody(body)
	return &Destination{
		Type:    DestinationTypeEmail,
		Details: details,
	}
}

// NewArchiveDestination Create archive destination
func NewArchiveDestination() *Destination {
	return &Destination{
		Type:    DestinationTypeArchive,
		Details: &DestinationDetails{},
	}
}

// NewPeppolDestination Create PEPPOL destination
func NewPeppolDestination(participantID, processID, documentType string) *Destination {
	details := &DestinationDetails{}
	details.SetParticipantID(participantID)
	details.SetProcessID(processID)
	details.SetDocumentType(documentType)
	return &Destination{
		Type:    DestinationTypePeppol,
		Details: details,
	}
}

// GetType getter for type
func (d *Destination) GetType() DestinationType {
	return d.Type
}

// GetDetails getter for details
func (d *Destination) GetDetails() *DestinationDetails {
	return d.Details
}

// CircuitBreakerConfig model matching Python SDK
type CircuitBreakerConfig struct {
	FailureThreshold   int `json:"failure_threshold"`
	TimeoutDurationMs int `json:"timeout_duration_ms"`
}

// NewCircuitBreakerConfig creates a new circuit breaker config
func NewCircuitBreakerConfig(failureThreshold, timeoutDurationMs int) *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		FailureThreshold:   failureThreshold,
		TimeoutDurationMs: timeoutDurationMs,
	}
}

// GetFailureThreshold getter for failure threshold
func (c *CircuitBreakerConfig) GetFailureThreshold() int {
	return c.FailureThreshold
}

// GetTimeout getter for timeout
func (c *CircuitBreakerConfig) GetTimeout() int {
	return c.TimeoutDurationMs
}

// RetryConfig model matching Python SDK
type RetryConfig struct {
	MaxAttempts              int         `json:"max_attempts"`
	BaseDelayMs              int         `json:"base_delay_ms"`
	MaxDelayMs               int         `json:"max_delay_ms"`
	BackoffMultiplier        float64     `json:"backoff_multiplier"`
	JitterFactor             float64     `json:"jitter_factor"`
	RetryableErrors          []ErrorCode `json:"retryable_errors"`
	RetryableHTTPCodes       []int       `json:"retryable_http_codes"`
	CircuitBreakerEnabled    bool        `json:"circuit_breaker_enabled"`
	FailureThreshold         int         `json:"failure_threshold"`
	CircuitBreakerTimeoutMs int         `json:"circuit_breaker_timeout_ms"`
}

// NewDefaultRetryConfig Create default retry configuration
func NewDefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:              5,
		BaseDelayMs:              500,
		MaxDelayMs:               30000,
		BackoffMultiplier:        2.0,
		JitterFactor:             0.1,
		RetryableErrors: []ErrorCode{
			ErrorCodeNetworkError,
			ErrorCodeTimeoutError,
			ErrorCodeRateLimitExceeded,
			ErrorCodeInternalServerError,
			ErrorCodeServiceUnavailable,
		},
		RetryableHTTPCodes:       []int{408, 429, 500, 502, 503, 504},
		CircuitBreakerEnabled:    true,
		FailureThreshold:         3,
		CircuitBreakerTimeoutMs: 60000,
	}
}

// NewAggressiveRetryConfig Create aggressive retry configuration
func NewAggressiveRetryConfig() *RetryConfig {
	config := NewDefaultRetryConfig()
	config.MaxAttempts = 7
	config.BaseDelayMs = 200
	config.MaxDelayMs = 60000
	config.BackoffMultiplier = 1.5
	return config
}

// NewConservativeRetryConfig Create conservative retry configuration
func NewConservativeRetryConfig() *RetryConfig {
	config := NewDefaultRetryConfig()
	config.MaxAttempts = 3
	config.BaseDelayMs = 1000
	config.MaxDelayMs = 10000
	config.BackoffMultiplier = 2.5
	return config
}

// NewNoRetryConfig Create no-retry configuration
func NewNoRetryConfig() *RetryConfig {
	config := NewDefaultRetryConfig()
	config.MaxAttempts = 1
	return config
}

// GetCircuitBreakerConfig Get circuit breaker configuration
func (r *RetryConfig) GetCircuitBreakerConfig() *CircuitBreakerConfig {
	return NewCircuitBreakerConfig(r.FailureThreshold, r.CircuitBreakerTimeoutMs)
}

// ShouldRetry Check if error should be retried
func (r *RetryConfig) ShouldRetry(errorCode ErrorCode) bool {
	for _, retryableError := range r.RetryableErrors {
		if retryableError == errorCode {
			return true
		}
	}
	return false
}

// ShouldRetryHTTPCode Check if HTTP code should be retried
func (r *RetryConfig) ShouldRetryHTTPCode(httpCode int) bool {
	for _, retryableCode := range r.RetryableHTTPCodes {
		if retryableCode == httpCode {
			return true
		}
	}
	return false
}
