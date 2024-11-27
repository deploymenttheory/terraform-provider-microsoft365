package graphBetaDeviceShellScript

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

var (
	// mutex needed to lock Create requests during parallel runs to avoid overwhelming api and resulting in stating issues
	mu sync.Mutex

	// object is the resource model for the device management script resource
	object DeviceShellScriptResourceModel
)

// Create handles the Create operation.
func (r *DeviceShellScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	mu.Lock()
	defer mu.Unlock()

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
		DeviceShellScripts().
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
			DeviceShellScripts().
			ByDeviceShellScriptId(object.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

	// resource and assignments are found within the same call
	respResource, err := r.client.
		DeviceManagement().
		DeviceShellScripts().
		ByDeviceShellScriptId(object.ID.ValueString()).
		Get(context.Background(), &devicemanagement.DeviceShellScriptsDeviceShellScriptItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DeviceShellScriptsDeviceShellScriptItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
			},
		})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation for device management script resources.
//
//   - Retrieves the current state from the read request
//   - Gets the resource details including assignments from the API using expand
//   - Maps both resource and assignment details to Terraform state
//
// The function ensures all components are properly read and mapped into the
// Terraform state in a single API call, providing a complete view of the
// resource's current configuration on the server.
func (r *DeviceShellScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	// Read resource with expanded assignments
	respResource, err := r.client.
		DeviceManagement().
		DeviceShellScripts().
		ByDeviceShellScriptId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.DeviceShellScriptsDeviceShellScriptItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DeviceShellScriptsDeviceShellScriptItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
			},
		})

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

// Update handles the Update operation for Device Management Script resources.
//
// The function performs the following operations:
//   - Patches the existing script resource with updated settings using PATCH
//   - Updates assignments using POST if they are defined
//   - Retrieves the updated resource with expanded assignments
//   - Maps the remote state back to Terraform
//
// The Microsoft Graph Beta API supports direct updates of device shell script resources
// through PATCH operations for the base resource, while assignments are handled through
// a separate POST operation to the assign endpoint. This allows for atomic updates
// of both the script properties and its assignments.
func (r *DeviceShellScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	var state DeviceShellScriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
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

	_, err = r.client.
		DeviceManagement().
		DeviceShellScripts().
		ByDeviceShellScriptId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	if object.Assignments != nil {
		requestAssignment, err := constructAssignment(ctx, &object)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for update method",
				fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			DeviceShellScripts().
			ByDeviceShellScriptId(state.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update - Assignments", r.WritePermissions)
			return
		}
	}

	respResource, err := r.client.
		DeviceManagement().
		DeviceShellScripts().
		ByDeviceShellScriptId(state.ID.ValueString()).
		Get(ctx, &devicemanagement.DeviceShellScriptsDeviceShellScriptItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DeviceShellScriptsDeviceShellScriptItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
			},
		})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update - Get", r.WritePermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for Device Management Script resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *DeviceShellScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

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
		DeviceShellScripts().
		ByDeviceShellScriptId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.ReadPermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
