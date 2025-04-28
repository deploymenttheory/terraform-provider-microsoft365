package graphBetaWindowsFeatureUpdateProfile

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Create handles the Create operation.
func (r *WindowsFeatureUpdateProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsFeatureUpdateProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

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
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	common.GraphSDKMutex.Lock()
	createdResource, err := r.client.
		DeviceManagement().
		WindowsFeatureUpdateProfiles().
		Post(ctx, requestBody, nil)
	common.GraphSDKMutex.Unlock()

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	if object.Assignments != nil && len(object.Assignments) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Assignments detected, constructing assignment request for policy ID: %s", object.ID.ValueString()))

		assignRequestBody, err := constructAssignments(ctx, &object)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Failed to construct assignments for policy: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			WindowsFeatureUpdateProfiles().
			ByWindowsFeatureUpdateProfileId(object.ID.ValueString()).
			Assign().
			Post(ctx, assignRequestBody, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "CreateAssignments", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully posted assignments for policy ID: %s", object.ID.ValueString()))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}, readResp)

	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State = readResp.State

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation for macos platform scripts resources.
//
//   - Retrieves the current state from the read request
//   - Gets the resource details including assignments from the API using expand
//   - Maps both resource and assignment details to Terraform state
//
// The function ensures all components are properly read and mapped into the
// Terraform state in a single API call, providing a complete view of the
// resource's current configuration on the server.
func (r *WindowsFeatureUpdateProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsFeatureUpdateProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	common.GraphSDKMutex.Lock()
	respResource, err := r.client.
		DeviceManagement().
		WindowsFeatureUpdateProfiles().
		ByWindowsFeatureUpdateProfileId(object.ID.ValueString()).
		Get(ctx, nil)
	common.GraphSDKMutex.Unlock()

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	assignmentsResp, err := r.client.
		DeviceManagement().
		WindowsFeatureUpdateProfiles().
		ByWindowsFeatureUpdateProfileId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteAssignmentsToTerraform(ctx, &object, assignmentsResp.GetValue())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for windows driver update profile resources.
//
// The function performs the following operations:
//   - Patches the existing script resource with updated settings using PATCH
//   - Updates assignments using POST if they are defined
//   - Retrieves the updated resource with expanded assignments
//   - Maps the remote state back to Terraform
//
// The Microsoft Graph Beta API supports direct updates of windows driver update profile resources
// through PATCH operations for the base resource, while assignments are handled through
// a separate POST operation to the assign endpoint. This allows for atomic updates
// of both the script properties and its assignments.
func (r *WindowsFeatureUpdateProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object WindowsFeatureUpdateProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	requestBody, err := constructResource(ctx, &object, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}
	common.GraphSDKMutex.Lock()
	_, err = r.client.
		DeviceManagement().
		WindowsFeatureUpdateProfiles().
		ByWindowsFeatureUpdateProfileId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)
	common.GraphSDKMutex.Unlock()

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	if object.Assignments != nil && len(object.Assignments) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Assignments detected, constructing assignment request for policy ID: %s", object.ID.ValueString()))

		assignRequestBody, err := constructAssignments(ctx, &object)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Failed to construct assignments for policy: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			WindowsFeatureUpdateProfiles().
			ByWindowsFeatureUpdateProfileId(object.ID.ValueString()).
			Assign().
			Post(ctx, assignRequestBody, nil)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "UpdateAssignments", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully posted assignments for policy ID: %s", object.ID.ValueString()))
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Update: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error verifying update",
			fmt.Sprintf("Failed to verify updated resource: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for windows driver update profile resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *WindowsFeatureUpdateProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsFeatureUpdateProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		WindowsFeatureUpdateProfiles().
		ByWindowsFeatureUpdateProfileId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
