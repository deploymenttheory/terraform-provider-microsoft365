package crud

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/stretchr/testify/assert"
)

// Helper functions to create test errors
func createRetryableError() error {
	return &url.Error{
		Op:  constants.TfOperationDelete,
		URL: "https://graph.microsoft.com/v1.0/resource",
		Err: fmt.Errorf("context deadline exceeded"),
	}
}

func createNonRetryableError() error {
	return &url.Error{
		Op:  constants.TfOperationDelete,
		URL: "https://graph.microsoft.com/v1.0/resource",
		Err: fmt.Errorf("resource not found"),
	}
}

func TestDefaultDeleteWithRetryOptions(t *testing.T) {
	opts := DefaultDeleteWithRetryOptions()

	assert.Equal(t, 10, opts.MaxRetries)
	assert.Equal(t, 30*time.Second, opts.RetryInterval)
	assert.Equal(t, constants.TfOperationDelete, opts.Operation)
	assert.Equal(t, "", opts.ResourceTypeName)
	assert.Equal(t, "", opts.ResourceID)
}

func TestDeleteWithRetry_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		return nil // Success on first attempt
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount, "Delete function should be called exactly once")
}

func TestDeleteWithRetry_SuccessAfterRetries(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()
	opts.MaxRetries = 3
	opts.RetryInterval = 100 * time.Millisecond // Shorter for testing
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		if callCount < 3 {
			// Return retryable timeout error for first 2 attempts
			return &url.Error{
				Op:  constants.TfOperationDelete,
				URL: "https://graph.microsoft.com/v1.0/resource",
				Err: fmt.Errorf("context deadline exceeded"),
			}
		}
		return nil // Success on third attempt
	}

	start := time.Now()
	err := DeleteWithRetry(ctx, deleteFunc, opts)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, 3, callCount, "Delete function should be called 3 times")
	assert.GreaterOrEqual(t, duration, 200*time.Millisecond, "Should wait between retries")
}

func TestDeleteWithRetry_NonRetryableError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		// Return a standard URL error which will be treated as non-retryable
		return &url.Error{
			Op:  constants.TfOperationDelete,
			URL: "https://graph.microsoft.com/v1.0/resource",
			Err: fmt.Errorf("resource not found"),
		}
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "non-retryable error")
	assert.Equal(t, 1, callCount, "Delete function should be called exactly once")
}

func TestDeleteWithRetry_MaxRetriesExceeded(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()
	opts.MaxRetries = 2
	opts.RetryInterval = 50 * time.Millisecond // Shorter for testing
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		// Always return retryable timeout error
		return &url.Error{
			Op:  constants.TfOperationDelete,
			URL: "https://graph.microsoft.com/v1.0/resource",
			Err: fmt.Errorf("context deadline exceeded"),
		}
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to delete resource TestResource after 3 attempts")
	assert.Equal(t, 3, callCount, "Delete function should be called 3 times (initial + 2 retries)")
}

func TestDeleteWithRetry_ContextTimeout(t *testing.T) {
	// Context with enough time for the function to start but not complete retries
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()
	opts.RetryInterval = 1 * time.Second
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		return createRetryableError()
	}

	start := time.Now()
	err := DeleteWithRetry(ctx, deleteFunc, opts)
	duration := time.Since(start)

	assert.Error(t, err)
	// Should have at least one call
	assert.GreaterOrEqual(t, callCount, 1, "Should have at least one attempt")
	// Should complete in reasonable time due to context timeout
	assert.LessOrEqual(t, duration, 3*time.Second, "Should be limited by context timeout")
}

func TestDeleteWithRetry_ContextCancellation(t *testing.T) {
	// Start with a context that has a deadline to satisfy the requirement
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	opts := DefaultDeleteWithRetryOptions()
	opts.MaxRetries = 5
	opts.RetryInterval = 100 * time.Millisecond
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		if callCount == 2 {
			// Cancel context during second attempt
			cancel()
			// Give some time for cancellation to propagate
			time.Sleep(10 * time.Millisecond)
		}
		return createRetryableError()
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.Error(t, err)
	// The error could be either context cancellation or the retry failure
	assert.True(t, callCount >= 2, "Should have at least 2 attempts before cancellation")
}

