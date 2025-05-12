package graphBetaWindowsDriverUpdateInventory

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Create handles the Create operation.
func (r *WindowsDriverUpdateInventoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsDriverUpdateInventoryResourceModel

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

	// Validate that Windows Driver Update Profile ID is provided
	if object.WindowsDriverUpdateProfileID.IsNull() || object.WindowsDriverUpdateProfileID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to create a driver inventory.",
		)
		return
	}

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	constants.GraphSDKMutex.Lock()
	createdResource, err := r.client.
		DeviceManagement().
		WindowsDriverUpdateProfiles().
		ByWindowsDriverUpdateProfileId(object.WindowsDriverUpdateProfileID.ValueString()).
		DriverInventories().
		Post(ctx, requestBody, nil)
	constants.GraphSDKMutex.Unlock()

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

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

// Read handles the Read operation.
func (r *WindowsDriverUpdateInventoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsDriverUpdateInventoryResourceModel

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

	// Validate that Windows Driver Update Profile ID is provided
	if object.WindowsDriverUpdateProfileID.IsNull() || object.WindowsDriverUpdateProfileID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to read a driver inventory.",
		)
		return
	}

	constants.GraphSDKMutex.Lock()
	respResource, err := r.client.
		DeviceManagement().
		WindowsDriverUpdateProfiles().
		ByWindowsDriverUpdateProfileId(object.WindowsDriverUpdateProfileID.ValueString()).
		DriverInventories().
		ByWindowsDriverUpdateInventoryId(object.ID.ValueString()).
		Get(ctx, nil)
	constants.GraphSDKMutex.Unlock()

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *WindowsDriverUpdateInventoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object WindowsDriverUpdateInventoryResourceModel

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

	// Validate that Windows Driver Update Profile ID is provided
	if object.WindowsDriverUpdateProfileID.IsNull() || object.WindowsDriverUpdateProfileID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to update a driver inventory.",
		)
		return
	}

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	constants.GraphSDKMutex.Lock()
	_, err = r.client.
		DeviceManagement().
		WindowsDriverUpdateProfiles().
		ByWindowsDriverUpdateProfileId(object.WindowsDriverUpdateProfileID.ValueString()).
		DriverInventories().
		ByWindowsDriverUpdateInventoryId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)
	constants.GraphSDKMutex.Unlock()

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Update Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource update",
			fmt.Sprintf("Failed to verify resource update: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *WindowsDriverUpdateInventoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsDriverUpdateInventoryResourceModel

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

	// Validate that Windows Driver Update Profile ID is provided
	if object.WindowsDriverUpdateProfileID.IsNull() || object.WindowsDriverUpdateProfileID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to delete a driver inventory.",
		)
		return
	}

	constants.GraphSDKMutex.Lock()
	err := r.client.
		DeviceManagement().
		WindowsDriverUpdateProfiles().
		ByWindowsDriverUpdateProfileId(object.WindowsDriverUpdateProfileID.ValueString()).
		DriverInventories().
		ByWindowsDriverUpdateInventoryId(object.ID.ValueString()).
		Delete(ctx, nil)
	constants.GraphSDKMutex.Unlock()

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
