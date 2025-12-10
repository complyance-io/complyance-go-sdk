package models

// SourceType represents the type of source system
type SourceType string

const (
	// SourceTypeFirstParty represents a first-party source system
	SourceTypeFirstParty SourceType = "FIRST_PARTY"

	// SourceTypeThirdParty represents a third-party source system
	SourceTypeThirdParty SourceType = "THIRD_PARTY"

	// SourceTypeMarketplace represents a marketplace source system
	SourceTypeMarketplace SourceType = "MARKETPLACE"
)

// String returns the string representation of the source type
func (st SourceType) String() string {
	return string(st)
}