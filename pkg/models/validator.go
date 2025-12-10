package models

// Validator is an interface for objects that can validate themselves
type Validator interface {
	// Validate checks if the object is valid
	Validate() error
}

// ValidateAll validates a slice of validators
func ValidateAll(validators ...Validator) error {
	for _, v := range validators {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateField validates a field against a validation function
func ValidateField(field string, value interface{}, validationFunc func(interface{}) error) error {
	if err := validationFunc(value); err != nil {
		return err
	}
	return nil
}