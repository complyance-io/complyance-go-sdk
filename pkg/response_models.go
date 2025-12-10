/*
Response models for the Complyance SDK matching Python SDK exactly.
*/
package complyancesdk

import (
	"fmt"
	"strings"
	"time"
)

// ErrorDetail model matching Python SDK
type ErrorDetail struct {
	Code               *ErrorCode             `json:"code,omitempty"`
	Message            *string                `json:"message,omitempty"`
	Suggestion         *string                `json:"suggestion,omitempty"`
	DocumentationURL   *string                `json:"documentation_url,omitempty"`
	Field              *string                `json:"field,omitempty"`
	FieldValue         interface{}            `json:"field_value,omitempty"`
	Context            map[string]interface{} `json:"context,omitempty"`
	ValidationErrors   []map[string]string    `json:"validation_errors,omitempty"`
	Retryable          bool                   `json:"retryable"`
	RetryAfterSeconds  *int                   `json:"retry_after_seconds,omitempty"`
	Timestamp          *string                `json:"timestamp,omitempty"`
}

// NewErrorDetail creates a new ErrorDetail
func NewErrorDetail() *ErrorDetail {
	now := time.Now().UTC().Format(time.RFC3339)
	return &ErrorDetail{
		Context:          make(map[string]interface{}),
		ValidationErrors: []map[string]string{},
		Retryable:        false,
		Timestamp:        &now,
	}
}

// NewErrorDetailWithCode creates an ErrorDetail with code and message
func NewErrorDetailWithCode(code ErrorCode, message string) *ErrorDetail {
	detail := NewErrorDetail()
	detail.Code = &code
	detail.Message = &message
	detail.Retryable = detail.isRetryableByDefault(code)
	return detail
}

// isRetryableByDefault Check if error code is retryable by default
func (e *ErrorDetail) isRetryableByDefault(code ErrorCode) bool {
	retryableCodes := map[ErrorCode]bool{
		ErrorCodeNetworkError:                  true,
		ErrorCodeTimeoutError:                  true,
		ErrorCodeRateLimitExceeded:             true,
		ErrorCodeAPIError:                      true,
		ErrorCodeInternalServerError:           true,
		ErrorCodeServiceUnavailable:            true,
		ErrorCodeDatabaseError:                 true,
		ErrorCodeQueueError:                    true,
		ErrorCodeGovernmentSystemUnavailable:   true,
		ErrorCodeSubmissionTimeout:             true,
		ErrorCodeCircuitBreakerOpen:            true,
	}
	return retryableCodes[code]
}

// AddContextValue Add context value
func (e *ErrorDetail) AddContextValue(key string, value interface{}) {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
}

// AddValidationError Add validation error
func (e *ErrorDetail) AddValidationError(field, message, code string) {
	validationError := map[string]string{
		"field":   field,
		"message": message,
		"code":    code,
	}
	e.ValidationErrors = append(e.ValidationErrors, validationError)
}

// GetContextValue Get context value
func (e *ErrorDetail) GetContextValue(key string) interface{} {
	if e.Context == nil {
		return nil
	}
	return e.Context[key]
}

// NewAPIErrorDetail Create API error detail
func NewAPIErrorDetail(httpStatus int, responseBody string) *ErrorDetail {
	error := NewErrorDetailWithCode(ErrorCodeAPIError, fmt.Sprintf("API request failed with HTTP %d", httpStatus))
	error.AddContextValue("httpStatus", httpStatus)
	error.AddContextValue("responseBody", responseBody)
	error.Retryable = httpStatus >= 500 || httpStatus == 429
	return error
}

// GetCode getter for code
func (e *ErrorDetail) GetCode() *ErrorCode {
	return e.Code
}

// GetMessage getter for message
func (e *ErrorDetail) GetMessage() *string {
	return e.Message
}

// GetSuggestion getter for suggestion
func (e *ErrorDetail) GetSuggestion() *string {
	return e.Suggestion
}

// IsRetryable getter for retryable
func (e *ErrorDetail) IsRetryable() bool {
	return e.Retryable
}

