package tests

import (
	"encoding/json"
	"testing"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestFieldMapping(t *testing.T) {
	// Test basic field mapping
	mapping := models.NewFieldMapping("$.invoice.number", "$.invoice_number")
	assert.Equal(t, "$.invoice.number", mapping.SourcePath)
	assert.Equal(t, "$.invoice_number", mapping.TargetPath)
	assert.False(t, mapping.Required)

	// Test builder pattern
	mapping.WithTransformation("toUpperCase")
	mapping.WithDefaultValue("DEFAULT")
	mapping.WithRequired(true)
	mapping.WithDescription("Invoice number mapping")

	assert.Equal(t, "toUpperCase", mapping.Transformation)
	assert.Equal(t, "DEFAULT", mapping.DefaultValue)
	assert.True(t, mapping.Required)
	assert.Equal(t, "Invoice number mapping", mapping.Description)

	// Test validation
	err := mapping.Validate()
	assert.NoError(t, err)

	// Test invalid mapping
	invalidMapping := models.NewFieldMapping("", "")
	err = invalidMapping.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "source path is required")

	invalidMapping = models.NewFieldMapping("$.invoice.number", "")
	err = invalidMapping.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "target path is required")

	// Test string representation
	str := mapping.String()
	assert.Contains(t, str, "SourcePath=$.invoice.number")
	assert.Contains(t, str, "TargetPath=$.invoice_number")
	assert.Contains(t, str, "Transformation=toUpperCase")
	assert.Contains(t, str, "DefaultValue=DEFAULT")
	assert.Contains(t, str, "Required=true")
}

func TestFieldMappingSet(t *testing.T) {
	// Test basic field mapping set
	mappingSet := models.NewFieldMappingSet("invoice-mapping", "SA", models.DocumentTypeTaxInvoice)
	assert.Equal(t, "invoice-mapping", mappingSet.Name)
	assert.Equal(t, "SA", mappingSet.Country)
	assert.Equal(t, models.DocumentTypeTaxInvoice, mappingSet.DocumentType)
	assert.Empty(t, mappingSet.Mappings)

	// Test builder pattern
	mappingSet.WithDescription("Invoice mapping for Saudi Arabia")
	assert.Equal(t, "Invoice mapping for Saudi Arabia", mappingSet.Description)

	// Add mappings
	mapping1 := models.NewFieldMapping("$.invoice.number", "$.invoice_number").WithRequired(true)
	mapping2 := models.NewFieldMapping("$.invoice.date", "$.issue_date").WithRequired(true)
	mapping3 := models.NewFieldMapping("$.invoice.customer.name", "$.customer_name").WithRequired(false)

	mappingSet.AddMapping(mapping1)
	mappingSet.AddMapping(mapping2)
	mappingSet.AddMapping(mapping3)

	assert.Equal(t, 3, len(mappingSet.Mappings))

	// Test AddSimpleMapping
	mappingSet.AddSimpleMapping("$.invoice.total", "$.total_amount", true)
	assert.Equal(t, 4, len(mappingSet.Mappings))
	assert.Equal(t, "$.invoice.total", mappingSet.Mappings[3].SourcePath)
	assert.Equal(t, "$.total_amount", mappingSet.Mappings[3].TargetPath)
	assert.True(t, mappingSet.Mappings[3].Required)

	// Test validation
	err := mappingSet.Validate()
	assert.NoError(t, err)

	// Test invalid mapping set
	invalidMappingSet := models.NewFieldMappingSet("", "", "")
	err = invalidMappingSet.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	invalidMappingSet = models.NewFieldMappingSet("invoice-mapping", "", "")
	err = invalidMappingSet.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country is required")

	invalidMappingSet = models.NewFieldMappingSet("invoice-mapping", "SAU", "")
	err = invalidMappingSet.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "country must be a 2-letter ISO code")

	invalidMappingSet = models.NewFieldMappingSet("invoice-mapping", "SA", "")
	err = invalidMappingSet.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "document type is required")

	invalidMappingSet = models.NewFieldMappingSet("invoice-mapping", "SA", "INVALID_TYPE")
	err = invalidMappingSet.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid document type")

	invalidMappingSet = models.NewFieldMappingSet("invoice-mapping", "SA", models.DocumentTypeTaxInvoice)
	err = invalidMappingSet.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one mapping is required")

	// Test with invalid mapping
	invalidMappingSet = models.NewFieldMappingSet("invoice-mapping", "SA", models.DocumentTypeTaxInvoice)
	invalidMappingSet.AddMapping(models.NewFieldMapping("", ""))
	err = invalidMappingSet.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid mapping at index 0")

	// Test GetRequiredMappings
	requiredMappings := mappingSet.GetRequiredMappings()
	assert.Equal(t, 3, len(requiredMappings))

	// Test GetOptionalMappings
	optionalMappings := mappingSet.GetOptionalMappings()
	assert.Equal(t, 1, len(optionalMappings))
	assert.Equal(t, "$.invoice.customer.name", optionalMappings[0].SourcePath)

	// Test GetMappingByTargetPath
	mapping := mappingSet.GetMappingByTargetPath("$.invoice_number")
	assert.NotNil(t, mapping)
	assert.Equal(t, "$.invoice.number", mapping.SourcePath)

	mapping = mappingSet.GetMappingByTargetPath("non-existent")
	assert.Nil(t, mapping)

	// Test GetMappingBySourcePath
	mapping = mappingSet.GetMappingBySourcePath("$.invoice.number")
	assert.NotNil(t, mapping)
	assert.Equal(t, "$.invoice_number", mapping.TargetPath)

	mapping = mappingSet.GetMappingBySourcePath("non-existent")
	assert.Nil(t, mapping)
}

