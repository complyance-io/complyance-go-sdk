package models

// Operation represents the processing operation
type Operation string

const (
	// OperationSingle represents a single document operation
	OperationSingle Operation = "SINGLE"

	// OperationBatch represents a batch document operation
	OperationBatch Operation = "BATCH"
)

// String returns the string representation of the operation
func (o Operation) String() string {
	return string(o)
}