# Complyance Go SDK Testing Suite

This directory contains comprehensive tests for the Complyance Go SDK. The testing suite includes unit tests, integration tests, table-driven tests, benchmarks, and race condition detection tests.

## Running Tests

### Basic Test Execution

To run all tests:

```bash
go test ./...
```

To run tests in a specific file:

```bash
go test ./tests/sdk_test.go
```

To run a specific test:

```bash
go test -run TestSDKConfiguration
```

### Running Tests with Race Detection

Race detection helps identify potential race conditions in concurrent code:

```bash
go test -race ./...
```

For specific race condition tests:

```bash
go test -race ./tests/race_condition_test.go
```

### Running Benchmarks

To run all benchmarks:

```bash
go test -bench=. ./...
```

To run specific benchmarks:

```bash
go test -bench=BenchmarkSDKClient ./tests/benchmark_test.go
```

To get more detailed benchmark results:

```bash
go test -bench=. -benchmem ./...
```

### Running Short Tests

For quick testing during development:

```bash
go test -short ./...
```

## Test Coverage

To generate test coverage reports:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Structure

The testing suite is organized as follows:

1. **Unit Tests**: Basic tests for individual components
2. **Integration Tests**: Tests that verify components work together
3. **Table-Driven Tests**: Comprehensive tests with multiple test cases
4. **Benchmarks**: Performance tests for critical components
5. **Race Condition Tests**: Tests designed to detect race conditions in concurrent code

## Test Helpers

The `test_helpers.go` file provides common utilities for testing, including:

- `MockServer`: A test server that returns predefined responses
- `CreateTestSource`: Creates a test source for use in tests
- `CreateTestPayload`: Creates a test payload with realistic data
- `SkipIfShort`: Helper to skip tests when running in short mode

## Best Practices

When adding new tests, follow these best practices:

1. Use table-driven tests for comprehensive test coverage
2. Include both positive and negative test cases
3. Test edge cases and error conditions
4. Use benchmarks for performance-critical code
5. Include race condition tests for concurrent code
6. Keep tests independent and isolated
7. Use descriptive test names and clear assertions
