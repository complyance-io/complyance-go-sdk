package retry

import (
	"sync"
	"sync/atomic"
	"time"
)

// CircuitState represents the state of the circuit breaker
type CircuitState int32

const (
	// CircuitClosed indicates the circuit is closed and requests are allowed
	CircuitClosed CircuitState = iota
	
	// CircuitOpen indicates the circuit is open and requests are blocked
	CircuitOpen
	
	// CircuitHalfOpen indicates the circuit is testing if it can be closed
	CircuitHalfOpen
)

// CircuitBreaker implements the circuit breaker pattern with atomic operations
type CircuitBreaker struct {
	// state is the current state of the circuit breaker
	state int32
	
	// failureCount is the current count of consecutive failures
	failureCount int32
	
	// failureThreshold is the number of failures before opening the circuit
	failureThreshold int32
	
	// timeout is the duration to keep the circuit open
	timeout time.Duration
	
	// lastStateChange is the time of the last state change
	lastStateChange time.Time
	
	// mutex protects lastStateChange
	mutex sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(failureThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:           int32(CircuitClosed),
		failureCount:    0,
		failureThreshold: int32(failureThreshold),
		timeout:         timeout,
		lastStateChange: time.Now(),
	}
}

// IsOpen returns true if the circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	state := CircuitState(atomic.LoadInt32(&cb.state))
	
	// If circuit is open, check if timeout has elapsed
	if state == CircuitOpen {
		cb.mutex.RLock()
		elapsed := time.Since(cb.lastStateChange)
		cb.mutex.RUnlock()
		
		// If timeout has elapsed, transition to half-open
		if elapsed >= cb.timeout {
			cb.transitionToHalfOpen()
			return false
		}
		return true
	}
	
	return false
}

// RecordSuccess records a successful operation
func (cb *CircuitBreaker) RecordSuccess() {
	state := CircuitState(atomic.LoadInt32(&cb.state))
	
	// Reset failure count
	atomic.StoreInt32(&cb.failureCount, 0)
	
	// If circuit is half-open, close it
	if state == CircuitHalfOpen {
		cb.transitionToClosed()
	}
}

// RecordFailure records a failed operation
func (cb *CircuitBreaker) RecordFailure() {
	state := CircuitState(atomic.LoadInt32(&cb.state))
	
	// Increment failure count
	newCount := atomic.AddInt32(&cb.failureCount, 1)
	
	// If circuit is closed and failure threshold is reached, open it
	if state == CircuitClosed && newCount >= cb.failureThreshold {
		cb.transitionToOpen()
	}
	
	// If circuit is half-open, open it again
	if state == CircuitHalfOpen {
		cb.transitionToOpen()
	}
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	atomic.StoreInt32(&cb.failureCount, 0)
	cb.transitionToClosed()
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() CircuitState {
	return CircuitState(atomic.LoadInt32(&cb.state))
}

// GetFailureCount returns the current failure count
func (cb *CircuitBreaker) GetFailureCount() int {
	return int(atomic.LoadInt32(&cb.failureCount))
}

// GetFailureThreshold returns the failure threshold
func (cb *CircuitBreaker) GetFailureThreshold() int {
	return int(cb.failureThreshold)
}

// GetTimeout returns the timeout duration
func (cb *CircuitBreaker) GetTimeout() time.Duration {
	return cb.timeout
}

// GetLastStateChange returns the time of the last state change
func (cb *CircuitBreaker) GetLastStateChange() time.Time {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.lastStateChange
}

// transitionToOpen changes the circuit state to open
func (cb *CircuitBreaker) transitionToOpen() {
	// Only transition if not already open
	if atomic.CompareAndSwapInt32(&cb.state, int32(CircuitClosed), int32(CircuitOpen)) ||
	   atomic.CompareAndSwapInt32(&cb.state, int32(CircuitHalfOpen), int32(CircuitOpen)) {
		cb.mutex.Lock()
		cb.lastStateChange = time.Now()
		cb.mutex.Unlock()
	}
}

// transitionToHalfOpen changes the circuit state to half-open
func (cb *CircuitBreaker) transitionToHalfOpen() {
	// Only transition if currently open
	if atomic.CompareAndSwapInt32(&cb.state, int32(CircuitOpen), int32(CircuitHalfOpen)) {
		cb.mutex.Lock()
		cb.lastStateChange = time.Now()
		cb.mutex.Unlock()
	}
}

// transitionToClosed changes the circuit state to closed
func (cb *CircuitBreaker) transitionToClosed() {
	// Transition from any state to closed
	oldState := atomic.SwapInt32(&cb.state, int32(CircuitClosed))
	if oldState != int32(CircuitClosed) {
		cb.mutex.Lock()
		cb.lastStateChange = time.Now()
		cb.mutex.Unlock()
	}
}