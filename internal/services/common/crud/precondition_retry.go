package crud

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

const (
	defaultPreconditionMaxRetries    = 3
	defaultPreconditionRetryInterval = 5 * time.Second
	defaultPreconditionStatusCode    = 412
	defaultPreconditionErrorCode     = "PreconditionFailed"
)

// PreconditionRetryOptions configures WriteWithPreconditionRetry.
type PreconditionRetryOptions struct {
	// MaxRetries is the number of retries after the first attempt (default: 3, so at most 4 attempts).
	MaxRetries int
	// RetryInterval is the time to wait between attempts (default: 5 seconds).
	RetryInterval time.Duration
	// StatusCodes are the HTTP status codes that trigger a retry (default: 412).
	StatusCodes []int
	// ErrorCodes are the OData error codes that trigger a retry (default: "PreconditionFailed").
	// An error is retried when its status code OR its error code matches.
	ErrorCodes []string
	// Operation is the name of the write for logging (e.g. constants.TfOperationUpdate).
	Operation string
	// ResourceTypeName is the optional resource type name for logging.
	ResourceTypeName string
	// LogFields are extra fields added to each retry log entry (e.g. the parent policy ID).
	LogFields map[string]any
}

// DefaultPreconditionRetryOptions returns the options used by the network policy rule
// resources: 3 retries, 5 seconds apart, on HTTP 412 or OData error code PreconditionFailed.
func DefaultPreconditionRetryOptions() PreconditionRetryOptions {
	return PreconditionRetryOptions{
		MaxRetries:    defaultPreconditionMaxRetries,
		RetryInterval: defaultPreconditionRetryInterval,
		StatusCodes:   []int{defaultPreconditionStatusCode},
		ErrorCodes:    []string{defaultPreconditionErrorCode},
		Operation:     "Operation",
	}
}

// WriteWithPreconditionRetry calls write and retries it while Graph rejects it on a
// precondition, for example while a parent policy is still settling after another write.
//
// It returns nil as soon as write succeeds. An error whose status code and OData error
// code both miss the configured lists is returned at once. When the retries are used up it
// returns the last error. If the context ends while waiting, it returns an error wrapping
// both the last write error and the context error.
//
// A MaxRetries or RetryInterval of zero or less, and empty StatusCodes and ErrorCodes
// together, fall back to DefaultPreconditionRetryOptions.
func WriteWithPreconditionRetry(ctx context.Context, write func(ctx context.Context) error, opts PreconditionRetryOptions) error {
	defaults := DefaultPreconditionRetryOptions()
	if opts.MaxRetries <= 0 {
		opts.MaxRetries = defaults.MaxRetries
	}
	if opts.RetryInterval <= 0 {
		opts.RetryInterval = defaults.RetryInterval
	}
	if len(opts.StatusCodes) == 0 && len(opts.ErrorCodes) == 0 {
		opts.StatusCodes = defaults.StatusCodes
		opts.ErrorCodes = defaults.ErrorCodes
	}
	if opts.Operation == "" {
		opts.Operation = defaults.Operation
	}
	resourceType := opts.ResourceTypeName
	if resourceType == "" {
		resourceType = "resource"
	}

	for attempt := 0; ; attempt++ {
		err := write(ctx)
		if err == nil {
			return nil
		}

		errorInfo := errors.GraphError(ctx, err)
		if !slices.Contains(opts.StatusCodes, errorInfo.StatusCode) && !slices.Contains(opts.ErrorCodes, errorInfo.ErrorCode) {
			return err
		}
		if attempt >= opts.MaxRetries {
			return err
		}

		fields := map[string]any{
			"operation":     opts.Operation,
			"resource_type": resourceType,
			"status_code":   errorInfo.StatusCode,
			"error_code":    errorInfo.ErrorCode,
			"attempt":       attempt + 1,
			"max_retries":   opts.MaxRetries,
		}
		maps.Copy(fields, opts.LogFields)
		tflog.Warn(ctx, fmt.Sprintf("Retrying %s %s after Graph precondition failure", resourceType, opts.Operation), fields)

		timer := time.NewTimer(opts.RetryInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return fmt.Errorf("context ended while retrying %s %s after a precondition failure: %w: %w", resourceType, opts.Operation, err, ctx.Err())
		case <-timer.C:
		}
	}
}
