package models

import (
	"errors"
	"fmt"
)

// UnifyRequest represents a request to the Complyance Unified API
type UnifyRequest struct {
	// Source is the source system for the document
	Source *Source `json:"source"`

	// DocumentType is the type of document being processed
	DocumentType DocumentType `json:"document_type"`

	// Country is the ISO country code for country-specific processing
	Country string `json:"country"`

	// Operation specifies single or batch processing
	Operation Operation `json:"operation"`

	// Mode specifies document or template processing
	Mode Mode `json:"mode"`

	// Purpose specifies the processing purpose
	Purpose Purpose `json:"purpose"`

	// Payload contains the document data
	Payload map[string]interface{} `json:"payload"`

	// Destinations specifies where to send the processed document
	Destinations []*Destination `json:"destinations,omitempty"`

	// Metadata contains additional request information
	Metadata *RequestMetadata `json:"metadata,omitempty"`

	// ValidationResults contains validation results if validation was performed
	ValidationResults *ValidationResults `json:"validation_results,omitempty"`
}

// Validate checks if the request is valid
func (r *UnifyRequest) Validate() error {
	if r.Source == nil {
		return errors.New("source is required")
	}

	if err := r.Source.Validate(); err != nil {
		return fmt.Errorf("invalid source: %w", err)
	}

	if r.DocumentType == "" {
		return errors.New("document type is required")
	}

	// Validate document type
	validDocTypes := []DocumentType{
		DocumentTypeTaxInvoice,
		DocumentTypeCreditNote,
		DocumentTypeDebitNote,
	}

	validDocType := false
	for _, dt := range validDocTypes {
		if r.DocumentType == dt {
			validDocType = true
			break
		}
	}

	if !validDocType {
		return errors.New("invalid document type: " + string(r.DocumentType))
	}

	if r.Country == "" {
		return errors.New("country is required")
	}

	if len(r.Country) != 2 {
		return errors.New("country must be a 2-letter ISO code")
	}

	if r.Operation == "" {
		return errors.New("operation is required")
	}

	// Validate operation
	validOps := []Operation{
		OperationSingle,
		OperationBatch,
	}

	validOp := false
	for _, op := range validOps {
		if r.Operation == op {
			validOp = true
			break
		}
	}

	if !validOp {
		return errors.New("invalid operation: " + string(r.Operation))
	}

	if r.Mode == "" {
		return errors.New("mode is required")
	}

	// Validate mode
	validModes := []Mode{
		ModeDocuments,
		ModeTemplates,
	}

	validMode := false
	for _, m := range validModes {
		if r.Mode == m {
			validMode = true
			break
		}
	}

	if !validMode {
		return errors.New("invalid mode: " + string(r.Mode))
	}

	if r.Purpose == "" {
		return errors.New("purpose is required")
	}

	// Validate purpose
	validPurposes := []Purpose{
		PurposeMapping,
		PurposeInvoicing,
		PurposeValidation,
	}

	validPurpose := false
	for _, p := range validPurposes {
		if r.Purpose == p {
			validPurpose = true
			break
		}
	}

	if !validPurpose {
		return errors.New("invalid purpose: " + string(r.Purpose))
	}

	if r.Payload == nil {
		return errors.New("payload is required")
	}

	// Validate destinations if present
	if len(r.Destinations) > 0 {
		for i, dest := range r.Destinations {
			if err := dest.Validate(); err != nil {
				return fmt.Errorf("invalid destination at index %d: %w", i, err)
			}
		}
	}

	return nil
}

