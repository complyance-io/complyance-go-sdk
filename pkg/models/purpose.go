package models

// Purpose represents the processing purpose
type Purpose string

const (
	// PurposeMapping represents field mapping purpose
	PurposeMapping Purpose = "MAPPING"

	// PurposeInvoicing represents invoicing purpose
	PurposeInvoicing Purpose = "INVOICING"

	// PurposeValidation represents validation purpose
	PurposeValidation Purpose = "VALIDATION"
)

// String returns the string representation of the purpose
func (p Purpose) String() string {
	return string(p)
}