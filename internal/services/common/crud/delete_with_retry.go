package crud

import (
	"context"
	"fmt"
	"time"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// DeleteWithRetryOptions configures the retry behavior for delete operations
type DeleteWithRetryOptions struct {
	// MaxRetries is the maximum number of retry attempts (default: 10)
	MaxRetries int
	// RetryInterval is the time to wait between retries (default: 30 seconds)
	RetryInterval time.Duration
	// Operation is the name of the operation for logging (e.g., "Delete")
	Operation string
	// ResourceTypeName is the optional resource type name for logging
	ResourceTypeName string
	// ResourceID is the resource ID for logging
	ResourceID string
}

// DefaultDeleteWithRetryOptions returns sensible default options for delete operations
func DefaultDeleteWithRetryOptions() DeleteWithRetryOptions {
	return DeleteWithRetryOptions{
		MaxRetries:    10,
		RetryInterval: 30 * time.Second,
		Operation:     "Delete",
	}
}

// DeleteWithRetry executes a delete operation with retry logic within the context timeout
// It repeatedly calls the provided delete function until success or context timeout
func DeleteWithRetry(
	ctx context.Context,
	deleteFunc func(ctx context.Context) error,
	opts DeleteWithRetryOptions,
) error {
	resourceType := opts.ResourceTypeName
	if resourceType == "" {
		resourceType = "resource"
	}

	resourceID := opts.ResourceID
	if resourceID == "" {
		resourceID = "unknown"
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting delete with retry for %s operation", opts.Operation), map[string]interface{}{
		"resource_id":   resourceID,
		"resource_type": resourceType,
	})

	// Ensure we have reasonable defaults
	if opts.MaxRetries <= 0 {
		opts.MaxRetries = 10
	}
	if opts.RetryInterval <= 0 {
		opts.RetryInterval = 30 * time.Second
	}
	if opts.Operation == "" {
		opts.Operation = "Delete"
	}

	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		return fmt.Errorf("context must have a deadline for retry operations")
	}

	timeRemaining := time.Until(deadline) - time.Second
	if timeRemaining <= 0 {
		return fmt.Errorf("insufficient time remaining in context for retry operation")
	}

	maxPossibleRetries := int(timeRemaining / opts.RetryInterval)
	if maxPossibleRetries < opts.MaxRetries {
		opts.MaxRetries = maxPossibleRetries
	}

	tflog.Debug(ctx, fmt.Sprintf("Will attempt up to %d retries with %s intervals", opts.MaxRetries, opts.RetryInterval), map[string]interface{}{
		"resource_id":   resourceID,
		"resource_type": resourceType,
	})

	var lastErr error
	for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during retry attempt %d: %w", attempt, ctx.Err())
		default:
		}

		if time.Until(deadline) < opts.RetryInterval {
			tflog.Debug(ctx, "Insufficient time remaining for another retry attempt", map[string]interface{}{
				"resource_id":   resourceID,
				"resource_type": resourceType,
			})
			break
		}

		tflog.Debug(ctx, fmt.Sprintf("Delete retry attempt %d/%d", attempt+1, opts.MaxRetries+1), map[string]interface{}{
			"resource_id":   resourceID,
			"resource_type": resourceType,
		})

		err := deleteFunc(ctx)

		if err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Delete successful on attempt %d", attempt+1), map[string]interface{}{
				"resource_id":   resourceID,
				"resource_type": resourceType,
			})
			return nil
		}

		lastErr = err

		// Extract error information and check if retryable
		errorInfo := errors.GraphError(ctx, err)

		// Check for non-retryable errors first (permanent failures or success)
		if errors.IsNonRetryableDeleteError(&errorInfo) {
			tflog.Error(ctx, fmt.Sprintf("Delete failed on attempt %d (non-retryable error)", attempt+1), map[string]interface{}{
				"resource_id":   resourceID,
				"resource_type": resourceType,
				"status_code":   errorInfo.StatusCode,
				"error_code":    errorInfo.ErrorCode,
				"error":         err.Error(),
			})
			return fmt.Errorf("delete operation failed with non-retryable error: %w", err)
		}

		// Check if this error should trigger a retry
		if errors.IsRetryableDeleteError(&errorInfo) {
			if attempt < opts.MaxRetries {
				tflog.Warn(ctx, fmt.Sprintf("Delete failed on attempt %d (retryable error), waiting %s before retry", attempt+1, opts.RetryInterval), map[string]interface{}{
					"resource_id":   resourceID,
					"resource_type": resourceType,
					"status_code":   errorInfo.StatusCode,
					"error_code":    errorInfo.ErrorCode,
					"error":         err.Error(),
				})

				select {
				case <-time.After(opts.RetryInterval):
				case <-ctx.Done():
					return fmt.Errorf("context cancelled during retry wait: %w", ctx.Err())
				}
			} else {
				tflog.Error(ctx, fmt.Sprintf("Delete failed on final attempt %d", attempt+1), map[string]interface{}{
					"resource_id":   resourceID,
					"resource_type": resourceType,
					"status_code":   errorInfo.StatusCode,
					"error_code":    errorInfo.ErrorCode,
					"error":         err.Error(),
				})
			}
		} else {
			// Unknown error type, fail immediately
			tflog.Error(ctx, fmt.Sprintf("Delete failed on attempt %d (unknown error type)", attempt+1), map[string]interface{}{
				"resource_id":   resourceID,
				"resource_type": resourceType,
				"error":         err.Error(),
			})
			return fmt.Errorf("delete operation failed with unknown error: %w", err)
		}
	}

	if lastErr != nil {
		return fmt.Errorf("failed to delete resource %s after %d attempts: %w", resourceType, opts.MaxRetries+1, lastErr)
	}

	return fmt.Errorf("failed to delete resource %s after %d attempts", resourceType, opts.MaxRetries+1)
}
