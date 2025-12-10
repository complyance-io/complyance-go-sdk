package models

// Mode represents the processing mode
type Mode string

const (
	// ModeDocuments represents document processing mode
	ModeDocuments Mode = "DOCUMENTS"

	// ModeTemplates represents template processing mode
	ModeTemplates Mode = "TEMPLATES"
)

// String returns the string representation of the mode
func (m Mode) String() string {
	return string(m)
}