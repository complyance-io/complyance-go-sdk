/*
Request models for the Complyance SDK matching Python SDK exactly.
*/
package complyancesdk

import (
	"math/rand"
	"strconv"
	"time"
)

// UnifyRequest model matching Python SDK
type UnifyRequest struct {
	Source             *Source        `json:"source"`
	DocumentType       DocumentType   `json:"document_type"`
	DocumentTypeString *string        `json:"document_type_string,omitempty"`
	Country            string         `json:"country"`
	Operation          *Operation     `json:"operation,omitempty"`
	Mode               *Mode          `json:"mode,omitempty"`
	Purpose            *Purpose       `json:"purpose,omitempty"`
	Payload            map[string]interface{} `json:"payload,omitempty"`
	APIKey             *string        `json:"api_key,omitempty"`
	RequestID          *string        `json:"request_id,omitempty"`
	Timestamp          *string        `json:"timestamp,omitempty"`
	Env                *string        `json:"env,omitempty"`
	Destinations       []*Destination `json:"destinations,omitempty"`
	CorrelationID      *string        `json:"correlation_id,omitempty"`
}

// NewUnifyRequest creates a new UnifyRequest
func NewUnifyRequest() *UnifyRequest {
	now := time.Now().UTC().Format(time.RFC3339)
	requestID := "req_" + strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10) + "_" + strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
	
	return &UnifyRequest{
		Payload:      make(map[string]interface{}),
		Destinations: []*Destination{},
		RequestID:    &requestID,
		Timestamp:    &now,
	}
}

// NewUnifyRequestBuilder Create builder for UnifyRequest
func NewUnifyRequestBuilder() *UnifyRequestBuilder {
	return &UnifyRequestBuilder{
		country:      "",
		payload:      make(map[string]interface{}),
		destinations: []*Destination{},
	}
}

// GetSource getter for source
func (u *UnifyRequest) GetSource() *Source {
	return u.Source
}

// GetDocumentType getter for document type
func (u *UnifyRequest) GetDocumentType() DocumentType {
	return u.DocumentType
}

// GetDocumentTypeString getter for document type string
func (u *UnifyRequest) GetDocumentTypeString() *string {
	return u.DocumentTypeString
}

// GetCountry getter for country
func (u *UnifyRequest) GetCountry() string {
	return u.Country
}

// GetOperation getter for operation
func (u *UnifyRequest) GetOperation() *Operation {
	return u.Operation
}

// GetMode getter for mode
func (u *UnifyRequest) GetMode() *Mode {
	return u.Mode
}

// GetPurpose getter for purpose
func (u *UnifyRequest) GetPurpose() *Purpose {
	return u.Purpose
}

// GetPayload getter for payload
func (u *UnifyRequest) GetPayload() map[string]interface{} {
	return u.Payload
}

// GetAPIKey getter for API key
func (u *UnifyRequest) GetAPIKey() *string {
	return u.APIKey
}

// GetRequestID getter for request ID
func (u *UnifyRequest) GetRequestID() *string {
	return u.RequestID
}

// GetTimestamp getter for timestamp
func (u *UnifyRequest) GetTimestamp() *string {
	return u.Timestamp
}

// GetEnv getter for env
func (u *UnifyRequest) GetEnv() *string {
	return u.Env
}

// GetDestinations getter for destinations
func (u *UnifyRequest) GetDestinations() []*Destination {
	return u.Destinations
}

// GetCorrelationID getter for correlation ID
func (u *UnifyRequest) GetCorrelationID() *string {
	return u.CorrelationID
}

// SetSource setter for source
func (u *UnifyRequest) SetSource(source *Source) {
	u.Source = source
}

// SetDocumentType setter for document type
func (u *UnifyRequest) SetDocumentType(documentType DocumentType) {
	u.DocumentType = documentType
}

// SetDocumentTypeString setter for document type string
func (u *UnifyRequest) SetDocumentTypeString(documentTypeString string) {
	u.DocumentTypeString = &documentTypeString
}

// SetCountry setter for country
func (u *UnifyRequest) SetCountry(country string) {
	u.Country = country
}

