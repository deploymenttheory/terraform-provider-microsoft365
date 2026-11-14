package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// taskArgumentsInvalidMessage is the Graph error returned when a task argument references an
// object Graph cannot resolve. It is also returned for a directory object created moments
// earlier (for example a group used as a groupID argument) until it has replicated.
const taskArgumentsInvalidMessage = "task arguments missing or invalid"

// taskArgumentRetryTimeout bounds how long a write is retried while task argument references replicate.
// Probes against a tenant showed a new group is accepted about 10 seconds after creation.
const taskArgumentRetryTimeout = 90 * time.Second

// taskArgumentRetryInterval is the wait between retries of a write rejected for task arguments.
const taskArgumentRetryInterval = 5 * time.Second

// writeWithTaskArgumentRetry runs a workflow write (POST or createNewVersion) and retries it while
// Graph rejects it with "Task arguments missing or invalid". Any other error is returned at once.
// When the retry window ends, the last Graph error is returned unchanged so it can be passed to
// HandleKiotaGraphError.
func writeWithTaskArgumentRetry(ctx context.Context, write func(ctx context.Context) error) error {
	var lastErr error

	retryCtx, cancel := context.WithTimeout(ctx, taskArgumentRetryTimeout)
	defer cancel()

	pollErr := crud.PollUntil(retryCtx, taskArgumentRetryInterval, func(ctx context.Context) (bool, error) {
		err := write(ctx)
		if err == nil {
			return true, nil
		}
		info := errors.GraphError(ctx, err)
		if info.StatusCode == http.StatusBadRequest && strings.Contains(strings.ToLower(info.ErrorMessage), taskArgumentsInvalidMessage) {
			tflog.Debug(ctx, "Graph rejected task arguments, retrying in case a referenced object has not replicated yet", map[string]any{"error": info.ErrorMessage})
			lastErr = err
			return false, err
		}
		return false, &crud.FatalPollError{Err: err}
	})

	if pollErr != nil && lastErr != nil && retryCtx.Err() != nil {
		return lastErr
	}

	return pollErr //nolint:wrapcheck // a fatal Graph error is passed unchanged to HandleKiotaGraphError
}

// versionedAttributes can only be changed with createNewVersion; PATCH rejects them.
var versionedAttributes = []string{"category", "execution_conditions", "tasks"}

// basicAttributes are the only properties Graph accepts on PATCH.
var basicAttributes = []string{"display_name", "description", "is_enabled", "is_scheduling_enabled"}

// updateKind describes which Graph call an update needs.
type updateKind int

const (
	updateNone updateKind = iota
	updateBasic
	updateNewVersion
)

// planUpdateKind compares plan and state and returns whether the update needs a new workflow
// version, a PATCH of the basic properties, or no Graph write at all (for example when only
// timeouts changed).
func planUpdateKind(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (updateKind, diag.Diagnostics) {
	var diags diag.Diagnostics

	for _, name := range versionedAttributes {
		planValue, stateValue, d := attributePair(ctx, plan, state, name)
		diags.Append(d...)
		if diags.HasError() {
			return updateNone, diags
		}
		if !planValue.Equal(stateValue) {
			return updateNewVersion, diags
		}
	}

	for _, name := range basicAttributes {
		planValue, stateValue, d := attributePair(ctx, plan, state, name)
		diags.Append(d...)
		if diags.HasError() {
			return updateNone, diags
		}
		if !planValue.Equal(stateValue) {
			return updateBasic, diags
		}
	}

	return updateNone, diags
}

func attributePair(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, name string) (attr.Value, attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics
	var planValue, stateValue attr.Value

	switch name {
	case "execution_conditions":
		var p, s types.Object
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		planValue, stateValue = p, s
	case "tasks":
		var p, s types.List
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		planValue, stateValue = p, s
	case "is_enabled", "is_scheduling_enabled":
		var p, s types.Bool
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		planValue, stateValue = p, s
	default:
		var p, s types.String
		diags.Append(plan.GetAttribute(ctx, path.Root(name), &p)...)
		diags.Append(state.GetAttribute(ctx, path.Root(name), &s)...)
		planValue, stateValue = p, s
	}

	return planValue, stateValue, diags
}
