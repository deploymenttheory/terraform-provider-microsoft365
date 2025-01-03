package graphBetaWindowsPlatformScript

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
func (r *WindowsPlatformScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsPlatformScriptResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	// create resource
	requestResource, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*requestResource.GetId())

	// create assignments
	if object.Assignments != nil {
		requestAssignment, err := constructAssignment(ctx, &object)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for create method",
				fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			DeviceManagementScripts().
			ByDeviceManagementScriptId(object.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

	// Update state with the new resource ID
	respResource, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Get(context.Background(), nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}
	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	respAssignments, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create - Assignments Fetch", r.ReadPermissions)
		return
	}
	MapRemoteAssignmentStateToTerraform(ctx, &object, respAssignments)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation for device management script resources.
//
//   - Retrieves the current state from the read request
//   - Gets the base resource details from the API
//   - Maps the base resource details to Terraform state
//   - Gets the assignments configuration from the API
//   - Maps the assignments configuration to Terraform state
//
// The function ensures that all components (base resource and assignments)
// are properly read and mapped into the Terraform state, providing a complete view
// of the resource's current configuration on the server.
func (r *WindowsPlatformScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsPlatformScriptResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Read base resource
	respResource, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	// Read assignments
	respAssignments, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteAssignmentStateToTerraform(ctx, &object, respAssignments)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for Device Management Script resources.
//
//   - Deletes the existing resource
//   - Creates a new resource with the updated settings and assignments
//   - Updates the Terraform state with the new resource ID and settings
//
// NOTE: The Update operation is implemented as a delete followed by a create operation.
// This is because the Microsoft Graph Beta API does not support updating the settings of a device management script
// resource directly despite what the SDK and documentation may suggest.
// https://learn.microsoft.com/en-us/graph/api/intune-shared-devicemanagementscript-update?view=graph-rest-beta
// Produces 400 ODATA errors when attempting to update the resource settings using PATCH.
// Tested with custom PUT, SDK PATCH and custom POST requests, all failed.
func (r *WindowsPlatformScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object WindowsPlatformScriptResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	// Get current state
	var state WindowsPlatformScriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get planned changes
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// First delete the existing resource using the API but don't remove from state
	err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(state.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update - Delete", r.WritePermissions)
		return
	}

	// Then create new resource
	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	// create resource
	requestResource, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*requestResource.GetId())

	// create assignments
	if object.Assignments != nil {
		requestAssignment, err := constructAssignment(ctx, &object)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for create method",
				fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			DeviceManagementScripts().
			ByDeviceManagementScriptId(object.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

	// Update state with the new resource ID
	respResource, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Get(context.Background(), nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}
	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	respAssignments, err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create - Assignments Fetch", r.ReadPermissions)
		return
	}
	MapRemoteAssignmentStateToTerraform(ctx, &object, respAssignments)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for Device Management Script resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *WindowsPlatformScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsPlatformScriptResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
