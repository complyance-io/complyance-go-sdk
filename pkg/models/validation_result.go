package models

// ValidationSeverity represents the severity level of a validation result
type ValidationSeverity string

const (
	// ValidationSeverityError represents a validation error
	ValidationSeverityError ValidationSeverity = "ERROR"

	// ValidationSeverityWarning represents a validation warning
	ValidationSeverityWarning ValidationSeverity = "WARNING"

	// ValidationSeverityInfo represents a validation information
	ValidationSeverityInfo ValidationSeverity = "INFO"
)

// String returns the string representation of the validation severity
func (vs ValidationSeverity) String() string {
	return string(vs)
}

// ValidationResult represents the result of a validation check
type ValidationResult struct {
	// Field is the name of the validated field
	Field string `json:"field"`

	// Message is the validation message
	Message string `json:"message"`

	// Severity is the validation severity level
	Severity ValidationSeverity `json:"severity"`

	// Code is an optional validation code
	Code string `json:"code,omitempty"`

	// Path is the JSON path to the field
	Path string `json:"path,omitempty"`

	// Value is the value that failed validation
	Value interface{} `json:"value,omitempty"`

	// Expected is the expected value or pattern
	Expected interface{} `json:"expected,omitempty"`
}

// NewValidationResult creates a new ValidationResult with the provided values
func NewValidationResult(field string, message string, severity ValidationSeverity) *ValidationResult {
	return &ValidationResult{
		Field:    field,
		Message:  message,
		Severity: severity,
	}
}

// WithCode adds a code to the validation result
func (vr *ValidationResult) WithCode(code string) *ValidationResult {
	vr.Code = code
	return vr
}

// WithPath adds a path to the validation result
func (vr *ValidationResult) WithPath(path string) *ValidationResult {
	vr.Path = path
	return vr
}

// WithValue adds the actual value to the validation result
func (vr *ValidationResult) WithValue(value interface{}) *ValidationResult {
	vr.Value = value
	return vr
}

// WithExpected adds the expected value to the validation result
func (vr *ValidationResult) WithExpected(expected interface{}) *ValidationResult {
	vr.Expected = expected
	return vr
}

// IsError returns true if the validation result is an error
func (vr *ValidationResult) IsError() bool {
	return vr.Severity == ValidationSeverityError
}

// IsWarning returns true if the validation result is a warning
func (vr *ValidationResult) IsWarning() bool {
	return vr.Severity == ValidationSeverityWarning
}

// IsInfo returns true if the validation result is an info
func (vr *ValidationResult) IsInfo() bool {
	return vr.Severity == ValidationSeverityInfo
}

// ValidationResults represents a collection of validation results
type ValidationResults struct {
	// Results is the list of validation results
	Results []*ValidationResult `json:"results"`
}

// NewValidationResults creates a new empty ValidationResults
func NewValidationResults() *ValidationResults {
	return &ValidationResults{
		Results: make([]*ValidationResult, 0),
	}
}

// AddResult adds a validation result to the collection
func (vrs *ValidationResults) AddResult(result *ValidationResult) *ValidationResults {
	vrs.Results = append(vrs.Results, result)
	return vrs
}

// AddError adds an error validation result
func (vrs *ValidationResults) AddError(field string, message string) *ValidationResults {
	vrs.Results = append(vrs.Results, NewValidationResult(field, message, ValidationSeverityError))
	return vrs
}

// AddWarning adds a warning validation result
func (vrs *ValidationResults) AddWarning(field string, message string) *ValidationResults {
	vrs.Results = append(vrs.Results, NewValidationResult(field, message, ValidationSeverityWarning))
	return vrs
}

// AddInfo adds an info validation result
func (vrs *ValidationResults) AddInfo(field string, message string) *ValidationResults {
	vrs.Results = append(vrs.Results, NewValidationResult(field, message, ValidationSeverityInfo))
	return vrs
}

// HasErrors returns true if there are any error results
func (vrs *ValidationResults) HasErrors() bool {
	for _, result := range vrs.Results {
		if result.IsError() {
			return true
		}
	}
	return false
}

// HasWarnings returns true if there are any warning results
func (vrs *ValidationResults) HasWarnings() bool {
	for _, result := range vrs.Results {
		if result.IsWarning() {
			return true
		}
	}
	return false
}

// Count returns the total number of validation results
func (vrs *ValidationResults) Count() int {
	return len(vrs.Results)
}

// ErrorCount returns the number of error results
func (vrs *ValidationResults) ErrorCount() int {
	count := 0
	for _, result := range vrs.Results {
		if result.IsError() {
			count++
		}
	}
	return count
}

// WarningCount returns the number of warning results
func (vrs *ValidationResults) WarningCount() int {
	count := 0
	for _, result := range vrs.Results {
		if result.IsWarning() {
			count++
		}
	}
	return count
}

// InfoCount returns the number of info results
func (vrs *ValidationResults) InfoCount() int {
	count := 0
	for _, result := range vrs.Results {
		if result.IsInfo() {
			count++
		}
	}
	return count
}