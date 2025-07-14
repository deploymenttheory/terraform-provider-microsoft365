package graphBetaWindowsAutopilotDeviceIdentity

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the Create operation.
func (r *WindowsAutopilotDeviceIdentityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsAutopilotDeviceIdentityResourceModel

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

	requestBody, err := constructResource(ctx, &object, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	// Handle user assignment if specified
	if object.UserAssignment != nil && !object.UserAssignment.UserPrincipalName.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf("Assigning user to device with ID: %s", object.ID.ValueString()))

		err := r.assignUser(ctx, object.ID.ValueString(), object.UserAssignment.UserPrincipalName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error assigning user to device",
				fmt.Sprintf("Could not assign user to device: %s: %s", ResourceName, err.Error()),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation.
func (r *WindowsAutopilotDeviceIdentityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsAutopilotDeviceIdentityResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := "Read"
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	respResource, err := r.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	// Handle user assignment information
	if respResource.GetUserPrincipalName() != nil && *respResource.GetUserPrincipalName() != "" {
		if object.UserAssignment == nil {
			object.UserAssignment = &UserAssignmentModel{}
		}
		object.UserAssignment.UserPrincipalName = types.StringValue(*respResource.GetUserPrincipalName())
		if respResource.GetAddressableUserName() != nil {
			object.UserAssignment.AddressableUserName = types.StringValue(*respResource.GetAddressableUserName())
		}
	} else {
		object.UserAssignment = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *WindowsAutopilotDeviceIdentityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsAutopilotDeviceIdentityResourceModel
	var state WindowsAutopilotDeviceIdentityResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

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

	requestBody, err := constructResource(ctx, &plan, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Handle user assignment changes
	if shouldUpdateUserAssignment(&state, &plan) {
		if plan.UserAssignment != nil && !plan.UserAssignment.UserPrincipalName.IsNull() {
			tflog.Debug(ctx, fmt.Sprintf("Assigning user to device with ID: %s", state.ID.ValueString()))
			err := r.assignUser(ctx, state.ID.ValueString(), plan.UserAssignment.UserPrincipalName.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error assigning user to device",
					fmt.Sprintf("Could not assign user to device: %s: %s", ResourceName, err.Error()),
				)
				return
			}
		} else if state.UserAssignment != nil && !state.UserAssignment.UserPrincipalName.IsNull() {
			tflog.Debug(ctx, fmt.Sprintf("Unassigning user from device with ID: %s", state.ID.ValueString()))
			err := r.unassignUser(ctx, state.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error unassigning user from device",
					fmt.Sprintf("Could not unassign user from device: %s: %s", ResourceName, err.Error()),
				)
				return
			}
		}
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation.
func (r *WindowsAutopilotDeviceIdentityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsAutopilotDeviceIdentityResourceModel

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

	// If there's a user assigned, unassign it first
	if object.UserAssignment != nil && !object.UserAssignment.UserPrincipalName.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf("Unassigning user from device with ID: %s before deletion", object.ID.ValueString()))
		err := r.unassignUser(ctx, object.ID.ValueString())
		if err != nil {
			// Log the error but continue with deletion
			tflog.Error(ctx, fmt.Sprintf("Error unassigning user from device before deletion: %s", err.Error()))
		}
	}

	err := r.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// assignUser assigns a user to a Windows Autopilot device
func (r *WindowsAutopilotDeviceIdentityResource) assignUser(ctx context.Context, deviceId string, userPrincipalName string) error {
	tflog.Debug(ctx, fmt.Sprintf("Assigning user %s to device %s", userPrincipalName, deviceId))

	requestBody := devicemanagement.NewWindowsAutopilotDeviceIdentitiesItemAssignUserToDevicePostRequestBody()
	requestBody.SetUserPrincipalName(&userPrincipalName)

	err := r.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(deviceId).
		AssignUserToDevice().
		Post(ctx, requestBody, nil)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error assigning user to device: %s", err.Error()))
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned user %s to device %s", userPrincipalName, deviceId))
	return nil
}

// unassignUser unassigns a user from a Windows Autopilot device
func (r *WindowsAutopilotDeviceIdentityResource) unassignUser(ctx context.Context, deviceId string) error {
	tflog.Debug(ctx, fmt.Sprintf("Unassigning user from device %s", deviceId))

	err := r.client.
		DeviceManagement().
		WindowsAutopilotDeviceIdentities().
		ByWindowsAutopilotDeviceIdentityId(deviceId).
		UnassignUserFromDevice().
		Post(ctx, nil)

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error unassigning user from device: %s", err.Error()))
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully unassigned user from device %s", deviceId))
	return nil
}

// shouldUpdateUserAssignment determines if the user assignment needs to be updated
func shouldUpdateUserAssignment(state *WindowsAutopilotDeviceIdentityResourceModel, plan *WindowsAutopilotDeviceIdentityResourceModel) bool {
	// Case 1: State has no user assignment but plan does
	if (state.UserAssignment == nil || state.UserAssignment.UserPrincipalName.IsNull()) &&
		(plan.UserAssignment != nil && !plan.UserAssignment.UserPrincipalName.IsNull()) {
		return true
	}

	// Case 2: State has user assignment but plan doesn't
	if (state.UserAssignment != nil && !state.UserAssignment.UserPrincipalName.IsNull()) &&
		(plan.UserAssignment == nil || plan.UserAssignment.UserPrincipalName.IsNull()) {
		return true
	}

	// Case 3: Both have user assignments but they're different
	if state.UserAssignment != nil && plan.UserAssignment != nil &&
		!state.UserAssignment.UserPrincipalName.IsNull() && !plan.UserAssignment.UserPrincipalName.IsNull() &&
		state.UserAssignment.UserPrincipalName.ValueString() != plan.UserAssignment.UserPrincipalName.ValueString() {
		return true
	}

	return false
}
