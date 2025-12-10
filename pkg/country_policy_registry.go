/*
Country Policy Registry implementation matching Python SDK exactly.
*/
package complyancesdk

import (
	"strings"
)

// CountryPolicyRegistry Country policy registry matching Python SDK
type CountryPolicyRegistry struct{}

// Evaluate country policy and return policy result
func (c *CountryPolicyRegistry) Evaluate(country Country, logicalType LogicalDocType) *PolicyResult {
	// Default base type mapping
	baseType := DocumentTypeTaxInvoice
	documentType := string(logicalType)
	metaConfigFlags := make(map[string]interface{})

	// Map logical document types to base document types
	logicalName := string(logicalType)
	if strings.Contains(logicalName, "CREDIT_NOTE") {
		if strings.Contains(logicalName, "SIMPLIFIED") {
			baseType = DocumentTypeSimplifiedCreditNote
		} else {
			baseType = DocumentTypeCreditNote
		}
	} else if strings.Contains(logicalName, "DEBIT_NOTE") {
		if strings.Contains(logicalName, "SIMPLIFIED") {
			baseType = DocumentTypeSimplifiedDebitNote
		} else {
			baseType = DocumentTypeDebitNote
		}
	} else if strings.Contains(logicalName, "SIMPLIFIED") {
		baseType = DocumentTypeSimplifiedInvoice
	} else {
		baseType = DocumentTypeTaxInvoice
	}

	// Set meta config flags based on logical type
	if strings.Contains(logicalName, "EXPORT") {
		metaConfigFlags["isExport"] = true
	} else {
		metaConfigFlags["isExport"] = false
	}

	if strings.Contains(logicalName, "SELF_BILLED") {
		metaConfigFlags["isSelfBilled"] = true
	} else {
		metaConfigFlags["isSelfBilled"] = false
	}

	if strings.Contains(logicalName, "THIRD_PARTY") {
		metaConfigFlags["isThirdParty"] = true
	} else {
		metaConfigFlags["isThirdParty"] = false
	}

	if strings.Contains(logicalName, "NOMINAL_SUPPLY") {
		metaConfigFlags["isNominal"] = true
	} else {
		metaConfigFlags["isNominal"] = false
	}

	if strings.Contains(logicalName, "SUMMARY") {
		metaConfigFlags["isSummary"] = true
	} else {
		metaConfigFlags["isSummary"] = false
	}

	// Set B2B/B2C flags based on logical type
	if strings.Contains(logicalName, "SIMPLIFIED") {
		metaConfigFlags["isB2B"] = false
	} else {
		metaConfigFlags["isB2B"] = true
	}

	// Set default flags
	metaConfigFlags["isPrepayment"] = false
	metaConfigFlags["isAdjusted"] = false
	metaConfigFlags["isReceipt"] = false

	// Country-specific adjustments
	switch country {
	case CountrySA:
		// Saudi Arabia specific logic
		documentType = c.getSaudiDocumentType(logicalType)
	case CountryMY:
		// Malaysia specific logic
		documentType = c.getMalaysiaDocumentType(logicalType)
	case CountryAE:
		// UAE specific logic
		documentType = c.getUAEDocumentType(logicalType)
	case CountrySG:
		// Singapore specific logic
		documentType = c.getSingaporeDocumentType(logicalType)
	}

	return NewPolicyResult(baseType, documentType, metaConfigFlags)
}

// getSaudiDocumentType Get Saudi-specific document type
func (c *CountryPolicyRegistry) getSaudiDocumentType(logicalType LogicalDocType) string {
	switch logicalType {
	case LogicalDocTypeTaxInvoice:
		return "tax_invoice"
	case LogicalDocTypeSimplifiedTaxInvoice:
		return "tax_invoice"
	case LogicalDocTypeTaxInvoiceCreditNote:
		return "credit_note"
	case LogicalDocTypeSimplifiedTaxInvoiceCreditNote:
		return "credit_note"
	case LogicalDocTypeTaxInvoiceDebitNote:
		return "debit_note"
	case LogicalDocTypeSimplifiedTaxInvoiceDebitNote:
		return "debit_note"
	default:
		return string(logicalType)
	}
}

// getMalaysiaDocumentType Get Malaysia-specific document type
func (c *CountryPolicyRegistry) getMalaysiaDocumentType(logicalType LogicalDocType) string {
	// Malaysia typically uses tax invoices
	logicalName := string(logicalType)
	if strings.Contains(logicalName, "CREDIT_NOTE") {
		return "credit_note"
	} else if strings.Contains(logicalName, "DEBIT_NOTE") {
		return "debit_note"
	} else {
		return "tax_invoice"
	}
}

// getUAEDocumentType Get UAE-specific document type
func (c *CountryPolicyRegistry) getUAEDocumentType(logicalType LogicalDocType) string {
	// UAE follows similar patterns to Saudi Arabia
	return c.getSaudiDocumentType(logicalType)
}

// getSingaporeDocumentType Get Singapore-specific document type
func (c *CountryPolicyRegistry) getSingaporeDocumentType(logicalType LogicalDocType) string {
	// Singapore typically uses tax invoices
	logicalName := string(logicalType)
	if strings.Contains(logicalName, "CREDIT_NOTE") {
		return "credit_note"
	} else if strings.Contains(logicalName, "DEBIT_NOTE") {
		return "debit_note"
	} else {
		return "tax_invoice"
	}
}

// Global registry instance
var CountryPolicyRegistryInstance = &CountryPolicyRegistry{}
