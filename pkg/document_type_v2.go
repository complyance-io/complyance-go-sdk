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
	GetsDocumentModifierB2B                 GetsDocumentModifier = "b2b"
	GetsDocumentModifierB2C                 GetsDocumentModifier = "b2c"
	GetsDocumentModifierB2G                 GetsDocumentModifier = "b2g"
	GetsDocumentModifierExport              GetsDocumentModifier = "export"
	GetsDocumentModifierSelfBilled          GetsDocumentModifier = "self_billed"
	GetsDocumentModifierThirdParty          GetsDocumentModifier = "third_party"
	GetsDocumentModifierNominal             GetsDocumentModifier = "nominal"
	GetsDocumentModifierNominalSupply       GetsDocumentModifier = "nominal_supply"
	GetsDocumentModifierSummary             GetsDocumentModifier = "summary"
	GetsDocumentModifierPrepayment          GetsDocumentModifier = "prepayment"
	GetsDocumentModifierAdjusted            GetsDocumentModifier = "adjusted"
	GetsDocumentModifierReceipt             GetsDocumentModifier = "receipt"
	GetsDocumentModifierZeroRated           GetsDocumentModifier = "zero_rated"
	GetsDocumentModifierReverseCharge       GetsDocumentModifier = "reverse_charge"
	GetsDocumentModifierContinuousSupply    GetsDocumentModifier = "continuous_supply"
	GetsDocumentModifierFreeTradeZone       GetsDocumentModifier = "free_trade_zone"
	GetsDocumentModifierIntraCommunitySupply GetsDocumentModifier = "intra_community_supply"
	GetsDocumentModifierConsolidated        GetsDocumentModifier = "consolidated"
)

type GetsDocumentVariant string

const (
	GetsDocumentVariantStandard             GetsDocumentVariant = "standard"
	GetsDocumentVariantPartial              GetsDocumentVariant = "partial"
	GetsDocumentVariantPartialConstruction  GetsDocumentVariant = "partial_construction"
	GetsDocumentVariantPartialFinalConstruction GetsDocumentVariant = "partial_final_construction"
	GetsDocumentVariantFinalConstruction    GetsDocumentVariant = "final_construction"
)

type GetsDocumentTypeV2 struct {
	Base      string   `json:"base"`
	Modifiers []string `json:"modifiers,omitempty"`
	Variant   *string  `json:"variant,omitempty"`
}

type GetsDocumentType = GetsDocumentTypeV2

var BASE = struct {
	TaxInvoice        GetsDocumentBase
	SimplifiedInvoice GetsDocumentBase
	CreditNote        GetsDocumentBase
	DebitNote         GetsDocumentBase
}{
	TaxInvoice:        GetsDocumentBaseTaxInvoice,
	SimplifiedInvoice: GetsDocumentBaseSimplifiedInvoice,
	CreditNote:        GetsDocumentBaseCreditNote,
	DebitNote:         GetsDocumentBaseDebitNote,
}

var MODIFIER = struct {
	B2B                  GetsDocumentModifier
	B2C                  GetsDocumentModifier
	B2G                  GetsDocumentModifier
	Export               GetsDocumentModifier
	SelfBilled           GetsDocumentModifier
	ThirdParty           GetsDocumentModifier
	Nominal              GetsDocumentModifier
	NominalSupply        GetsDocumentModifier
	Summary              GetsDocumentModifier
	Prepayment           GetsDocumentModifier
	Adjusted             GetsDocumentModifier
	Receipt              GetsDocumentModifier
	ZeroRated            GetsDocumentModifier
	ReverseCharge        GetsDocumentModifier
	ContinuousSupply     GetsDocumentModifier
	FreeTradeZone        GetsDocumentModifier
	IntraCommunitySupply GetsDocumentModifier
	Consolidated         GetsDocumentModifier
}{
	B2B:                  GetsDocumentModifierB2B,
	B2C:                  GetsDocumentModifierB2C,
	B2G:                  GetsDocumentModifierB2G,
	Export:               GetsDocumentModifierExport,
	SelfBilled:           GetsDocumentModifierSelfBilled,
	ThirdParty:           GetsDocumentModifierThirdParty,
	Nominal:              GetsDocumentModifierNominal,
	NominalSupply:        GetsDocumentModifierNominalSupply,
	Summary:              GetsDocumentModifierSummary,
	Prepayment:           GetsDocumentModifierPrepayment,
	Adjusted:             GetsDocumentModifierAdjusted,
	Receipt:              GetsDocumentModifierReceipt,
	ZeroRated:            GetsDocumentModifierZeroRated,
	ReverseCharge:        GetsDocumentModifierReverseCharge,
	ContinuousSupply:     GetsDocumentModifierContinuousSupply,
	FreeTradeZone:        GetsDocumentModifierFreeTradeZone,
	IntraCommunitySupply: GetsDocumentModifierIntraCommunitySupply,
	Consolidated:         GetsDocumentModifierConsolidated,
}

type GetsDocumentTypeBuilder struct {
	base      string
	modifiers []string
	variant   *string
}

func NewGetsDocumentTypeBuilder() *GetsDocumentTypeBuilder {
	return &GetsDocumentTypeBuilder{
		base:      string(GetsDocumentBaseTaxInvoice),
		modifiers: []string{},
	}
}

func (b *GetsDocumentTypeBuilder) Base(base string) *GetsDocumentTypeBuilder {
	b.base = strings.ToLower(strings.TrimSpace(base))
	return b
}

func (b *GetsDocumentTypeBuilder) Modifier(modifier string) *GetsDocumentTypeBuilder {
	return b.Modifiers([]string{modifier})
}

func (b *GetsDocumentTypeBuilder) AddModifier(modifier string) *GetsDocumentTypeBuilder {
	return b.Modifier(modifier)
}

func (b *GetsDocumentTypeBuilder) Modifiers(modifiers []string) *GetsDocumentTypeBuilder {
	current := append([]string{}, b.modifiers...)
	current = append(current, modifiers...)
	normalized := make([]string, 0, len(current))
	seen := map[string]bool{}
	for _, modifier := range current {
		value := strings.ToLower(strings.TrimSpace(modifier))
		if value == "" || seen[value] {
			continue
		}
		normalized = append(normalized, value)
		seen[value] = true
	}
	b.modifiers = normalized
	return b
}

func (b *GetsDocumentTypeBuilder) Variant(variant *string) *GetsDocumentTypeBuilder {
	if variant == nil {
		b.variant = nil
		return b
	}
	value := strings.ToLower(strings.TrimSpace(*variant))
	if value == "" {
		b.variant = nil
		return b
	}
	b.variant = &value
	return b
}

func (b *GetsDocumentTypeBuilder) Build() *GetsDocumentType {
	return NewGetsDocumentTypeV2(b.base, b.modifiers, b.variant)
}

type docTypeFactory struct{}

var DocType = docTypeFactory{}

func (docTypeFactory) Of(base GetsDocumentBase, modifiers ...GetsDocumentModifier) *GetsDocumentType {
	modifierStrings := make([]string, 0, len(modifiers))
	for _, modifier := range modifiers {
		modifierStrings = append(modifierStrings, string(modifier))
	}
	return NewGetsDocumentTypeV2(string(base), modifierStrings, nil)
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
