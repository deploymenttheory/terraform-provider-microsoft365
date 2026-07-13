package crud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/assert"
)

// mockGraphStatusError simulates a Graph API error with a specific HTTP status code.
// It satisfies the GetStatusCode/GetErrorEscaped interface recognized by errors.GraphError.
type mockGraphStatusError struct {
	statusCode int
}

func (e *mockGraphStatusError) Error() string {
	return fmt.Sprintf("mock graph error with status %d", e.statusCode)
}

func (e *mockGraphStatusError) GetStatusCode() int {
	return e.statusCode
}

func (e *mockGraphStatusError) GetErrorEscaped() odataerrors.MainErrorable {
	return odataerrors.NewMainError()
}

func TestDefaultWriteWithRetryOptions(t *testing.T) {
	opts := DefaultWriteWithRetryOptions()

	assert.Equal(t, 10, opts.MaxRetries)
	assert.Equal(t, 5*time.Second, opts.RetryInterval)
	assert.Equal(t, constants.TfOperationCreate, opts.Operation)
	assert.Equal(t, "", opts.ResourceTypeName)
	assert.Equal(t, "", opts.ResourceID)
}

func TestWriteWithRetry_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultWriteWithRetryOptions()
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	writeFunc := func(ctx context.Context) error {
		callCount++
		return nil // Success on first attempt
	}

	err := WriteWithRetry(ctx, writeFunc, opts)

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount, "Write function should be called exactly once")
}

func TestWriteWithRetry_SuccessAfter404Retries(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultWriteWithRetryOptions()
	opts.MaxRetries = 3
	opts.RetryInterval = 100 * time.Millisecond // Shorter for testing
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	writeFunc := func(ctx context.Context) error {
		callCount++
		if callCount < 3 {
			// 404: referenced object not yet propagated across Entra replicas
			return &mockGraphStatusError{statusCode: 404}
		}
		return nil // Success on third attempt
	}

	start := time.Now()
	err := WriteWithRetry(ctx, writeFunc, opts)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, 3, callCount, "Write function should be called 3 times")
	assert.GreaterOrEqual(t, duration, 200*time.Millisecond, "Should wait between retries")
}

func TestWriteWithRetry_NonRetryableError(t *testing.T) {
	nonRetryableStatusCodes := []int{400, 401, 403, 409, 422}

	for _, statusCode := range nonRetryableStatusCodes {
		t.Run(fmt.Sprintf("status_%d", statusCode), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			opts := DefaultWriteWithRetryOptions()
			opts.RetryInterval = 50 * time.Millisecond
			opts.ResourceTypeName = "TestResource"
			opts.ResourceID = "test-123"

			callCount := 0
			writeFunc := func(ctx context.Context) error {
				callCount++
				return &mockGraphStatusError{statusCode: statusCode}
			}

			err := WriteWithRetry(ctx, writeFunc, opts)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "non-retryable error")
			assert.Equal(t, 1, callCount, "Write function should be called exactly once")
		})
	}
}

func TestWriteWithRetry_MaxRetriesExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultWriteWithRetryOptions()
	opts.MaxRetries = 2
	opts.RetryInterval = 50 * time.Millisecond // Shorter for testing
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	writeFunc := func(ctx context.Context) error {
		callCount++
		// Persistent 404 (e.g. the referenced object genuinely does not exist)
		return &mockGraphStatusError{statusCode: 404}
	}

	err := WriteWithRetry(ctx, writeFunc, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to write resource TestResource after 3 attempts")
	assert.Equal(t, 3, callCount, "Write function should be called 3 times (initial + 2 retries)")
}

func TestWriteWithRetry_RetryableServerErrors(t *testing.T) {
	retryableStatusCodes := []int{429, 500, 502, 503, 504}

	for _, statusCode := range retryableStatusCodes {
		t.Run(fmt.Sprintf("status_%d", statusCode), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			opts := DefaultWriteWithRetryOptions()
			opts.MaxRetries = 1
			opts.RetryInterval = 50 * time.Millisecond
			opts.ResourceTypeName = "TestResource"
			opts.ResourceID = "test-123"

			callCount := 0
			writeFunc := func(ctx context.Context) error {
				callCount++
				if callCount == 1 {
					return &mockGraphStatusError{statusCode: statusCode}
				}
				return nil
			}

			err := WriteWithRetry(ctx, writeFunc, opts)

			assert.NoError(t, err)
			assert.Equal(t, 2, callCount, "Write function should succeed on the retry")
		})
	}
}