// SetOperation setter for operation
func (u *UnifyRequest) SetOperation(operation Operation) {
	u.Operation = &operation
}

// SetMode setter for mode
func (u *UnifyRequest) SetMode(mode Mode) {
	u.Mode = &mode
}

// SetPurpose setter for purpose
func (u *UnifyRequest) SetPurpose(purpose Purpose) {
	u.Purpose = &purpose
}

// SetPayload setter for payload
func (u *UnifyRequest) SetPayload(payload map[string]interface{}) {
	u.Payload = payload
}

// SetAPIKey setter for API key
func (u *UnifyRequest) SetAPIKey(apiKey string) {
	u.APIKey = &apiKey
}

// SetRequestID setter for request ID
func (u *UnifyRequest) SetRequestID(requestID string) {
	u.RequestID = &requestID
}

// SetTimestamp setter for timestamp
func (u *UnifyRequest) SetTimestamp(timestamp string) {
	u.Timestamp = &timestamp
}

// SetEnv setter for env
func (u *UnifyRequest) SetEnv(env string) {
	u.Env = &env
}

// SetDestinations setter for destinations
func (u *UnifyRequest) SetDestinations(destinations []*Destination) {
	u.Destinations = destinations
}

// SetCorrelationID setter for correlation ID
func (u *UnifyRequest) SetCorrelationID(correlationID string) {
	u.CorrelationID = &correlationID
}

// UnifyRequestBuilder Builder for UnifyRequest matching Python SDK
type UnifyRequestBuilder struct {
	source             *Source
	documentType       *DocumentType
	documentTypeString *string
	country            string
	operation          *Operation
	mode               *Mode
	purpose            *Purpose
	payload            map[string]interface{}
	apiKey             *string
	requestID          *string
	timestamp          *string
	env                *string
	destinations       []*Destination
	correlationID      *string
}

// Source setter for source
func (b *UnifyRequestBuilder) Source(source *Source) *UnifyRequestBuilder {
	b.source = source
	return b
}

// DocumentType setter for document type
func (b *UnifyRequestBuilder) DocumentType(documentType DocumentType) *UnifyRequestBuilder {
	b.documentType = &documentType
	return b
}

// DocumentTypeString setter for document type string
func (b *UnifyRequestBuilder) DocumentTypeString(documentTypeString string) *UnifyRequestBuilder {
	b.documentTypeString = &documentTypeString
	return b
}

// Country setter for country
func (b *UnifyRequestBuilder) Country(country string) *UnifyRequestBuilder {
	b.country = country
	return b
}

// Operation setter for operation
func (b *UnifyRequestBuilder) Operation(operation Operation) *UnifyRequestBuilder {
	b.operation = &operation
	return b
}

// Mode setter for mode
func (b *UnifyRequestBuilder) Mode(mode Mode) *UnifyRequestBuilder {
	b.mode = &mode
	return b
}

// Purpose setter for purpose
func (b *UnifyRequestBuilder) Purpose(purpose Purpose) *UnifyRequestBuilder {
	b.purpose = &purpose
	return b
}

// Payload setter for payload
func (b *UnifyRequestBuilder) Payload(payload map[string]interface{}) *UnifyRequestBuilder {
	b.payload = payload
	return b
}

// APIKey setter for API key
func (b *UnifyRequestBuilder) APIKey(apiKey string) *UnifyRequestBuilder {
	b.apiKey = &apiKey
	return b
}

// RequestID setter for request ID
func (b *UnifyRequestBuilder) RequestID(requestID string) *UnifyRequestBuilder {
	b.requestID = &requestID
	return b
}

// Timestamp setter for timestamp
func (b *UnifyRequestBuilder) Timestamp(timestamp string) *UnifyRequestBuilder {
	b.timestamp = &timestamp
	return b
}

// Env setter for env
func (b *UnifyRequestBuilder) Env(env string) *UnifyRequestBuilder {
	b.env = &env
	return b
}

// Destinations setter for destinations
func (b *UnifyRequestBuilder) Destinations(destinations []*Destination) *UnifyRequestBuilder {
	b.destinations = destinations
	return b
}

