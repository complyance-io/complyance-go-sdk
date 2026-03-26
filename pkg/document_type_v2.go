package complyancesdk

import "strings"

type GetsDocumentBase string

const (
	GetsDocumentBaseTaxInvoice        GetsDocumentBase = "tax_invoice"
	GetsDocumentBaseSimplifiedInvoice GetsDocumentBase = "simplified_invoice"
	GetsDocumentBaseCreditNote        GetsDocumentBase = "credit_note"
	GetsDocumentBaseDebitNote         GetsDocumentBase = "debit_note"
)

type GetsDocumentModifier string

const (
	GetsDocumentModifierB2B           GetsDocumentModifier = "b2b"
	GetsDocumentModifierB2C           GetsDocumentModifier = "b2c"
	GetsDocumentModifierExport        GetsDocumentModifier = "export"
	GetsDocumentModifierSelfBilled    GetsDocumentModifier = "self_billed"
	GetsDocumentModifierThirdParty    GetsDocumentModifier = "third_party"
	GetsDocumentModifierNominalSupply GetsDocumentModifier = "nominal_supply"
	GetsDocumentModifierSummary       GetsDocumentModifier = "summary"
	GetsDocumentModifierB2G           GetsDocumentModifier = "b2g"
)

type GetsDocumentVariant string

const (
	GetsDocumentVariantStandard GetsDocumentVariant = "standard"
)

type GetsDocumentTypeV2 struct {
	Base      string   `json:"base"`
	Modifiers []string `json:"modifiers,omitempty"`
	Variant   *string  `json:"variant,omitempty"`
}

func NewGetsDocumentTypeV2(base string, modifiers []string, variant *string) *GetsDocumentTypeV2 {
	normalizedBase := strings.ToLower(strings.TrimSpace(base))

	normalizedModifiers := make([]string, 0, len(modifiers))
	seen := map[string]bool{}
	for _, modifier := range modifiers {
		value := strings.ToLower(strings.TrimSpace(modifier))
		if value == "" || seen[value] {
			continue
		}
		normalizedModifiers = append(normalizedModifiers, value)
		seen[value] = true
	}

	var normalizedVariant *string
	if variant != nil {
		value := strings.ToLower(strings.TrimSpace(*variant))
		if value != "" {
			normalizedVariant = &value
		}
	}

	return &GetsDocumentTypeV2{
		Base:      normalizedBase,
		Modifiers: normalizedModifiers,
		Variant:   normalizedVariant,
	}
}

func MapLogicalDocTypeToGetsV2(logicalType LogicalDocType) *GetsDocumentTypeV2 {
	name := string(logicalType)

	base := string(GetsDocumentBaseTaxInvoice)
	if strings.Contains(name, "CREDIT_NOTE") {
		base = string(GetsDocumentBaseCreditNote)
	} else if strings.Contains(name, "DEBIT_NOTE") {
		base = string(GetsDocumentBaseDebitNote)
	} else if strings.HasPrefix(name, "SIMPLIFIED") {
		base = string(GetsDocumentBaseSimplifiedInvoice)
	}

	modifiers := []string{}
	if strings.HasPrefix(name, "SIMPLIFIED") {
		modifiers = append(modifiers, string(GetsDocumentModifierB2C))
	}
	if strings.Contains(name, "EXPORT") {
		modifiers = append(modifiers, string(GetsDocumentModifierExport))
	}
	if strings.Contains(name, "SELF_BILLED") {
		modifiers = append(modifiers, string(GetsDocumentModifierSelfBilled))
	}
	if strings.Contains(name, "THIRD_PARTY") {
		modifiers = append(modifiers, string(GetsDocumentModifierThirdParty))
	}
	if strings.Contains(name, "NOMINAL_SUPPLY") {
		modifiers = append(modifiers, string(GetsDocumentModifierNominalSupply))
	}
	if strings.Contains(name, "SUMMARY") {
		modifiers = append(modifiers, string(GetsDocumentModifierSummary))
	}
	if strings.Contains(name, "B2G") {
		modifiers = append(modifiers, string(GetsDocumentModifierB2G))
	}

	return NewGetsDocumentTypeV2(base, modifiers, nil)
}
