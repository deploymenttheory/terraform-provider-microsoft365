package crud

import (
	"context"
	"fmt"
	"strings"
	"time"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ReadWithRetryOptions configures the retry behavior for reading resource state
type ReadWithRetryOptions struct {
	// MaxRetries is the maximum number of retry attempts (default: 30)
	MaxRetries int
	// RetryInterval is the time to wait between retries (default: 2 seconds)
	RetryInterval time.Duration
	// Operation is the name of the operation for logging (e.g., constants.TfOperationCreate, constants.TfOperationUpdate)
	Operation string
	// ResourceTypeName is the optional resource type name for logging
	ResourceTypeName string
	// ConsistencyPredicate is an optional function that inspects the read state and returns
	// true when the state is consistent with the preceding write. When non-nil, ReadWithRetry
	// will continue retrying even on a successful (no-error) read until the predicate returns
	// true or the retry limit is reached. Pass nil for the default fast-exit behaviour.
	//
	// This is the recommended approach for Microsoft Entra eventual consistency: trust the
	// successful write, then poll until the read reflects the written state, treating stale
	// reads as transient rather than using fixed sleeps.
	// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
	ConsistencyPredicate func(ctx context.Context, state tfsdk.State) bool
}

// DefaultReadWithRetryOptions returns sensible default options for most use cases
func DefaultReadWithRetryOptions() ReadWithRetryOptions {
	return ReadWithRetryOptions{
		MaxRetries:    30,
		RetryInterval: 2 * time.Second,
		Operation:     "Operation",
	}
}

// StateContainer interface for anything that has a State field and optionally an Identity field
type StateContainer interface {
	GetState() tfsdk.State
	SetState(tfsdk.State)
	GetIdentity() any
	SetIdentity(any)
}

// CreateResponseContainer wraps resource.CreateResponse to implement StateContainer
type CreateResponseContainer struct {
	*resource.CreateResponse
}

func (c *CreateResponseContainer) GetState() tfsdk.State {
	return c.State
}

func (c *CreateResponseContainer) SetState(state tfsdk.State) {
	c.State = state
}

func (c *CreateResponseContainer) GetIdentity() any {
	return c.Identity
}

func (c *CreateResponseContainer) SetIdentity(identity any) {
	if id, ok := identity.(*tfsdk.ResourceIdentity); ok {
		c.Identity = id
	}
}

// UpdateResponseContainer wraps resource.UpdateResponse to implement StateContainer
type UpdateResponseContainer struct {
	*resource.UpdateResponse
}

func (c *UpdateResponseContainer) GetState() tfsdk.State {
	return c.State
}

func (c *UpdateResponseContainer) SetState(state tfsdk.State) {
	c.State = state
}

func (c *UpdateResponseContainer) GetIdentity() any {
	return c.Identity
}

func (c *UpdateResponseContainer) SetIdentity(identity any) {
	if id, ok := identity.(*tfsdk.ResourceIdentity); ok {
		c.Identity = id
	}
}

// extractResourceID attempts to extract the ID from the state for logging purposes
func extractResourceID(ctx context.Context, state tfsdk.State) (result string) {

	defer func() {
		if r := recover(); r != nil {

			result = "unknown"
		}
	}()

	if state.Raw.IsNull() || !state.Raw.IsKnown() {
		return "unknown"
	}

	var idValue types.String
	diags := state.GetAttribute(ctx, path.Root("id"), &idValue)
	if diags.HasError() || idValue.IsNull() || idValue.IsUnknown() {
		return "unknown"
	}
	return idValue.ValueString()
}

// extractErrorFromDiagnostics analyzes Terraform diagnostics to extract HTTP error information
// for intelligent retry decisions
func extractErrorFromDiagnostics(diagnostics diag.Diagnostics) errors.GraphErrorInfo {
	if !diagnostics.HasError() {
		return errors.GraphErrorInfo{}
	}

	// Iterate through diagnostics to find HTTP error information
	for _, d := range diagnostics.Errors() {
		summary := d.Summary()
		detail := d.Detail()

		errorText := summary + " " + detail
		errorTextLower := strings.ToLower(errorText)

		if strings.Contains(errorTextLower, "404") || strings.Contains(errorTextLower, "not found") {
			return errors.GraphErrorInfo{StatusCode: 404, ErrorCode: "NotFound"}
		}
		if strings.Contains(errorTextLower, "400") || strings.Contains(errorTextLower, "bad request") {
			return errors.GraphErrorInfo{StatusCode: 400, ErrorCode: "BadRequest"}
		}
		if strings.Contains(errorTextLower, "401") || strings.Contains(errorTextLower, "unauthorized") {
			return errors.GraphErrorInfo{StatusCode: 401, ErrorCode: "Unauthorized"}
		}
		if strings.Contains(errorTextLower, "403") || strings.Contains(errorTextLower, "forbidden") {
			return errors.GraphErrorInfo{StatusCode: 403, ErrorCode: "Forbidden"}
		}
		if strings.Contains(errorTextLower, "409") || strings.Contains(errorTextLower, "conflict") {
			return errors.GraphErrorInfo{StatusCode: 409, ErrorCode: "Conflict"}
		}
		if strings.Contains(errorTextLower, "423") || strings.Contains(errorTextLower, "locked") {
			return errors.GraphErrorInfo{StatusCode: 423, ErrorCode: "Locked"}
		}
		if strings.Contains(errorTextLower, "429") || strings.Contains(errorTextLower, "too many requests") || strings.Contains(errorTextLower, "throttl") {
			return errors.GraphErrorInfo{StatusCode: 429, ErrorCode: "TooManyRequests"}
		}
		if strings.Contains(errorTextLower, "500") || strings.Contains(errorTextLower, "internal server error") {
			return errors.GraphErrorInfo{StatusCode: 500, ErrorCode: "InternalServerError"}
		}
		if strings.Contains(errorTextLower, "502") || strings.Contains(errorTextLower, "bad gateway") {
			return errors.GraphErrorInfo{StatusCode: 502, ErrorCode: "BadGateway"}
		}
		if strings.Contains(errorTextLower, "503") || strings.Contains(errorTextLower, "service unavailable") {
			return errors.GraphErrorInfo{StatusCode: 503, ErrorCode: "ServiceUnavailable"}
		}
		if strings.Contains(errorTextLower, "504") || strings.Contains(errorTextLower, "gateway timeout") {
			return errors.GraphErrorInfo{StatusCode: 504, ErrorCode: "GatewayTimeout"}
		}

		// Look for specific Graph API error patterns
		if strings.Contains(errorTextLower, "service unavailable") {
			return errors.GraphErrorInfo{StatusCode: 503, ErrorCode: "ServiceUnavailable"}
		}
		if strings.Contains(errorTextLower, "request throttled") {
			return errors.GraphErrorInfo{StatusCode: 429, ErrorCode: "RequestThrottled"}
		}
		if strings.Contains(errorTextLower, "resource not found") {
			return errors.GraphErrorInfo{StatusCode: 404, ErrorCode: "ResourceNotFound"}
		}
		if strings.Contains(errorTextLower, "network error") {
			return errors.GraphErrorInfo{StatusCode: 500, ErrorCode: "NetworkError"}
		}
		if strings.Contains(errorTextLower, "timeout") {
			return errors.GraphErrorInfo{StatusCode: 504, ErrorCode: "RequestTimeout"}
		}
	}

	return errors.GraphErrorInfo{}
}

// ReadWithRetry executes a read operation with retry logic within the context timeout
// It repeatedly calls the provided read function until success or context timeout
func ReadWithRetry(
	ctx context.Context,
	readFunc func(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse),
	readReq resource.ReadRequest,
	stateContainer StateContainer,
	opts ReadWithRetryOptions,
) error {
	resourceID := extractResourceID(ctx, stateContainer.GetState())
	resourceType := opts.ResourceTypeName
	if resourceType == "" {
		resourceType = "resource"
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting read with retry for %s operation", opts.Operation), map[string]any{
		"resource_id":   resourceID,
		"resource_type": resourceType,
	})

	if opts.MaxRetries <= 0 {
		opts.MaxRetries = 30
	}
	if opts.RetryInterval <= 0 {
		opts.RetryInterval = 2 * time.Second
	}
	if opts.Operation == "" {
		opts.Operation = "Operation"
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

	tflog.Debug(ctx, fmt.Sprintf("Will attempt up to %d retries with %s intervals", opts.MaxRetries, opts.RetryInterval), map[string]any{
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
			tflog.Debug(ctx, "Insufficient time remaining for another retry attempt", map[string]any{
				"resource_id":   resourceID,
				"resource_type": resourceType,
			})
			break
		}

		tflog.Debug(ctx, fmt.Sprintf("Read retry attempt %d/%d", attempt+1, opts.MaxRetries+1), map[string]any{
			"resource_id":   resourceID,
			"resource_type": resourceType,
		})

		readResp := &resource.ReadResponse{
			State: stateContainer.GetState(),
		}

		if containerIdentity := stateContainer.GetIdentity(); containerIdentity != nil {
			if identityState, ok := containerIdentity.(*tfsdk.ResourceIdentity); ok {
				readResp.Identity = identityState
			}
		}

		ctxWithOp := context.WithValue(ctx, "retry_operation", opts.Operation)
		readFunc(ctxWithOp, readReq, readResp)

		if !readResp.Diagnostics.HasError() {
			if opts.ConsistencyPredicate == nil || opts.ConsistencyPredicate(ctxWithOp, readResp.State) {
				tflog.Debug(ctx, fmt.Sprintf("Read successful on attempt %d", attempt+1), map[string]any{
					"resource_id":   resourceID,
					"resource_type": resourceType,
				})
				stateContainer.SetState(readResp.State)
				// Propagate Identity back to the container
				if readResp.Identity != nil {
					stateContainer.SetIdentity(readResp.Identity)
				}
				return nil
			}

			// Read succeeded but the state does not yet reflect the written change.
			// This is the eventual-consistency case: the API returned 2xx but replication
			// across replicas has not yet completed. Continue retrying.
			lastErr = fmt.Errorf("consistency predicate not yet satisfied after %s on attempt %d — awaiting propagation",
				opts.Operation, attempt+1)
			tflog.Debug(ctx, fmt.Sprintf("Read succeeded on attempt %d but consistency predicate not satisfied, retrying in %s",
				attempt+1, opts.RetryInterval), map[string]any{
				"resource_id":   resourceID,
				"resource_type": resourceType,
			})
			if attempt < opts.MaxRetries {
				select {
				case <-time.After(opts.RetryInterval):
				case <-ctx.Done():
					return fmt.Errorf("context cancelled during consistency wait: %w", ctx.Err())
				}
			}
			continue
		}

		lastErr = fmt.Errorf("error reading resource state after %s method on attempt %d: %s",
			opts.Operation, attempt+1, readResp.Diagnostics.Errors())

		errorInfo := extractErrorFromDiagnostics(readResp.Diagnostics)

		if errors.IsNonRetryableReadError(&errorInfo) {
			tflog.Error(ctx, fmt.Sprintf("Read failed on attempt %d (non-retryable error)", attempt+1), map[string]any{
				"resource_id":   resourceID,
				"resource_type": resourceType,
				"status_code":   errorInfo.StatusCode,
				"error_code":    errorInfo.ErrorCode,
				"diagnostics":   readResp.Diagnostics.Errors(),
			})
			return fmt.Errorf("read operation failed with non-retryable error: %w", lastErr)
		}

		if errors.IsRetryableReadError(&errorInfo) {
			if attempt < opts.MaxRetries {
				tflog.Warn(ctx, fmt.Sprintf("Read failed on attempt %d (retryable error), waiting %s before retry", attempt+1, opts.RetryInterval), map[string]any{
					"resource_id":   resourceID,
					"resource_type": resourceType,
					"status_code":   errorInfo.StatusCode,
					"error_code":    errorInfo.ErrorCode,
					"diagnostics":   readResp.Diagnostics.Errors(),
				})

				select {
				case <-time.After(opts.RetryInterval):
				case <-ctx.Done():
					return fmt.Errorf("context cancelled during retry wait: %w", ctx.Err())
				}
			} else {
				tflog.Error(ctx, fmt.Sprintf("Read failed on final attempt %d", attempt+1), map[string]any{
					"resource_id":   resourceID,
					"resource_type": resourceType,
					"status_code":   errorInfo.StatusCode,
					"error_code":    errorInfo.ErrorCode,
					"diagnostics":   readResp.Diagnostics.Errors(),
				})
			}
		} else {
			if attempt < opts.MaxRetries {
				tflog.Debug(ctx, fmt.Sprintf("Read failed on attempt %d (unknown error type, continuing retry)", attempt+1), map[string]any{
					"resource_id":   resourceID,
					"resource_type": resourceType,
					"diagnostics":   readResp.Diagnostics.Errors(),
				})

				select {
				case <-time.After(opts.RetryInterval):
				case <-ctx.Done():
					return fmt.Errorf("context cancelled during retry wait: %w", ctx.Err())
				}
			}
		}
	}

	if lastErr != nil {
		return fmt.Errorf("failed to read resource state for %s after %d attempts: %w", resourceType, opts.MaxRetries+1, lastErr)
	}

	return fmt.Errorf("failed to read resource state for %s after %d attempts", resourceType, opts.MaxRetries+1)
}