// WithSuggestion adds a suggestion to the error detail and returns the error detail for chaining
func (e *ErrorDetail) WithSuggestion(suggestion string) *ErrorDetail {
	e.Suggestion = &suggestion
	return e
}

// String string representation
func (e *ErrorDetail) String() string {
	codeStr := "nil"
	if e.Code != nil {
		codeStr = string(*e.Code)
	}
	messageStr := "nil"
	if e.Message != nil {
		messageStr = *e.Message
	}
	return fmt.Sprintf("ErrorDetail{code=%s, message='%s', retryable=%t}", codeStr, messageStr, e.Retryable)
}

// SourceResponse model matching Python SDK
type SourceResponse struct {
	SourceID *string `json:"source_id,omitempty"`
	Sourceid *string `json:"sourceid,omitempty"` // API returns lowercase version
	Type     *string `json:"type,omitempty"`
	Name     *string `json:"name,omitempty"`
	Version  *string `json:"version,omitempty"`
	Created  bool    `json:"created"`
	ID       *string `json:"id,omitempty"`
}

// GetSourceID getter for source ID
func (s *SourceResponse) GetSourceID() *string {
	return s.SourceID
}

// GetSourceid getter for sourceid
func (s *SourceResponse) GetSourceid() *string {
	return s.Sourceid
}

// GetType getter for type
func (s *SourceResponse) GetType() *string {
	return s.Type
}

// GetName getter for name
func (s *SourceResponse) GetName() *string {
	return s.Name
}

// GetVersion getter for version
func (s *SourceResponse) GetVersion() *string {
	return s.Version
}

// IsCreated getter for created
func (s *SourceResponse) IsCreated() bool {
	return s.Created
}

// GetID getter for ID
func (s *SourceResponse) GetID() *string {
	return s.ID
}

// AnalysisResponse model matching Python SDK
type AnalysisResponse struct {
	HasNested bool     `json:"has_nested"`
	Keys      []string `json:"keys,omitempty"`
	Size      *int     `json:"size,omitempty"`
}

// IsHasNested getter for has nested
func (a *AnalysisResponse) IsHasNested() bool {
	return a.HasNested
}

// GetKeys getter for keys
func (a *AnalysisResponse) GetKeys() []string {
	return a.Keys
}

// GetSize getter for size
func (a *AnalysisResponse) GetSize() *int {
	return a.Size
}

// PayloadResponse model matching Python SDK
type PayloadResponse struct {
	PayloadID   *string           `json:"payload_id,omitempty"`
	DocumentType *string          `json:"document_type,omitempty"`
	Country     *string           `json:"country,omitempty"`
	Environment *string           `json:"environment,omitempty"`
	StoredAt    *string           `json:"stored_at,omitempty"`
	Analysis    *AnalysisResponse `json:"analysis,omitempty"`
}

// GetPayloadID getter for payload ID
func (p *PayloadResponse) GetPayloadID() *string {
	return p.PayloadID
}

// GetDocumentType getter for document type
func (p *PayloadResponse) GetDocumentType() *string {
	return p.DocumentType
}

// GetCountry getter for country
func (p *PayloadResponse) GetCountry() *string {
	return p.Country
}

// GetEnvironment getter for environment
func (p *PayloadResponse) GetEnvironment() *string {
	return p.Environment
}

// GetStoredAt getter for stored at
func (p *PayloadResponse) GetStoredAt() *string {
	return p.StoredAt
}

// GetAnalysis getter for analysis
func (p *PayloadResponse) GetAnalysis() *AnalysisResponse {
	return p.Analysis
}

// TemplateResponse model matching Python SDK
type TemplateResponse struct {
	TemplateID             *string `json:"template_id,omitempty"`
	TemplateName           *string `json:"template_name,omitempty"`
	MappingCompleted       bool    `json:"mapping_completed"`
	TotalMandatoryFields   *int    `json:"total_mandatory_fields,omitempty"`
	MappedMandatoryFields  *int    `json:"mapped_mandatory_fields,omitempty"`
	AIMappingApplied       *bool   `json:"ai_mapping_applied,omitempty"`
}