func TestDeleteWithRetry_NoDeadline(t *testing.T) {
	ctx := context.Background() // No deadline

	opts := DefaultDeleteWithRetryOptions()

	deleteFunc := func(ctx context.Context) error {
		return nil
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context must have a deadline")
}

func TestDeleteWithRetry_InsufficientTime(t *testing.T) {
	// Context with very short remaining time
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Sleep to reduce remaining time
	time.Sleep(400 * time.Millisecond)

	opts := DefaultDeleteWithRetryOptions()
	opts.RetryInterval = 1 * time.Second // Longer than remaining time

	deleteFunc := func(ctx context.Context) error {
		return nil
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient time remaining")
}

func TestDeleteWithRetry_TimeConstrainedRetries(t *testing.T) {
	// Context with limited time that allows only 2 retries
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()
	opts.MaxRetries = 10 // More than time allows
	opts.RetryInterval = 2 * time.Second
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		return &url.Error{
			Op:  constants.TfOperationDelete,
			URL: "https://graph.microsoft.com/v1.0/resource",
			Err: fmt.Errorf("context deadline exceeded"),
		}
	}

	start := time.Now()
	err := DeleteWithRetry(ctx, deleteFunc, opts)
	duration := time.Since(start)

	assert.Error(t, err)
	// Should be constrained by time, not MaxRetries
	assert.LessOrEqual(t, callCount, 3, "Should be limited by available time")
	assert.LessOrEqual(t, duration, 6*time.Second, "Should complete within reasonable time")
}

func TestDeleteWithRetry_DefaultValues(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Test with minimal options
	opts := DeleteWithRetryOptions{}

	deleteFunc := func(ctx context.Context) error {
		return nil
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.NoError(t, err)
}

func TestDeleteWithRetry_ErrorTypes(t *testing.T) {
	t.Run("Retryable timeout error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		opts := DefaultDeleteWithRetryOptions()
		opts.MaxRetries = 2
		opts.RetryInterval = 50 * time.Millisecond
		opts.ResourceTypeName = "TestResource"
		opts.ResourceID = "test-123"

		callCount := 0
		deleteFunc := func(ctx context.Context) error {
			callCount++
			return createRetryableError()
		}

		err := DeleteWithRetry(ctx, deleteFunc, opts)

		assert.Error(t, err, "Should fail after retries")
		assert.Equal(t, 3, callCount, "Should retry until max retries") // initial + 2 retries
	})

	t.Run("Non-retryable error", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		opts := DefaultDeleteWithRetryOptions()
		opts.MaxRetries = 2
		opts.RetryInterval = 50 * time.Millisecond
		opts.ResourceTypeName = "TestResource"
		opts.ResourceID = "test-123"

		callCount := 0
		deleteFunc := func(ctx context.Context) error {
			callCount++
			return createNonRetryableError()
		}

		err := DeleteWithRetry(ctx, deleteFunc, opts)

		assert.Error(t, err, "Should fail immediately")
		assert.Contains(t, err.Error(), "non-retryable error")
		assert.Equal(t, 1, callCount, "Should not retry")
	})
}

func TestDeleteWithRetry_UnknownErrorType(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()
	opts.MaxRetries = 2 // Limit retries to prevent long-running test
	opts.RetryInterval = 50 * time.Millisecond
	opts.ResourceTypeName = "TestResource"
	opts.ResourceID = "test-123"

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		// Return standard Go error - this will be treated as 500 (retryable)
		return fmt.Errorf("unknown error type")
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown error")
	// Standard Go errors are treated as 500 (retryable), so should retry until max attempts
	assert.Equal(t, 3, callCount, "Should retry for unknown error type (treated as 500)")
}

func TestDeleteWithRetry_ResourceLogging(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Test with empty resource names (should use defaults)
	opts := DeleteWithRetryOptions{
		MaxRetries:    1,
		RetryInterval: 50 * time.Millisecond,
		Operation:     "", // Should default to constants.TfOperationDelete
		// ResourceTypeName and ResourceID left empty
	}

	callCount := 0
	deleteFunc := func(ctx context.Context) error {
		callCount++
		return nil
	}

	err := DeleteWithRetry(ctx, deleteFunc, opts)

	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

// Benchmark tests
func BenchmarkDeleteWithRetry_Success(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()

	for i := 0; i < b.N; i++ {
		deleteFunc := func(ctx context.Context) error {
			return nil
		}

		err := DeleteWithRetry(ctx, deleteFunc, opts)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkDeleteWithRetry_WithRetries(b *testing.B) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	opts := DefaultDeleteWithRetryOptions()
	opts.MaxRetries = 2
	opts.RetryInterval = 10 * time.Millisecond

	for i := 0; i < b.N; i++ {
		callCount := 0
		deleteFunc := func(ctx context.Context) error {
			callCount++
			if callCount < 3 {
				return createRetryableError()
			}
			return nil
		}

		err := DeleteWithRetry(ctx, deleteFunc, opts)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}
