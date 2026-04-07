/*
Circuit Breaker implementation matching Python SDK exactly.
*/
package complyancesdk

import (
	"log"
	"time"
)

// CircuitState Circuit breaker states
type CircuitState string

const (
	CircuitStateClosed   CircuitState = "CLOSED"
	CircuitStateOpen     CircuitState = "OPEN"
	CircuitStateHalfOpen CircuitState = "HALF_OPEN"
)

// CircuitBreaker Circuit breaker implementation matching Python SDK
type CircuitBreaker struct {
	config          *CircuitBreakerConfig
	state           CircuitState
	failureCount    int
	lastFailureTime int64
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(config *CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		config:          config,
		state:           CircuitStateClosed,
		failureCount:    0,
		lastFailureTime: 0,
	}
}

// Execute operation with circuit breaker
func (c *CircuitBreaker) Execute(operation func() (interface{}, error)) (interface{}, error) {
	if c.state == CircuitStateOpen {
		currentTime := time.Now().UnixNano() / int64(time.Millisecond) // Convert to milliseconds
		timeSinceLastFailure := currentTime - c.lastFailureTime
		remainingTime := 60000 - timeSinceLastFailure // 1 minute timeout

		if c.shouldAttemptReset() {
			c.state = CircuitStateHalfOpen
		} else {
			return nil, NewSDKError(NewErrorDetailWithCode(
				ErrorCodeCircuitBreakerOpen,
				"Circuit breaker is open - "+string(remainingTime/1000)+" seconds remaining",
			))
		}
	}

	result, err := operation()
	if err != nil {
		c.onFailure()
		return nil, err
	} else {
		c.onSuccess()
		return result, nil
	}
}

// onSuccess Handle successful operation
func (c *CircuitBreaker) onSuccess() {
	if c.state == CircuitStateHalfOpen {
		c.state = CircuitStateClosed
		c.failureCount = 0
	}
}

// onFailure Handle failed operation
func (c *CircuitBreaker) onFailure() {
	c.failureCount++
	c.lastFailureTime = time.Now().UnixNano() / int64(time.Millisecond) // Convert to milliseconds

	if c.failureCount >= c.config.GetFailureThreshold() {
		c.state = CircuitStateOpen
	}
}

// shouldAttemptReset Check if circuit breaker should attempt reset
func (c *CircuitBreaker) shouldAttemptReset() bool {
	currentTime := time.Now().UnixNano() / int64(time.Millisecond) // Convert to milliseconds
	timeSinceLastFailure := currentTime - c.lastFailureTime
	timeoutMillis := int64(c.config.GetTimeout())

	if timeSinceLastFailure >= timeoutMillis {
		log.Printf("Circuit breaker timeout expired (%d ms) - attempting reset", timeSinceLastFailure)
		return true
	} else {
		remainingTime := timeoutMillis - timeSinceLastFailure
		log.Printf("Circuit breaker timeout not expired - %d seconds remaining", remainingTime/1000)
		return false
	}
}

// GetState Get circuit breaker state
func (c *CircuitBreaker) GetState() CircuitState {
	return c.state
}

// GetFailureCount Get failure count
func (c *CircuitBreaker) GetFailureCount() int {
	return c.failureCount
}

// GetLastFailureTime Get last failure time
func (c *CircuitBreaker) GetLastFailureTime() int64 {
	return c.lastFailureTime
}

// IsOpen Check if circuit breaker is open
func (c *CircuitBreaker) IsOpen() bool {
	return c.state == CircuitStateOpen
}

// IsClosed Check if circuit breaker is closed
func (c *CircuitBreaker) IsClosed() bool {
	return c.state == CircuitStateClosed
}

// IsHalfOpen Check if circuit breaker is half open
func (c *CircuitBreaker) IsHalfOpen() bool {
	return c.state == CircuitStateHalfOpen
}