// GetTemplateID getter for template ID
func (t *TemplateResponse) GetTemplateID() *string {
	return t.TemplateID
}

// GetTemplateName getter for template name
func (t *TemplateResponse) GetTemplateName() *string {
	return t.TemplateName
}

// IsMappingCompleted getter for mapping completed
func (t *TemplateResponse) IsMappingCompleted() bool {
	return t.MappingCompleted
}

// GetTotalMandatoryFields getter for total mandatory fields
func (t *TemplateResponse) GetTotalMandatoryFields() *int {
	return t.TotalMandatoryFields
}

// GetMappedMandatoryFields getter for mapped mandatory fields
func (t *TemplateResponse) GetMappedMandatoryFields() *int {
	return t.MappedMandatoryFields
}

// GetAIMappingApplied getter for AI mapping applied
func (t *TemplateResponse) GetAIMappingApplied() *bool {
	return t.AIMappingApplied
}

// ConversionResponse model matching Python SDK
type ConversionResponse struct {
	Success        bool                   `json:"success"`
	GetsDocument   map[string]interface{} `json:"gets_document,omitempty"`
	ConversionTime *int                   `json:"conversion_time,omitempty"`
	Errors         []string               `json:"errors,omitempty"`
}

// IsSuccess getter for success
func (c *ConversionResponse) IsSuccess() bool {
	return c.Success
}

// GetGetsDocument getter for GETS document
func (c *ConversionResponse) GetGetsDocument() map[string]interface{} {
	return c.GetsDocument
}

// GetConversionTime getter for conversion time
func (c *ConversionResponse) GetConversionTime() *int {
	return c.ConversionTime
}

// GetErrors getter for errors
func (c *ConversionResponse) GetErrors() []string {
	return c.Errors
}