// CorrelationID setter for correlation ID
func (b *UnifyRequestBuilder) CorrelationID(correlationID string) *UnifyRequestBuilder {
	b.correlationID = &correlationID
	return b
}

// Build builds the UnifyRequest
func (b *UnifyRequestBuilder) Build() *UnifyRequest {
	request := NewUnifyRequest()
	
	request.Source = b.source
	if b.documentType != nil {
		request.DocumentType = *b.documentType
	}
	request.DocumentTypeString = b.documentTypeString
	request.Country = b.country
	request.Operation = b.operation
	request.Mode = b.mode
	request.Purpose = b.purpose
	request.Payload = b.payload
	request.APIKey = b.apiKey
	if b.requestID != nil {
		request.RequestID = b.requestID
	}
	if b.timestamp != nil {
		request.Timestamp = b.timestamp
	}
	request.Env = b.env
	request.Destinations = b.destinations
	request.CorrelationID = b.correlationID
	
	return request
}

// SubmissionResponseOld Legacy submission response for backward compatibility
type SubmissionResponseOld struct {
	SubmissionID string            `json:"submission_id"`
	Status       SubmissionStatus  `json:"status"`
	Error        *ErrorDetail      `json:"error,omitempty"`
}

// GetSubmissionID getter for submission ID
func (s *SubmissionResponseOld) GetSubmissionID() string {
	return s.SubmissionID
}

// GetStatus getter for status
func (s *SubmissionResponseOld) GetStatus() SubmissionStatus {
	return s.Status
}

// GetError getter for error
func (s *SubmissionResponseOld) GetError() *ErrorDetail {
	return s.Error
}

// SourceRef model matching Python SDK
type SourceRef struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// NewSourceRef creates a new SourceRef
func NewSourceRef(name, version string) *SourceRef {
	return &SourceRef{
		Name:    name,
		Version: version,
	}
}

// GetName getter for name
func (s *SourceRef) GetName() string {
	return s.Name
}

// GetVersion getter for version
func (s *SourceRef) GetVersion() string {
	return s.Version
}

// PayloadSubmission model for queue matching Python SDK
type PayloadSubmission struct {
	Payload      string       `json:"payload"`
	Source       *Source      `json:"source"`
	Country      Country      `json:"country"`
	DocumentType DocumentType `json:"document_type"`
}

// NewPayloadSubmission creates a new PayloadSubmission
func NewPayloadSubmission(payload string, source *Source, country Country, documentType DocumentType) *PayloadSubmission {
	return &PayloadSubmission{
		Payload:      payload,
		Source:       source,
		Country:      country,
		DocumentType: documentType,
	}
}

// GetPayload getter for payload
func (p *PayloadSubmission) GetPayload() string {
	return p.Payload
}

// GetSource getter for source
func (p *PayloadSubmission) GetSource() *Source {
	return p.Source
}

// GetCountry getter for country
func (p *PayloadSubmission) GetCountry() Country {
	return p.Country
}

// GetDocumentType getter for document type
func (p *PayloadSubmission) GetDocumentType() DocumentType {
	return p.DocumentType
}

// PolicyResult model for logical document types matching Python SDK
type PolicyResult struct {
	BaseType        DocumentType           `json:"base_type"`
	DocumentType    string                 `json:"document_type"`
	MetaConfigFlags map[string]interface{} `json:"meta_config_flags"`
}

// NewPolicyResult creates a new PolicyResult
func NewPolicyResult(baseType DocumentType, documentType string, metaConfigFlags map[string]interface{}) *PolicyResult {
	return &PolicyResult{
		BaseType:        baseType,
		DocumentType:    documentType,
		MetaConfigFlags: metaConfigFlags,
	}
}

// GetBaseType getter for base type
func (p *PolicyResult) GetBaseType() DocumentType {
	return p.BaseType
}

// GetDocumentType getter for document type
func (p *PolicyResult) GetDocumentType() string {
	return p.DocumentType
}

// GetMetaConfigFlags getter for meta config flags
func (p *PolicyResult) GetMetaConfigFlags() map[string]interface{} {
	return p.MetaConfigFlags
}
