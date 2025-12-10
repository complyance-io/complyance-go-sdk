package models

// DocumentType represents the type of document being processed
type DocumentType string

const (
	// DocumentTypeTaxInvoice represents a tax invoice document
	DocumentTypeTaxInvoice DocumentType = "TAX_INVOICE"

	// DocumentTypeCreditNote represents a credit note document
	DocumentTypeCreditNote DocumentType = "CREDIT_NOTE"

	// DocumentTypeDebitNote represents a debit note document
	DocumentTypeDebitNote DocumentType = "DEBIT_NOTE"
)

// String returns the string representation of the document type
func (dt DocumentType) String() string {
	return string(dt)
}