// DocumentResponse model matching Python SDK
type DocumentResponse struct {
	DocumentID *string                `json:"document_id,omitempty"`
	DocumentType *string              `json:"document_type,omitempty"`
	CreatedAt  *string                `json:"created_at,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Status     *string                `json:"status,omitempty"`
}

// GetDocumentID getter for document ID
func (d *DocumentResponse) GetDocumentID() *string {
	return d.DocumentID
}

// GetDocumentType getter for document type
func (d *DocumentResponse) GetDocumentType() *string {
	return d.DocumentType
}

// GetCreatedAt getter for created at
func (d *DocumentResponse) GetCreatedAt() *string {
	return d.CreatedAt
}

// GetMetadata getter for metadata
func (d *DocumentResponse) GetMetadata() map[string]interface{} {
	return d.Metadata
}

// GetStatus getter for status
func (d *DocumentResponse) GetStatus() *string {
	return d.Status
}

// ValidationErrorModel model matching Python SDK
type ValidationErrorModel struct {
	Method  *string  `json:"method,omitempty"`
	Message *string  `json:"message,omitempty"`
	Code    *string  `json:"code,omitempty"`
	Path    []string `json:"path,omitempty"`
}

// GetMethod getter for method
func (v *ValidationErrorModel) GetMethod() *string {
	return v.Method
}

// GetMessage getter for message
func (v *ValidationErrorModel) GetMessage() *string {
	return v.Message
}

// GetCode getter for code
func (v *ValidationErrorModel) GetCode() *string {
	return v.Code
}

// GetPath getter for path
func (v *ValidationErrorModel) GetPath() []string {
	return v.Path
}

// ValidationResponse model matching Python SDK
type ValidationResponse struct {
	OverallSuccess bool                      `json:"overall_success"`
	Methods        []string                  `json:"methods,omitempty"`
	Errors         []*ValidationErrorModel   `json:"errors,omitempty"`
	ValidatedAt    *string                   `json:"validated_at,omitempty"`
	Success        *bool                     `json:"success,omitempty"`
}

// IsOverallSuccess getter for overall success
func (v *ValidationResponse) IsOverallSuccess() bool {
	return v.OverallSuccess
}

// GetMethods getter for methods
func (v *ValidationResponse) GetMethods() []string {
	return v.Methods
}

// GetErrors getter for errors
func (v *ValidationResponse) GetErrors() []*ValidationErrorModel {
	return v.Errors
}

// GetValidatedAt getter for validated at
func (v *ValidationResponse) GetValidatedAt() *string {
	return v.ValidatedAt
}

// GetSuccess getter for success
func (v *ValidationResponse) GetSuccess() *bool {
	return v.Success
}

// SubmissionResponseData model matching Python SDK
type SubmissionResponseData struct {
	ClearanceStatus    *string `json:"clearance_status,omitempty"`
	UUID              *string `json:"uuid,omitempty"`
	Hash              *string `json:"hash,omitempty"`
	QRCode            *string `json:"qr_code,omitempty"`
	SubmissionNumber  *string `json:"submission_number,omitempty"`
}

// GetClearanceStatus getter for clearance status
func (s *SubmissionResponseData) GetClearanceStatus() *string {
	return s.ClearanceStatus
}

// GetUUID getter for UUID
func (s *SubmissionResponseData) GetUUID() *string {
	return s.UUID
}

// GetHash getter for hash
func (s *SubmissionResponseData) GetHash() *string {
	return s.Hash
}

// GetQRCode getter for QR code
func (s *SubmissionResponseData) GetQRCode() *string {
	return s.QRCode
}

// GetSubmissionNumber getter for submission number
func (s *SubmissionResponseData) GetSubmissionNumber() *string {
	return s.SubmissionNumber
}

// SubmissionError model matching Python SDK
type SubmissionError struct {
	Code    *string `json:"code,omitempty"`
	Message *string `json:"message,omitempty"`
}

// GetCode getter for code
func (s *SubmissionError) GetCode() *string {
	return s.Code
}

// GetMessage getter for message
func (s *SubmissionError) GetMessage() *string {
	return s.Message
}

// SubmissionResponse model matching Python SDK
type SubmissionResponse struct {
	SubmissionID       *string                     `json:"submission_id,omitempty"`
	Country           *string                     `json:"country,omitempty"`
	Authority         *string                     `json:"authority,omitempty"`
	Status            *string                     `json:"status,omitempty"`
	SubmittedAt       *string                     `json:"submitted_at,omitempty"`
	Response          *SubmissionResponseData     `json:"response,omitempty"`
	GovernmentResponse map[string]interface{}     `json:"government_response,omitempty"`
	Errors            []*SubmissionError          `json:"errors,omitempty"`
}

// IsAccepted Check if submission is accepted
func (s *SubmissionResponse) IsAccepted() bool {
	return s.Status != nil && *s.Status == "accepted"
}

// IsRejected Check if submission is rejected
func (s *SubmissionResponse) IsRejected() bool {
	return s.Status != nil && *s.Status == "rejected"
}

// IsFailed Check if submission failed
func (s *SubmissionResponse) IsFailed() bool {
	return s.Status != nil && *s.Status == "failed"
}

// IsSubmitted Check if submission was submitted
func (s *SubmissionResponse) IsSubmitted() bool {
	return s.Status != nil && *s.Status == "submitted"
}

// GetSubmissionID getter for submission ID
func (s *SubmissionResponse) GetSubmissionID() *string {
	return s.SubmissionID
}

// GetCountry getter for country
func (s *SubmissionResponse) GetCountry() *string {
	return s.Country
}

// GetAuthority getter for authority
func (s *SubmissionResponse) GetAuthority() *string {
	return s.Authority
}

// GetStatus getter for status
func (s *SubmissionResponse) GetStatus() *string {
	return s.Status
}

// GetSubmittedAt getter for submitted at
func (s *SubmissionResponse) GetSubmittedAt() *string {
	return s.SubmittedAt
}

// GetResponse getter for response
func (s *SubmissionResponse) GetResponse() *SubmissionResponseData {
	return s.Response
}

// GetGovernmentResponse getter for government response
func (s *SubmissionResponse) GetGovernmentResponse() map[string]interface{} {
	return s.GovernmentResponse
}

// GetErrors getter for errors
func (s *SubmissionResponse) GetErrors() []*SubmissionError {
	return s.Errors
}

// ProcessingResponse model matching Python SDK
type ProcessingResponse struct {
	Purpose               *string  `json:"purpose,omitempty"`
	CompletedSteps        []string `json:"completed_steps,omitempty"`
	TotalProcessingTime   *int     `json:"total_processing_time,omitempty"`
	CompletedAt           *string  `json:"completed_at,omitempty"`
	ProcessedAt           *string  `json:"processed_at,omitempty"`
	RequestID             *string  `json:"request_id,omitempty"`
	Status                *string  `json:"status,omitempty"`
}

// IsInvoicingPurpose check if invoicing purpose
func (p *ProcessingResponse) IsInvoicingPurpose() bool {
	return p.Purpose != nil && *p.Purpose == "invoicing"
}

// IsMappingPurpose check if mapping purpose
func (p *ProcessingResponse) IsMappingPurpose() bool {
	return p.Purpose != nil && *p.Purpose == "mapping"
}

// GetPurpose getter for purpose
func (p *ProcessingResponse) GetPurpose() *string {
	return p.Purpose
}

// GetCompletedSteps getter for completed steps
func (p *ProcessingResponse) GetCompletedSteps() []string {
	return p.CompletedSteps
}

// GetTotalProcessingTime getter for total processing time
func (p *ProcessingResponse) GetTotalProcessingTime() *int {
	return p.TotalProcessingTime
}

// GetCompletedAt getter for completed at
func (p *ProcessingResponse) GetCompletedAt() *string {
	return p.CompletedAt
}

// GetProcessedAt getter for processed at
func (p *ProcessingResponse) GetProcessedAt() *string {
	return p.ProcessedAt
}

// GetRequestID getter for request ID
func (p *ProcessingResponse) GetRequestID() *string {
	return p.RequestID
}

// GetStatus getter for status
func (p *ProcessingResponse) GetStatus() *string {
	return p.Status
}

// DestinationsResponse model matching Python SDK
type DestinationsResponse struct {
	Count  *int     `json:"count,omitempty"`
	Stored bool     `json:"stored"`
	Types  []string `json:"types,omitempty"`
	Valid  *int     `json:"valid,omitempty"`
}

// GetCount getter for count
func (d *DestinationsResponse) GetCount() *int {
	return d.Count
}

// IsStored getter for stored
func (d *DestinationsResponse) IsStored() bool {
	return d.Stored
}

// GetTypes getter for types
func (d *DestinationsResponse) GetTypes() []string {
	return d.Types
}

// GetValid getter for valid
func (d *DestinationsResponse) GetValid() *int {
	return d.Valid
}

// LogicalDocumentTypeResponse model matching Python SDK
type LogicalDocumentTypeResponse struct {
	OriginalType *string                `json:"original_type,omitempty"`
	MetaConfig   map[string]interface{} `json:"meta_config,omitempty"`
}

// GetOriginalType getter for original type
func (l *LogicalDocumentTypeResponse) GetOriginalType() *string {
	return l.OriginalType
}

// GetMetaConfig getter for meta config
func (l *LogicalDocumentTypeResponse) GetMetaConfig() map[string]interface{} {
	return l.MetaConfig
}

// MetaConfigFlags model matching Python SDK
type MetaConfigFlags struct {
	IsExport       *bool `json:"is_export,omitempty"`
	IsSelfBilled   *bool `json:"is_self_billed,omitempty"`
	IsThirdParty   *bool `json:"is_third_party,omitempty"`
	IsNominalSupply *bool `json:"is_nominal_supply,omitempty"`
	IsSummary      *bool `json:"is_summary,omitempty"`
}

// GetIsExport getter for is export
func (m *MetaConfigFlags) GetIsExport() *bool {
	return m.IsExport
}

// GetIsSelfBilled getter for is self billed
func (m *MetaConfigFlags) GetIsSelfBilled() *bool {
	return m.IsSelfBilled
}

// GetIsThirdParty getter for is third party
func (m *MetaConfigFlags) GetIsThirdParty() *bool {
	return m.IsThirdParty
}

// GetIsNominalSupply getter for is nominal supply
func (m *MetaConfigFlags) GetIsNominalSupply() *bool {
	return m.IsNominalSupply
}

// GetIsSummary getter for is summary
func (m *MetaConfigFlags) GetIsSummary() *bool {
	return m.IsSummary
}

// UnifyResponseData model matching Python SDK
type UnifyResponseData struct {
	Source               *SourceResponse              `json:"source,omitempty"`
	Payload              *PayloadResponse             `json:"payload,omitempty"`
	Template             *TemplateResponse            `json:"template,omitempty"`
	LogicalDocumentType  *LogicalDocumentTypeResponse `json:"logical_document_type,omitempty"`
	Conversion           *ConversionResponse          `json:"conversion,omitempty"`
	Document             *DocumentResponse            `json:"document,omitempty"`
	Validation           *ValidationResponse          `json:"validation,omitempty"`
	Submission           *SubmissionResponse          `json:"submission,omitempty"`
	Processing           *ProcessingResponse          `json:"processing,omitempty"`
	Destinations         *DestinationsResponse        `json:"destinations,omitempty"`
}

// GetSource getter for source
func (u *UnifyResponseData) GetSource() *SourceResponse {
	return u.Source
}

// GetPayload getter for payload
func (u *UnifyResponseData) GetPayload() *PayloadResponse {
	return u.Payload
}

// GetTemplate getter for template
func (u *UnifyResponseData) GetTemplate() *TemplateResponse {
	return u.Template
}

// GetLogicalDocumentType getter for logical document type
func (u *UnifyResponseData) GetLogicalDocumentType() *LogicalDocumentTypeResponse {
	return u.LogicalDocumentType
}

// GetConversion getter for conversion
func (u *UnifyResponseData) GetConversion() *ConversionResponse {
	return u.Conversion
}

// GetDocument getter for document
func (u *UnifyResponseData) GetDocument() *DocumentResponse {
	return u.Document
}

// GetValidation getter for validation
func (u *UnifyResponseData) GetValidation() *ValidationResponse {
	return u.Validation
}

// GetSubmission getter for submission
func (u *UnifyResponseData) GetSubmission() *SubmissionResponse {
	return u.Submission
}

// GetProcessing getter for processing
func (u *UnifyResponseData) GetProcessing() *ProcessingResponse {
	return u.Processing
}

// GetDestinations getter for destinations
func (u *UnifyResponseData) GetDestinations() *DestinationsResponse {
	return u.Destinations
}

// UnifyResponse model matching Python SDK
type UnifyResponse struct {
	Status   string                 `json:"status"`
	Message  *string                `json:"message,omitempty"`
	Data     *UnifyResponseData     `json:"data,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Error    *ErrorDetail           `json:"error,omitempty"`
}

