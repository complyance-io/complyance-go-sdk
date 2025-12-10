package models

import (
	"errors"
	"strings"
)

// Source represents a source system for document processing
type Source struct {
	// ID is the unique identifier for the source
	ID string `json:"id"`

	// Type is the type of source system
	Type SourceType `json:"type"`

	// Name is the display name of the source
	Name string `json:"name"`

	// Version is the version of the source system
	Version string `json:"version,omitempty"`

	// Metadata contains additional source-specific information
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Validate checks if the source is valid
func (s *Source) Validate() error {
	if s.ID == "" {
		return errors.New("source ID is required")
	}

	if s.Type == "" {
		return errors.New("source type is required")
	}

	if s.Name == "" {
		return errors.New("source name is required")
	}

	// Validate source type
	validTypes := []SourceType{
		SourceTypeFirstParty,
		SourceTypeThirdParty,
		SourceTypeMarketplace,
	}

	valid := false
	for _, t := range validTypes {
		if s.Type == t {
			valid = true
			break
		}
	}

	if !valid {
		return errors.New("invalid source type: " + string(s.Type))
	}

	return nil
}

// NewSource creates a new Source with the provided values
func NewSource(id string, sourceType SourceType, name string) *Source {
	return &Source{
		ID:   id,
		Type: sourceType,
		Name: name,
	}
}

// WithVersion adds a version to the source
func (s *Source) WithVersion(version string) *Source {
	s.Version = version
	return s
}

// WithMetadata adds metadata to the source
func (s *Source) WithMetadata(metadata map[string]interface{}) *Source {
	s.Metadata = metadata
	return s
}

// AddMetadata adds a single metadata key-value pair
func (s *Source) AddMetadata(key string, value interface{}) *Source {
	if s.Metadata == nil {
		s.Metadata = make(map[string]interface{})
	}
	s.Metadata[key] = value
	return s
}

// String returns a string representation of the source
func (s *Source) String() string {
	var sb strings.Builder
	sb.WriteString("Source{")
	sb.WriteString("ID=")
	sb.WriteString(s.ID)
	sb.WriteString(", Type=")
	sb.WriteString(string(s.Type))
	sb.WriteString(", Name=")
	sb.WriteString(s.Name)
	if s.Version != "" {
		sb.WriteString(", Version=")
		sb.WriteString(s.Version)
	}
	sb.WriteString("}")
	return sb.String()
}