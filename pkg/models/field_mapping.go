package models

import (
	"errors"
	"fmt"
	"strings"
)

// FieldMapping represents a mapping between source and target fields
type FieldMapping struct {
	// SourcePath is the JSON path to the source field
	SourcePath string `json:"source_path"`

	// TargetPath is the JSON path to the target field
	TargetPath string `json:"target_path"`

	// Transformation is an optional transformation to apply
	Transformation string `json:"transformation,omitempty"`

	// DefaultValue is an optional default value if source is missing
	DefaultValue interface{} `json:"default_value,omitempty"`

	// Required indicates if the field is required
	Required bool `json:"required"`

	// Description is an optional description of the mapping
	Description string `json:"description,omitempty"`
}

// FieldMappingSet represents a collection of field mappings
type FieldMappingSet struct {
	// Name is the name of the mapping set
	Name string `json:"name"`

	// Description is an optional description of the mapping set
	Description string `json:"description,omitempty"`

	// Country is the target country for the mapping
	Country string `json:"country"`

	// DocumentType is the target document type for the mapping
	DocumentType DocumentType `json:"document_type"`

	// Mappings is the list of field mappings
	Mappings []*FieldMapping `json:"mappings"`
}

// NewFieldMapping creates a new FieldMapping with the provided values
func NewFieldMapping(sourcePath string, targetPath string) *FieldMapping {
	return &FieldMapping{
		SourcePath: sourcePath,
		TargetPath: targetPath,
		Required:   false,
	}
}

// WithTransformation adds a transformation to the field mapping
func (fm *FieldMapping) WithTransformation(transformation string) *FieldMapping {
	fm.Transformation = transformation
	return fm
}

// WithDefaultValue adds a default value to the field mapping
func (fm *FieldMapping) WithDefaultValue(defaultValue interface{}) *FieldMapping {
	fm.DefaultValue = defaultValue
	return fm
}

// WithRequired sets the required flag on the field mapping
func (fm *FieldMapping) WithRequired(required bool) *FieldMapping {
	fm.Required = required
	return fm
}

// WithDescription adds a description to the field mapping
func (fm *FieldMapping) WithDescription(description string) *FieldMapping {
	fm.Description = description
	return fm
}

// Validate checks if the field mapping is valid
func (fm *FieldMapping) Validate() error {
	if fm.SourcePath == "" {
		return errors.New("source path is required")
	}

	if fm.TargetPath == "" {
		return errors.New("target path is required")
	}

	return nil
}

// String returns a string representation of the field mapping
func (fm *FieldMapping) String() string {
	var sb strings.Builder
	sb.WriteString("FieldMapping{")
	sb.WriteString("SourcePath=")
	sb.WriteString(fm.SourcePath)
	sb.WriteString(", TargetPath=")
	sb.WriteString(fm.TargetPath)
	if fm.Transformation != "" {
		sb.WriteString(", Transformation=")
		sb.WriteString(fm.Transformation)
	}
	if fm.DefaultValue != nil {
		sb.WriteString(", DefaultValue=")
		sb.WriteString(fmt.Sprintf("%v", fm.DefaultValue))
	}
	sb.WriteString(", Required=")
	sb.WriteString(fmt.Sprintf("%v", fm.Required))
	sb.WriteString("}")
	return sb.String()
}

// NewFieldMappingSet creates a new FieldMappingSet with the provided values
func NewFieldMappingSet(name string, country string, documentType DocumentType) *FieldMappingSet {
	return &FieldMappingSet{
		Name:         name,
		Country:      country,
		DocumentType: documentType,
		Mappings:     make([]*FieldMapping, 0),
	}
}

// WithDescription adds a description to the field mapping set
func (fms *FieldMappingSet) WithDescription(description string) *FieldMappingSet {
	fms.Description = description
	return fms
}

// AddMapping adds a field mapping to the set
func (fms *FieldMappingSet) AddMapping(mapping *FieldMapping) *FieldMappingSet {
	fms.Mappings = append(fms.Mappings, mapping)
	return fms
}

// AddSimpleMapping adds a simple field mapping to the set
func (fms *FieldMappingSet) AddSimpleMapping(sourcePath string, targetPath string, required bool) *FieldMappingSet {
	mapping := NewFieldMapping(sourcePath, targetPath).WithRequired(required)
	fms.Mappings = append(fms.Mappings, mapping)
	return fms
}

// Validate checks if the field mapping set is valid
func (fms *FieldMappingSet) Validate() error {
	if fms.Name == "" {
		return errors.New("name is required")
	}

	if fms.Country == "" {
		return errors.New("country is required")
	}

	if len(fms.Country) != 2 {
		return errors.New("country must be a 2-letter ISO code")
	}

	if fms.DocumentType == "" {
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
		if fms.DocumentType == dt {
			validDocType = true
			break
		}
	}

	if !validDocType {
		return errors.New("invalid document type: " + string(fms.DocumentType))
	}

	if len(fms.Mappings) == 0 {
		return errors.New("at least one mapping is required")
	}

	// Validate individual mappings
	for i, mapping := range fms.Mappings {
		if err := mapping.Validate(); err != nil {
			return fmt.Errorf("invalid mapping at index %d: %w", i, err)
		}
	}

	return nil
}

// GetRequiredMappings returns all required field mappings
func (fms *FieldMappingSet) GetRequiredMappings() []*FieldMapping {
	required := make([]*FieldMapping, 0)
	for _, mapping := range fms.Mappings {
		if mapping.Required {
			required = append(required, mapping)
		}
	}
	return required
}

// GetOptionalMappings returns all optional field mappings
func (fms *FieldMappingSet) GetOptionalMappings() []*FieldMapping {
	optional := make([]*FieldMapping, 0)
	for _, mapping := range fms.Mappings {
		if !mapping.Required {
			optional = append(optional, mapping)
		}
	}
	return optional
}

// GetMappingByTargetPath returns the mapping for the specified target path
func (fms *FieldMappingSet) GetMappingByTargetPath(targetPath string) *FieldMapping {
	for _, mapping := range fms.Mappings {
		if mapping.TargetPath == targetPath {
			return mapping
		}
	}
	return nil
}

// GetMappingBySourcePath returns the mapping for the specified source path
func (fms *FieldMappingSet) GetMappingBySourcePath(sourcePath string) *FieldMapping {
	for _, mapping := range fms.Mappings {
		if mapping.SourcePath == sourcePath {
			return mapping
		}
	}
	return nil
}