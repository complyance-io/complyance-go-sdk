/*
SDK Configuration for the Complyance SDK matching Python SDK exactly.
*/
package complyancesdk

// SDKConfig model matching Python SDK
type SDKConfig struct {
	APIKey                    string       `json:"api_key"`
	Environment               Environment  `json:"environment"`
	Sources                   []*Source    `json:"sources"`
	RetryConfig               *RetryConfig `json:"retry_config"`
	AutoGenerateTaxDestination bool         `json:"auto_generate_tax_destination"`
	CorrelationID             *string      `json:"correlation_id,omitempty"`
}

// NewSDKConfig creates a new SDK configuration
func NewSDKConfig(apiKey string, environment Environment, sources []*Source, retryConfig *RetryConfig) *SDKConfig {
	if retryConfig == nil {
		retryConfig = NewDefaultRetryConfig()
	}
	
	return &SDKConfig{
		APIKey:                    apiKey,
		Environment:               environment,
		Sources:                   sources,
		RetryConfig:               retryConfig,
		AutoGenerateTaxDestination: true,
		CorrelationID:             nil,
	}
}

// NewSDKConfigBuilder Create a builder for SDKConfig
func NewSDKConfigBuilder() *SDKConfigBuilder {
	return &SDKConfigBuilder{
		environment:               EnvironmentDev,
		sources:                   []*Source{},
		retryConfig:               nil,
		autoGenerateTaxDestination: true,
		correlationID:             nil,
	}
}

// GetAPIKey getter for API key
func (s *SDKConfig) GetAPIKey() string {
	return s.APIKey
}

// GetEnvironment getter for environment
func (s *SDKConfig) GetEnvironment() Environment {
	return s.Environment
}

// GetSources getter for sources
func (s *SDKConfig) GetSources() []*Source {
	return s.Sources
}

// GetRetryConfig getter for retry config
func (s *SDKConfig) GetRetryConfig() *RetryConfig {
	return s.RetryConfig
}

// IsAutoGenerateTaxDestination getter for auto generate tax destination
func (s *SDKConfig) IsAutoGenerateTaxDestination() bool {
	return s.AutoGenerateTaxDestination
}

// GetCorrelationID getter for correlation ID
func (s *SDKConfig) GetCorrelationID() *string {
	return s.CorrelationID
}

// SetRetryConfig setter for retry config
func (s *SDKConfig) SetRetryConfig(retryConfig *RetryConfig) {
	if retryConfig != nil {
		s.RetryConfig = retryConfig
	} else {
		s.RetryConfig = NewDefaultRetryConfig()
	}
}

// SetAPIKey setter for API key
func (s *SDKConfig) SetAPIKey(apiKey string) {
	s.APIKey = apiKey
}

// SetEnvironment setter for environment
func (s *SDKConfig) SetEnvironment(environment Environment) {
	s.Environment = environment
}

// SetSources setter for sources
func (s *SDKConfig) SetSources(sources []*Source) {
	s.Sources = sources
}

// SetAutoGenerateTaxDestination setter for auto generate tax destination
func (s *SDKConfig) SetAutoGenerateTaxDestination(autoGenerateTaxDestination bool) {
	s.AutoGenerateTaxDestination = autoGenerateTaxDestination
}

// SetCorrelationID setter for correlation ID
func (s *SDKConfig) SetCorrelationID(correlationID string) {
	s.CorrelationID = &correlationID
}

// SDKConfigBuilder Builder for SDKConfig matching Python SDK
type SDKConfigBuilder struct {
	apiKey                    *string
	environment               Environment
	sources                   []*Source
	retryConfig               *RetryConfig
	autoGenerateTaxDestination bool
	correlationID             *string
}

// APIKey setter for API key
func (b *SDKConfigBuilder) APIKey(apiKey string) *SDKConfigBuilder {
	b.apiKey = &apiKey
	return b
}

// Environment setter for environment
func (b *SDKConfigBuilder) Environment(environment Environment) *SDKConfigBuilder {
	b.environment = environment
	return b
}

// Sources setter for sources
func (b *SDKConfigBuilder) Sources(sources []*Source) *SDKConfigBuilder {
	if sources != nil {
		b.sources = sources
	} else {
		b.sources = []*Source{}
	}
	return b
}

// RetryConfig setter for retry config
func (b *SDKConfigBuilder) RetryConfig(retryConfig *RetryConfig) *SDKConfigBuilder {
	b.retryConfig = retryConfig
	return b
}

// AutoGenerateTaxDestination setter for auto generate tax destination
func (b *SDKConfigBuilder) AutoGenerateTaxDestination(autoGenerate bool) *SDKConfigBuilder {
	b.autoGenerateTaxDestination = autoGenerate
	return b
}

// CorrelationID setter for correlation ID
func (b *SDKConfigBuilder) CorrelationID(correlationID string) *SDKConfigBuilder {
	b.correlationID = &correlationID
	return b
}

// Build builds the SDKConfig
func (b *SDKConfigBuilder) Build() *SDKConfig {
	apiKey := ""
	if b.apiKey != nil {
		apiKey = *b.apiKey
	}
	
	config := NewSDKConfig(apiKey, b.environment, b.sources, b.retryConfig)
	config.AutoGenerateTaxDestination = b.autoGenerateTaxDestination
	config.CorrelationID = b.correlationID
	return config
}