func TestWriteWithRetry_ContextTimeout(t *testing.T) {
	// Context with enough time for the function to start but not complete retries
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := DefaultWriteWithRetryOptions()
	opts.RetryInterval = 1 * time.Second
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	writeFunc := func(ctx context.Context) error {
		callCount++
		return &mockGraphStatusError{statusCode: 404}
	}

	start := time.Now()
	err := WriteWithRetry(ctx, writeFunc, opts)
	duration := time.Since(start)

	assert.Error(t, err)
	// Should have at least one call
	assert.GreaterOrEqual(t, callCount, 1, "Should have at least one attempt")
	// Should complete in reasonable time due to context timeout
	assert.LessOrEqual(t, duration, 3*time.Second, "Should be limited by context timeout")
}

func TestWriteWithRetry_ContextCancellation(t *testing.T) {
	// Start with a context that has a deadline to satisfy the requirement
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	opts := DefaultWriteWithRetryOptions()
	opts.MaxRetries = 5
	opts.RetryInterval = 100 * time.Millisecond
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	writeFunc := func(ctx context.Context) error {
		callCount++
		if callCount == 2 {
			// Cancel context during second attempt
			cancel()
			// Give some time for cancellation to propagate
			time.Sleep(10 * time.Millisecond)
		}
		return &mockGraphStatusError{statusCode: 404}
	}

	err := WriteWithRetry(ctx, writeFunc, opts)

	assert.Error(t, err)
	// The error could be either context cancellation or the retry failure
	assert.True(t, callCount >= 2, "Should have at least 2 attempts before cancellation")
}

func TestWriteWithRetry_NoDeadline(t *testing.T) {
	ctx := context.Background() // No deadline

	opts := DefaultWriteWithRetryOptions()

	writeFunc := func(ctx context.Context) error {
		return nil
	}

	err := WriteWithRetry(ctx, writeFunc, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context must have a deadline")
}

func TestWriteWithRetry_InsufficientTime(t *testing.T) {
	// Context with very short remaining time
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Sleep to reduce remaining time
	time.Sleep(400 * time.Millisecond)

	opts := DefaultWriteWithRetryOptions()
	opts.RetryInterval = 1 * time.Second // Longer than remaining time

	writeFunc := func(ctx context.Context) error {
		return nil
	}

	err := WriteWithRetry(ctx, writeFunc, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient time remaining")
}

func TestWriteWithRetry_TimeConstrainedRetries(t *testing.T) {
	// Context with limited time that allows only 2 retries
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := DefaultWriteWithRetryOptions()
	opts.MaxRetries = 10 // More than time allows
	opts.RetryInterval = 2 * time.Second
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	writeFunc := func(ctx context.Context) error {
		callCount++
		return &mockGraphStatusError{statusCode: 404}
	}

	start := time.Now()
	err := WriteWithRetry(ctx, writeFunc, opts)
	duration := time.Since(start)

	assert.Error(t, err)
	// Should be constrained by time, not MaxRetries
	assert.LessOrEqual(t, callCount, 3, "Should be limited by available time")
	assert.LessOrEqual(t, duration, 6*time.Second, "Should complete within reasonable time")
}

func TestWriteWithRetry_DefaultValues(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Test with minimal options
	opts := WriteWithRetryOptions{}

	writeFunc := func(ctx context.Context) error {
		return nil
	}

	err := WriteWithRetry(ctx, writeFunc, opts)

	assert.NoError(t, err)
}

func TestWriteWithRetry_UnknownErrorType(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultWriteWithRetryOptions()
	opts.MaxRetries = 2 // Limit retries to prevent long-running test
	opts.RetryInterval = 50 * time.Millisecond
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	writeFunc := func(ctx context.Context) error {
		callCount++
		// Return standard Go error - this will be treated as 500 (retryable)
		return fmt.Errorf("unknown error type")
	}

	err := WriteWithRetry(ctx, writeFunc, opts)

	assert.Error(t, err)
	// Standard Go errors are treated as 500 (retryable), so should retry until max attempts
	assert.Equal(t, 3, callCount, "Should retry for unknown error type (treated as 500)")
}

func TestWriteWithRetry_ResourceLogging(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Test with empty resource names (should use defaults)
	opts := WriteWithRetryOptions{
		MaxRetries:    1,
		RetryInterval: 50 * time.Millisecond,
		Operation:     "", // Should default to constants.TfOperationCreate
		// ResourceTypeName and ResourceID left empty
	}

	callCount := 0
	writeFunc := func(ctx context.Context) error {
		callCount++
		return nil
	}

	err := WriteWithRetry(ctx, writeFunc, opts)

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}
