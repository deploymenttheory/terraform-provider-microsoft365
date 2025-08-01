package graphBetaWindowsUpdateRingAction

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	custom_requests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for Windows Update Ring Action resources.
func (r *WindowsUpdateRingActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsUpdateRingActionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Generate a unique ID for this action resource
	object.ID = types.StringValue(uuid.New().String())

	// Validate that update_ring_id is provided
	if object.UpdateRingId.IsNull() || object.UpdateRingId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Field",
			"update_ring_id is required but was not provided",
		)
		return
	}

	// Determine what actions need to be performed
	actionsToPerform := determineActionsToPerform(ctx, &object)
	if len(actionsToPerform) == 0 {
		resp.Diagnostics.AddError(
			"No Actions Specified",
			"At least one action must be set to true to perform actions on the Windows Update Ring",
		)
		return
	}

	// Perform each action
	var lastActionPerformed ActionType
	for _, action := range actionsToPerform {
		tflog.Debug(ctx, fmt.Sprintf("Performing action: %s", action.ActionType))

		err := r.performAction(ctx, object.UpdateRingId.ValueString(), action.ActionType)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		lastActionPerformed = action.ActionType
		tflog.Debug(ctx, fmt.Sprintf("Successfully performed action: %s", action.ActionType))
	}

	// Update metadata
	updateActionMetadata(ctx, &object, lastActionPerformed)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for Windows Update Ring Action resources.
func (r *WindowsUpdateRingActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsUpdateRingActionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// For action resources, we mainly just return the current state
	// The actual state of the update ring can be checked via the main update ring resource
	tflog.Debug(ctx, fmt.Sprintf("Reading action resource state for ID: %s", object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Update Ring Action resources.
func (r *WindowsUpdateRingActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsUpdateRingActionResourceModel
	var state WindowsUpdateRingActionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Determine what new actions need to be performed based on plan changes
	actionsToPerform := determineActionsToPerform(ctx, &plan)
	if len(actionsToPerform) == 0 {
		// No actions to perform, just update the description if changed
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	// Perform each new action
	var lastActionPerformed ActionType
	for _, action := range actionsToPerform {
		tflog.Debug(ctx, fmt.Sprintf("Performing updated action: %s", action.ActionType))

		err := r.performAction(ctx, plan.UpdateRingId.ValueString(), action.ActionType)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}

		lastActionPerformed = action.ActionType
		tflog.Debug(ctx, fmt.Sprintf("Successfully performed updated action: %s", action.ActionType))
	}

	// Update metadata
	updateActionMetadata(ctx, &plan, lastActionPerformed)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Windows Update Ring Action resources.
func (r *WindowsUpdateRingActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsUpdateRingActionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// For action resources, deletion is a no-op since these are one-time actions
	// The actions have already been performed and cannot be undone by deleting this resource
	tflog.Debug(ctx, fmt.Sprintf("Action resource %s deleted from Terraform state - no API cleanup needed for one-time actions", object.ID.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// performAction executes the specified action on the Windows Update Ring
func (r *WindowsUpdateRingActionResource) performAction(ctx context.Context, updateRingId string, actionType ActionType) error {
	tflog.Debug(ctx, fmt.Sprintf("Performing action %s on update ring %s", actionType, updateRingId))

	switch actionType {
	case ActionExtendFeatureUpdates:
		// POST request to extend endpoint using custom request
		endpoint := fmt.Sprintf("deviceManagement/deviceConfigurations/%s/microsoft.graph.windowsUpdateForBusinessConfiguration/extendFeatureUpdatesPause", updateRingId)

		config := custom_requests.PostRequestConfig{
			APIVersion:  custom_requests.GraphAPIBeta,
			Endpoint:    endpoint,
			RequestBody: nil, // No body needed for extend operation
		}

		err := custom_requests.PostRequestNoContent(ctx, r.client.GetAdapter(), config)
		if err != nil {
			return fmt.Errorf("failed to extend feature updates pause: %w", err)
		}

	default:
		// PATCH request for all other actions
		requestBody, err := constructPatchRequest(ctx, actionType)
		if err != nil {
			return fmt.Errorf("failed to construct PATCH request: %w", err)
		}

		_, err = r.client.
			DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(updateRingId).
			Patch(ctx, requestBody, nil)

		if err != nil {
			return fmt.Errorf("failed to perform PATCH action %s: %w", actionType, err)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully performed action %s", actionType))
	return nil
}
