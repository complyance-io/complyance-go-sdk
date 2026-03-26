package erp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	BasePath    = "/documents/erp-exports"
	DefaultURL  = "https://api.complyance.io"
	DefaultTimeout = 30 * time.Second
)

type Client struct {
	APIKey     string
	BaseURL    string
	Environment string
	Timeout    time.Duration
	HTTPClient *http.Client
}

type Option func(*Client)

func WithEnvironment(env string) Option {
	return func(c *Client) {
		if env != "" {
			c.Environment = env
		}
	}
}

func WithBaseURL(url string) Option {
	return func(c *Client) {
		if url != "" {
			c.BaseURL = url
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if timeout > 0 {
			c.Timeout = timeout
		}
	}
}

func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		APIKey:      apiKey,
		BaseURL:     DefaultURL,
		Environment: "production",
		Timeout:     DefaultTimeout,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func NewClientFromEnv() (*Client, error) {
	apiKey := os.Getenv("COMPLYANCE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("COMPLYANCE_API_KEY environment variable is required")
	}

	env := os.Getenv("COMPLYANCE_ENVIRONMENT")
	if env == "" {
		env = "production"
	}

	baseURL := os.Getenv("COMPLYANCE_BASE_URL")

	return NewClient(apiKey, WithEnvironment(env), WithBaseURL(baseURL)), nil
}

func (c *Client) request(method, path string, body interface{}) (map[string]interface{}, error) {
	url := c.BaseURL + BasePath + path

	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *Client) ListJobs(limit *int) ([]map[string]interface{}, error) {
	path := fmt.Sprintf("/jobs?environment=%s", c.Environment)
	if limit != nil {
		path += fmt.Sprintf("&limit=%d", *limit)
	}

	resp, err := c.request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	jobs, ok := resp["jobs"].([]map[string]interface{})
	if !ok {
		return []map[string]interface{}{}, nil
	}

	return jobs, nil
}

func (c *Client) GetJob(jobID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/jobs/%s?environment=%s", jobID, c.Environment)
	resp, err := c.request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	job, ok := resp["job"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}

	return job, nil
}

func (c *Client) GetJobPayload(jobID string) (map[string]interface{}, error) {
	path := fmt.Sprintf("/jobs/%s/payload?environment=%s", jobID, c.Environment)
	resp, err := c.request("GET", path, nil)
	if err != nil {
		return nil, err
	}

	payload, ok := resp["payload"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, nil
	}

	return payload, nil
}

func (c *Client) AcknowledgeJob(jobID, status string, errorMsg *string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"status":      status,
		"environment": c.Environment,
	}
	if errorMsg != nil {
		body["error"] = *errorMsg
	}

	return c.request("POST", fmt.Sprintf("/jobs/%s/ack", jobID), body)
}

func (c *Client) TriggerManual(documentID string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"documentId":  documentID,
		"environment": c.Environment,
	}

	return c.request("POST", "/jobs/trigger-manual", body)
}

func (c *Client) GetConfig() (map[string]interface{}, error) {
	path := fmt.Sprintf("/config?environment=%s", c.Environment)
	return c.request("GET", path, nil)
}

func (c *Client) TestConnection(configID string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"environment": c.Environment,
	}

	return c.request("POST", fmt.Sprintf("/config/%s/test", configID), body)
}