// ValidateWithResults performs validation and returns detailed validation results
func (r *UnifyRequest) ValidateWithResults() *ValidationResults {
	results := NewValidationResults()

	// Validate source
	if r.Source == nil {
		results.AddError("source", "Source is required")
	} else {
		if r.Source.ID == "" {
			results.AddError("source.id", "Source ID is required")
		}
		if r.Source.Type == "" {
			results.AddError("source.type", "Source type is required")
		} else {
			validTypes := []SourceType{
				SourceTypeFirstParty,
				SourceTypeThirdParty,
				SourceTypeMarketplace,
			}

			valid := false
			for _, t := range validTypes {
				if r.Source.Type == t {
					valid = true
					break
				}
			}

			if !valid {
				results.AddError("source.type", "Invalid source type: "+string(r.Source.Type))
			}
		}
		if r.Source.Name == "" {
			results.AddError("source.name", "Source name is required")
		}
	}

	// Validate document type
	if r.DocumentType == "" {
		results.AddError("document_type", "Document type is required")
	} else {
		validDocTypes := []DocumentType{
			DocumentTypeTaxInvoice,
			DocumentTypeCreditNote,
			DocumentTypeDebitNote,
		}

		valid := false
		for _, dt := range validDocTypes {
			if r.DocumentType == dt {
				valid = true
				break
			}
		}

		if !valid {
			results.AddError("document_type", "Invalid document type: "+string(r.DocumentType))
		}
	}

	// Validate country
	if r.Country == "" {
		results.AddError("country", "Country is required")
	} else if len(r.Country) != 2 {
		results.AddError("country", "Country must be a 2-letter ISO code")
	}

	// Validate operation
	if r.Operation == "" {
		results.AddError("operation", "Operation is required")
	} else {
		validOps := []Operation{
			OperationSingle,
			OperationBatch,
		}

		valid := false
		for _, op := range validOps {
			if r.Operation == op {
				valid = true
				break
			}
		}

		if !valid {
			results.AddError("operation", "Invalid operation: "+string(r.Operation))
		}
	}

	// Validate mode
	if r.Mode == "" {
		results.AddError("mode", "Mode is required")
	} else {
		validModes := []Mode{
			ModeDocuments,
			ModeTemplates,
		}

		valid := false
		for _, m := range validModes {
			if r.Mode == m {
				valid = true
				break
			}
		}

		if !valid {
			results.AddError("mode", "Invalid mode: "+string(r.Mode))
		}
	}

	// Validate purpose
	if r.Purpose == "" {
		results.AddError("purpose", "Purpose is required")
	} else {
		validPurposes := []Purpose{
			PurposeMapping,
			PurposeInvoicing,
			PurposeValidation,
		}

		valid := false
		for _, p := range validPurposes {
			if r.Purpose == p {
				valid = true
				break
			}
		}

		if !valid {
			results.AddError("purpose", "Invalid purpose: "+string(r.Purpose))
		}
	}

	// Validate payload
	if r.Payload == nil {
		results.AddError("payload", "Payload is required")
	}

	// Validate destinations if present
	if len(r.Destinations) > 0 {
		for i, dest := range r.Destinations {
			if dest.Type == "" {
				results.AddError(fmt.Sprintf("destinations[%d].type", i), "Destination type is required")
			}
			if dest.Config == nil {
				results.AddError(fmt.Sprintf("destinations[%d].config", i), "Destination config is required")
			}
		}
	}

	r.ValidationResults = results
	return results
}

// NewUnifyRequest creates a new UnifyRequest with the provided values
func NewUnifyRequest(source *Source, documentType DocumentType, country string) *UnifyRequest {
	return &UnifyRequest{
		Source:       source,
		DocumentType: documentType,
		Country:      country,
		Operation:    OperationSingle,
		Mode:         ModeDocuments,
		Purpose:      PurposeInvoicing,
		Payload:      make(map[string]interface{}),
		Metadata:     NewRequestMetadata(),
	}
}

// WithOperation sets the operation
func (r *UnifyRequest) WithOperation(operation Operation) *UnifyRequest {
	r.Operation = operation
	return r
}

// WithMode sets the mode
func (r *UnifyRequest) WithMode(mode Mode) *UnifyRequest {
	r.Mode = mode
	return r
}

// WithPurpose sets the purpose
func (r *UnifyRequest) WithPurpose(purpose Purpose) *UnifyRequest {
	r.Purpose = purpose
	return r
}

// WithPayload sets the payload
func (r *UnifyRequest) WithPayload(payload map[string]interface{}) *UnifyRequest {
	r.Payload = payload
	return r
}

// AddPayloadField adds a single payload field
func (r *UnifyRequest) AddPayloadField(key string, value interface{}) *UnifyRequest {
	if r.Payload == nil {
		r.Payload = make(map[string]interface{})
	}
	r.Payload[key] = value
	return r
}

// AddDestination adds a destination
func (r *UnifyRequest) AddDestination(destinationType string, config map[string]interface{}) *UnifyRequest {
	destination := NewDestination(destinationType, config)
	r.Destinations = append(r.Destinations, destination)
	return r
}

// WithMetadata sets the metadata
func (r *UnifyRequest) WithMetadata(metadata *RequestMetadata) *UnifyRequest {
	r.Metadata = metadata
	return r
}