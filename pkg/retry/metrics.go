package retry

import (
	"expvar"
	"sync/atomic"
	"time"
)

// Metrics tracks retry statistics using expvar for observability
type Metrics struct {
	// Total number of attempts (including initial attempts)
	attempts int64
	
	// Number of successful operations
	successes int64
	
	// Number of failed operations
	failures int64
	
	// Number of times circuit breaker was open
	circuitOpens int64
	
	// Timestamp of the last retry
	lastRetryTime int64
	
	// Exported metrics via expvar
	expvarMap *expvar.Map
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	m := &Metrics{
		attempts:     0,
		successes:    0,
		failures:     0,
		circuitOpens: 0,
		lastRetryTime: 0,
	}
	
	// Create expvar map
	m.expvarMap = expvar.NewMap("complyance_sdk_retry_metrics")
	
	// Initialize expvar values
	m.updateExpvar()
	
	return m
}

// RecordAttempt increments the attempt counter
func (m *Metrics) RecordAttempt() {
	atomic.AddInt64(&m.attempts, 1)
	m.updateExpvar()
}

// RecordSuccess increments the success counter
func (m *Metrics) RecordSuccess() {
	atomic.AddInt64(&m.successes, 1)
	m.updateExpvar()
}

// RecordFailure increments the failure counter
func (m *Metrics) RecordFailure() {
	atomic.AddInt64(&m.failures, 1)
	atomic.StoreInt64(&m.lastRetryTime, time.Now().UnixNano())
	m.updateExpvar()
}

// RecordCircuitOpen increments the circuit open counter
func (m *Metrics) RecordCircuitOpen() {
	atomic.AddInt64(&m.circuitOpens, 1)
	m.updateExpvar()
}

// GetAttempts returns the total number of attempts
func (m *Metrics) GetAttempts() int64 {
	return atomic.LoadInt64(&m.attempts)
}

// GetSuccesses returns the number of successful operations
func (m *Metrics) GetSuccesses() int64 {
	return atomic.LoadInt64(&m.successes)
}

// GetFailures returns the number of failed operations
func (m *Metrics) GetFailures() int64 {
	return atomic.LoadInt64(&m.failures)
}

// GetCircuitOpens returns the number of times circuit breaker was open
func (m *Metrics) GetCircuitOpens() int64 {
	return atomic.LoadInt64(&m.circuitOpens)
}

// GetLastRetryTime returns the timestamp of the last retry
func (m *Metrics) GetLastRetryTime() time.Time {
	nanos := atomic.LoadInt64(&m.lastRetryTime)
	if nanos == 0 {
		return time.Time{}
	}
	return time.Unix(0, nanos)
}

// Reset resets all metrics to zero
func (m *Metrics) Reset() {
	atomic.StoreInt64(&m.attempts, 0)
	atomic.StoreInt64(&m.successes, 0)
	atomic.StoreInt64(&m.failures, 0)
	atomic.StoreInt64(&m.circuitOpens, 0)
	atomic.StoreInt64(&m.lastRetryTime, 0)
	m.updateExpvar()
}

// updateExpvar updates the expvar metrics
func (m *Metrics) updateExpvar() {
	m.expvarMap.Set("attempts", expvar.Func(func() interface{} {
		return atomic.LoadInt64(&m.attempts)
	}))
	
	m.expvarMap.Set("successes", expvar.Func(func() interface{} {
		return atomic.LoadInt64(&m.successes)
	}))
	
	m.expvarMap.Set("failures", expvar.Func(func() interface{} {
		return atomic.LoadInt64(&m.failures)
	}))
	
	m.expvarMap.Set("circuit_opens", expvar.Func(func() interface{} {
		return atomic.LoadInt64(&m.circuitOpens)
	}))
	
	m.expvarMap.Set("last_retry", expvar.Func(func() interface{} {
		return atomic.LoadInt64(&m.lastRetryTime)
	}))
	
	// Calculate success rate
	m.expvarMap.Set("success_rate", expvar.Func(func() interface{} {
		attempts := atomic.LoadInt64(&m.attempts)
		if attempts == 0 {
			return float64(0)
		}
		successes := atomic.LoadInt64(&m.successes)
		return float64(successes) / float64(attempts)
	}))
}