// IsSuccess Check if response indicates success
func (u *UnifyResponse) IsSuccess() bool {
	return strings.ToLower(u.Status) == "success"
}

// HasError Check if response has error
func (u *UnifyResponse) HasError() bool {
	return u.Error != nil || strings.ToLower(u.Status) == "error"
}

// GetStatus getter for status
func (u *UnifyResponse) GetStatus() string {
	return u.Status
}

// GetMessage getter for message
func (u *UnifyResponse) GetMessage() *string {
	return u.Message
}

// GetData getter for data
func (u *UnifyResponse) GetData() *UnifyResponseData {
	return u.Data
}

// GetMetadata getter for metadata
func (u *UnifyResponse) GetMetadata() map[string]interface{} {
	return u.Metadata
}

// GetError getter for error
func (u *UnifyResponse) GetError() *ErrorDetail {
	return u.Error
}

// SetStatus setter for status
func (u *UnifyResponse) SetStatus(status string) {
	u.Status = status
}

// SetMessage setter for message
func (u *UnifyResponse) SetMessage(message string) {
	u.Message = &message
}

// SetData setter for data
func (u *UnifyResponse) SetData(data *UnifyResponseData) {
	u.Data = data
}

// SetError setter for error
func (u *UnifyResponse) SetError(error *ErrorDetail) {
	u.Error = error
}
