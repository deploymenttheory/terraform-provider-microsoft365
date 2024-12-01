package graphBetaSettingsCatalog

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client/graphcustom"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/retry"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for Settings Catalog resources.
//
//   - Retrieves the planned configuration from the create request
//   - Constructs the resource request body from the plan
//   - Sends POST request to create the base resource and settings
//   - Captures the new resource ID from the response
//   - Constructs and sends assignment configuration if specified
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API
//   - Updates the final state with the fresh data from the API
//
// The function ensures that both the settings catalog profile and its assignments
// (if specified) are created properly. The settings must be defined during creation
// as they are required for a successful deployment, while assignments are optional.
func (r *SettingsCatalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object SettingsCatalogProfileResourceModel

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

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	err = retry.RetryableIntuneOperation(ctx, "create resource", retry.IntuneWrite, func() error {
		var opErr error
		requestBody, opErr = r.client.
			DeviceManagement().
			ConfigurationPolicies().
			Post(ctx, requestBody, nil)
		return opErr
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*requestBody.GetId())

	if object.Assignments != nil {
		requestAssignment, err := constructAssignment(ctx, &object)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for create method",
				fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
			)
			return
		}

		err = retry.RetryableAssignmentOperation(ctx, "create assignment", func() error {
			_, err := r.client.
				DeviceManagement().
				ConfigurationPolicies().
				ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
				Assign().
				Post(ctx, requestAssignment, nil)
			return err
		})

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
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

// Read handles the Read operation for Settings Catalog resources.
//
//   - Retrieves the current state from the read request
//   - Gets the base resource details from the API
//   - Maps the base resource details to Terraform state
//   - Gets the settings configuration from the API
//   - Maps the settings configuration to Terraform state
//   - Gets the assignments configuration from the API
//   - Maps the assignments configuration to Terraform state
//
// The function ensures that all components (base resource, settings, and assignments)
// are properly read and mapped into the Terraform state, providing a complete view
// of the resource's current configuration on the server.
func (r *SettingsCatalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object SettingsCatalogProfileResourceModel

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

	err := retry.RetryableIntuneOperation(ctx, "read resource", retry.IntuneRead, func() error {
		respResource, err := r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			Get(ctx, nil)
		if err != nil {
			return err
		}
		MapRemoteResourceStateToTerraform(ctx, &object, respResource)
		return nil
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	settingsConfig := graphcustom.GetRequestConfig{
		APIVersion:        graphcustom.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		EndpointSuffix:    "/settings",
		ResourceIDPattern: "('id')",
		ResourceID:        object.ID.ValueString(),
		QueryParameters: map[string]string{
			"$expand": "children",
		},
	}

	err = retry.RetryableIntuneOperation(ctx, "read resource", retry.IntuneRead, func() error {
		respSettings, err := graphcustom.GetRequestByResourceId(
			ctx,
			r.client.GetAdapter(),
			settingsConfig,
		)
		if err != nil {
			return err
		}
		MapRemoteSettingsStateToTerraform(ctx, &object, respSettings)
		return nil
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	err = retry.RetryableAssignmentOperation(ctx, "read resource", func() error {
		respAssignments, err := r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			Assignments().
			Get(ctx, nil)
		if err != nil {
			return err
		}
		MapRemoteAssignmentStateToTerraform(ctx, &object, respAssignments)
		return nil
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for Settings Catalog resources.
//
//   - Retrieves the planned changes from the update request
//   - Constructs the resource request body from the plan
//   - Sends PUT request to update the base resource and settings
//   - Constructs the assignment request body from the plan
//   - Sends POST request to update the assignments
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API
//   - Updates the final state with the fresh data from the API
//
// The function ensures that both the settings and assignments are updated atomically,
// and the final state reflects the actual state of the resource on the server.
func (r *SettingsCatalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object SettingsCatalogProfileResourceModel

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

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	putRequest := graphcustom.PutRequestConfig{
		APIVersion:  graphcustom.GraphAPIBeta,
		Endpoint:    r.ResourcePath,
		ResourceID:  object.ID.ValueString(),
		RequestBody: requestBody,
	}

	// Use retryableOperation for main resource update
	err = retry.RetryableIntuneOperation(ctx, "update resource", retry.IntuneWrite, func() error {
		return graphcustom.PutRequestByResourceId(ctx, r.client.GetAdapter(), putRequest)
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.ReadPermissions)
		return
	}

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for update method",
			fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	// Use retryableAssignmentOperation for assignment update
	err = retry.RetryableAssignmentOperation(ctx, "update assignment", func() error {
		_, err := r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)
		return err
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for Settings Catalog resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *SettingsCatalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object SettingsCatalogProfileResourceModel

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

	err := retry.RetryableIntuneOperation(ctx, "delete resource", retry.IntuneWrite, func() error {
		return r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			Delete(ctx, nil)
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.ReadPermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
