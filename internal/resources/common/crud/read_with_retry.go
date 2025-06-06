package crud

import (
	"context"
	"fmt"
	"time"

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
	// Operation is the name of the operation for logging (e.g., "Create", "Update")
	Operation string
	// ResourceTypeName is the optional resource type name for logging
	ResourceTypeName string
}

// DefaultReadWithRetryOptions returns sensible default options for most use cases
func DefaultReadWithRetryOptions() ReadWithRetryOptions {
	return ReadWithRetryOptions{
		MaxRetries:    30,
		RetryInterval: 2 * time.Second,
		Operation:     "Operation",
	}
}

// StateContainer interface for anything that has a State field
type StateContainer interface {
	GetState() tfsdk.State
	SetState(tfsdk.State)
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

	tflog.Debug(ctx, fmt.Sprintf("Starting read with retry for %s operation", opts.Operation), map[string]interface{}{
		"resource_id":   resourceID,
		"resource_type": resourceType,
	})

	// Ensure we have reasonable defaults
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

		tflog.Debug(ctx, fmt.Sprintf("Read retry attempt %d/%d", attempt+1, opts.MaxRetries+1), map[string]interface{}{
			"resource_id":   resourceID,
			"resource_type": resourceType,
		})

		readResp := &resource.ReadResponse{State: stateContainer.GetState()}

		readFunc(ctx, readReq, readResp)

		if !readResp.Diagnostics.HasError() {
			tflog.Debug(ctx, fmt.Sprintf("Read successful on attempt %d", attempt+1), map[string]interface{}{
				"resource_id":   resourceID,
				"resource_type": resourceType,
			})
			stateContainer.SetState(readResp.State)
			return nil
		}

		lastErr = fmt.Errorf("error reading resource state after %s method on attempt %d: %s",
			opts.Operation, attempt+1, readResp.Diagnostics.Errors())

		if attempt < opts.MaxRetries {
			tflog.Debug(ctx, fmt.Sprintf("Read failed on attempt %d, waiting %s before retry", attempt+1, opts.RetryInterval), map[string]interface{}{
				"resource_id":   resourceID,
				"resource_type": resourceType,
			})

			select {
			case <-time.After(opts.RetryInterval):
			case <-ctx.Done():
				return fmt.Errorf("context cancelled during retry wait: %w", ctx.Err())
			}
		}
	}

	if lastErr != nil {
		return fmt.Errorf("failed to read resource state for %s after %d attempts: %w", resourceType, opts.MaxRetries+1, lastErr)
	}

	return fmt.Errorf("failed to read resource state for %s after %d attempts", resourceType, opts.MaxRetries+1)
}