func TestFieldMappingJSON(t *testing.T) {
	// Create a field mapping
	mapping := models.NewFieldMapping("$.invoice.number", "$.invoice_number")
	mapping.WithTransformation("toUpperCase")
	mapping.WithDefaultValue("DEFAULT")
	mapping.WithRequired(true)
	mapping.WithDescription("Invoice number mapping")

	// Serialize to JSON
	jsonData, err := json.Marshal(mapping)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedMapping models.FieldMapping
	err = json.Unmarshal(jsonData, &deserializedMapping)
	assert.NoError(t, err)

	// Verify fields
	assert.Equal(t, "$.invoice.number", deserializedMapping.SourcePath)
	assert.Equal(t, "$.invoice_number", deserializedMapping.TargetPath)
	assert.Equal(t, "toUpperCase", deserializedMapping.Transformation)
	assert.Equal(t, "DEFAULT", deserializedMapping.DefaultValue)
	assert.True(t, deserializedMapping.Required)
	assert.Equal(t, "Invoice number mapping", deserializedMapping.Description)

	// Create a field mapping set
	mappingSet := models.NewFieldMappingSet("invoice-mapping", "SA", models.DocumentTypeTaxInvoice)
	mappingSet.WithDescription("Invoice mapping for Saudi Arabia")
	mappingSet.AddMapping(mapping)
	mappingSet.AddSimpleMapping("$.invoice.date", "$.issue_date", true)

	// Serialize to JSON
	jsonData, err = json.Marshal(mappingSet)
	assert.NoError(t, err)

	// Deserialize from JSON
	var deserializedMappingSet models.FieldMappingSet
	err = json.Unmarshal(jsonData, &deserializedMappingSet)
	assert.NoError(t, err)

	// Verify fields
	assert.Equal(t, "invoice-mapping", deserializedMappingSet.Name)
	assert.Equal(t, "SA", deserializedMappingSet.Country)
	assert.Equal(t, models.DocumentTypeTaxInvoice, deserializedMappingSet.DocumentType)
	assert.Equal(t, "Invoice mapping for Saudi Arabia", deserializedMappingSet.Description)
	assert.Equal(t, 2, len(deserializedMappingSet.Mappings))
	assert.Equal(t, "$.invoice.number", deserializedMappingSet.Mappings[0].SourcePath)
	assert.Equal(t, "$.invoice_number", deserializedMappingSet.Mappings[0].TargetPath)
	assert.Equal(t, "$.invoice.date", deserializedMappingSet.Mappings[1].SourcePath)
	assert.Equal(t, "$.issue_date", deserializedMappingSet.Mappings[1].TargetPath)
}