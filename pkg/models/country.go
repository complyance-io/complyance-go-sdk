package models

import (
	"errors"
	"fmt"
	"strings"
)

// Country represents a country for document processing
type Country struct {
	// Code is the ISO 3166-1 alpha-2 country code
	Code string `json:"code"`

	// Name is the country name
	Name string `json:"name"`

	// Extensions contains country-specific extensions
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// CountryCode represents ISO 3166-1 alpha-2 country codes
type CountryCode string

const (
	// CountryCodeSA represents Saudi Arabia
	CountryCodeSA CountryCode = "SA"

	// CountryCodeAE represents United Arab Emirates
	CountryCodeAE CountryCode = "AE"

	// CountryCodeMY represents Malaysia
	CountryCodeMY CountryCode = "MY"

	// CountryCodeSG represents Singapore
	CountryCodeSG CountryCode = "SG"

	// CountryCodeAU represents Australia
	CountryCodeAU CountryCode = "AU"

	// CountryCodeNZ represents New Zealand
	CountryCodeNZ CountryCode = "NZ"

	// CountryCodeFR represents France
	CountryCodeFR CountryCode = "FR"

	// CountryCodeDE represents Germany
	CountryCodeDE CountryCode = "DE"

	// CountryCodePL represents Poland
	CountryCodePL CountryCode = "PL"

	// CountryCodeBE represents Belgium
	CountryCodeBE CountryCode = "BE"

	// CountryCodeBR represents Brazil
	CountryCodeBR CountryCode = "BR"
)

// String returns the string representation of the country code
func (cc CountryCode) String() string {
	return string(cc)
}

// NewCountry creates a new Country with the provided values
func NewCountry(code CountryCode, name string) *Country {
	return &Country{
		Code: string(code),
		Name: name,
	}
}

// Validate checks if the country is valid
func (c *Country) Validate() error {
	if c.Code == "" {
		return errors.New("country code is required")
	}

	if len(c.Code) != 2 {
		return errors.New("country code must be a 2-letter ISO code")
	}

	if c.Name == "" {
		return errors.New("country name is required")
	}

	return nil
}

// WithExtensions adds extensions to the country
func (c *Country) WithExtensions(extensions map[string]interface{}) *Country {
	c.Extensions = extensions
	return c
}

// AddExtension adds a single extension
func (c *Country) AddExtension(key string, value interface{}) *Country {
	if c.Extensions == nil {
		c.Extensions = make(map[string]interface{})
	}
	c.Extensions[key] = value
	return c
}

// String returns a string representation of the country
func (c *Country) String() string {
	var sb strings.Builder
	sb.WriteString("Country{")
	sb.WriteString("Code=")
	sb.WriteString(c.Code)
	sb.WriteString(", Name=")
	sb.WriteString(c.Name)
	if len(c.Extensions) > 0 {
		sb.WriteString(", Extensions=")
		sb.WriteString(fmt.Sprintf("%v", c.Extensions))
	}
	sb.WriteString("}")
	return sb.String()
}

// IsValidCountryCode checks if the provided string is a valid ISO 3166-1 alpha-2 country code
func IsValidCountryCode(code string) bool {
	if len(code) != 2 {
		return false
	}

	// This is a simplified validation that just checks the length
	// A more comprehensive validation would check against a list of valid country codes
	return true
}