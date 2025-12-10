package models

import "errors"

// Destination represents a destination for processed documents
type Destination struct {
	// Type is the destination type
	Type string `json:"type"`

	// Config contains destination-specific configuration
	Config map[string]interface{} `json:"config"`
}

// NewDestination creates a new Destination with the provided values
func NewDestination(destinationType string, config map[string]interface{}) *Destination {
	return &Destination{
		Type:   destinationType,
		Config: config,
	}
}

// Validate checks if the destination is valid
func (d *Destination) Validate() error {
	if d.Type == "" {
		return errors.New("destination type is required")
	}

	if d.Config == nil {
		return errors.New("destination config is required")
	}

	return nil
}

// WithConfig sets the destination config
func (d *Destination) WithConfig(config map[string]interface{}) *Destination {
	d.Config = config
	return d
}

// AddConfigField adds a single config field
func (d *Destination) AddConfigField(key string, value interface{}) *Destination {
	if d.Config == nil {
		d.Config = make(map[string]interface{})
	}
	d.Config[key] = value
	return d
}