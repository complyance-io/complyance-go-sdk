package models

// Environment represents the API environment
type Environment string

const (
	// EnvironmentSandbox represents the sandbox environment
	EnvironmentSandbox Environment = "sandbox"

	// EnvironmentProduction represents the production environment
	EnvironmentProduction Environment = "production"

	// EnvironmentLocal represents a local development environment
	EnvironmentLocal Environment = "local"
)

// String returns the string representation of the environment
func (e Environment) String() string {
	return string(e)
}