package http

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/complyance-io/complyance-go-sdk/v3/pkg/config"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/errors"
	"github.com/complyance-io/complyance-go-sdk/v3/pkg/models"
)

// Service provides high-level API operations
type Service struct {
	client Client
	config *config.Config
}

// NewService creates a new API service
func NewService(cfg *config.Config) *Service {
	return &Service{
		client: NewClient(cfg),
		config: cfg,
	}
}

// PushToUnify sends a request to the Complyance Unified API
func (s *Service) PushToUnify(ctx context.Context, request *models.UnifyRequest) (*models.UnifyResponse, error) {
	// Validate request
	if request == nil {
		return nil, errors.NewValidationError("request cannot be nil", nil)
	}

	if err := request.Validate(); err != nil {
		return nil, errors.NewValidationError("invalid request", err)
	}

	// Add API key to metadata if not present
	if request.Metadata == nil {
		request.Metadata = models.NewRequestMetadata()
	}
	if request.Metadata.APIKey == "" {
		request.Metadata.APIKey = s.config.APIKey
	}

	// Set environment in metadata
	request.Metadata.Environment = string(s.config.Environment)

	// Send request
	resp, err := s.client.Post(ctx, "/unify", request, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var unifyResponse models.UnifyResponse
	if err := resp.JSON(&unifyResponse); err != nil {
		return nil, errors.NewAPIError("failed to parse response", err).
			AddContext("body", resp.String())
	}

	return &unifyResponse, nil
}

// GetStatus checks the status of a submission
func (s *Service) GetStatus(ctx context.Context, submissionID string) (*models.UnifyResponse, error) {
	if submissionID == "" {
		return nil, errors.NewValidationError("submission ID is required", nil)
	}

	// Send request
	path := fmt.Sprintf("/status/%s", submissionID)
	resp, err := s.client.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var unifyResponse models.UnifyResponse
	if err := resp.JSON(&unifyResponse); err != nil {
		return nil, errors.NewAPIError("failed to parse response", err).
			AddContext("body", resp.String())
	}

	return &unifyResponse, nil
}

// ValidateMapping validates a field mapping
func (s *Service) ValidateMapping(ctx context.Context, source *models.Source, country string, payload map[string]interface{}) (*models.UnifyResponse, error) {
	if source == nil {
		return nil, errors.NewValidationError("source cannot be nil", nil)
	}

	if err := source.Validate(); err != nil {
		return nil, errors.NewValidationError("invalid source", err)
	}

	if country == "" {
		return nil, errors.NewValidationError("country is required", nil)
	}

	if len(country) != 2 {
		return nil, errors.NewValidationError("country must be a 2-letter ISO code", nil)
	}

	if payload == nil {
		return nil, errors.NewValidationError("payload cannot be nil", nil)
	}

	// Create validation request
	request := map[string]interface{}{
		"source":  source,
		"country": country,
		"payload": payload,
	}

	// Send request
	resp, err := s.client.Post(ctx, "/validate/mapping", request, nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var unifyResponse models.UnifyResponse
	if err := json.Unmarshal(resp.Body, &unifyResponse); err != nil {
		return nil, errors.NewAPIError("failed to parse response", err).
			AddContext("body", string(resp.Body))
	}

	return &unifyResponse, nil
}

// WithClient sets the HTTP client
func (s *Service) WithClient(client Client) *Service {
	s.client = client
	return s
}