package graphBetaEndpointPrivilegeManagement

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	construct "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequest "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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
func (r *EndpointPrivilegeManagementResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object sharedmodels.SettingsCatalogProfileResourceModel

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

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	if object.Assignments != nil {
		requestAssignment, err := construct.ConstructConfigurationPolicyAssignment(ctx, object.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for Create Method",
				fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		_, err = r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
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
func (r *EndpointPrivilegeManagementResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object sharedmodels.SettingsCatalogProfileResourceModel
	var baseResource models.DeviceManagementConfigurationPolicyable
	var assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

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

	baseResource, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, baseResource)

	settingsConfig := customrequest.GetRequestConfig{
		APIVersion:        customrequest.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		EndpointSuffix:    "/settings",
		ResourceIDPattern: "('id')",
		ResourceID:        object.ID.ValueString(),
		QueryParameters: map[string]string{
			"$expand": "children",
		},
	}

	var settingsResponse []byte

	settingsResponse, err = customrequest.GetRequestByResourceId(
		ctx,
		r.client.GetAdapter(),
		settingsConfig,
	)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	sharedstater.StateConfigurationPolicySettings(ctx, &object, settingsResponse)

	assignmentsResponse, err = r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	sharedstater.StateConfigurationPolicyAssignment(ctx, object.Assignments, assignmentsResponse)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
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
func (r *EndpointPrivilegeManagementResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sharedmodels.SettingsCatalogProfileResourceModel
	var state sharedmodels.SettingsCatalogProfileResourceModel

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

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	putRequest := customrequest.PutRequestConfig{
		APIVersion:  customrequest.GraphAPIBeta,
		Endpoint:    r.ResourcePath,
		ResourceID:  state.ID.ValueString(),
		RequestBody: requestBody,
	}

	err = customrequest.PutRequestByResourceId(
		ctx,
		r.client.GetAdapter(),
		putRequest)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.ReadPermissions)
		return
	}

	if plan.Assignments != nil {
		requestAssignment, err := construct.ConstructConfigurationPolicyAssignment(ctx, plan.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for Update Method",
				fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		_, err = r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(state.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Settings Catalog resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *EndpointPrivilegeManagementResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object sharedmodels.SettingsCatalogProfileResourceModel

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

	err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
