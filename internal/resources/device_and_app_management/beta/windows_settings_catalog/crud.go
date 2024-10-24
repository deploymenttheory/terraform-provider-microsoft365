package graphBetaWindowsSettingsCatalog

import (
	"context"
	"fmt"
	"time"

	betaClient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client/beta"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the Create operation.
func (r *WindowsSettingsCatalogResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WindowsSettingsCatalogProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	resource, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		Post(context.Background(), requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for create method",
			fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(plan.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(*resource.GetId())

	// Set the state after creating the resource
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the Read method to refresh the state after creation
	r.Read(ctx, resource.ReadRequest{
		State: resp.State,
	}, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))

}

// Read handles the Read operation.
func (r *WindowsSettingsCatalogResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WindowsSettingsCatalogProfileResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Get base resource
	respResource, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(state.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &state, respResource)

	// Get settings
	respSettings, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(state.ID.ValueString()).
		Settings().
		Get(context.Background(), &msgraphsdk.ConfigurationPoliciesItemSettingsRequestBuilderGetRequestConfiguration{
			QueryParameters: &msgraphsdk.ConfigurationPoliciesItemSettingsRequestBuilderGetQueryParameters{
				Expand: []string{""}, // Expand all related settings
			},
		})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteSettingsStateToTerraform(ctx, &state, respSettings)

	// Get assignments
	respAssignments, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(state.ID.ValueString()).
		Assignments().
		Get(context.Background(), nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteAssignmentStateToTerraform(ctx, &state, respAssignments)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *WindowsSettingsCatalogResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsSettingsCatalogProfileResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	putRequest := betaClient.CustomPutRequestConfig{
		APIVersion:  betaClient.GraphAPIBeta,
		Endpoint:    "deviceManagement/configurationPolicies",
		ResourceID:  plan.ID.ValueString(),
		RequestBody: requestBody,
	}

	err = betaClient.SendCustomPutRequestByResourceId(ctx, r.client, putRequest)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// Step 2: Update assignments
	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for update method",
			fmt.Sprintf("Could not construct assignment: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(plan.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Step 3: Use Read to refresh the state
	readResp := resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &readResp)

	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))

}

// Delete handles the Delete operation.
func (r *WindowsSettingsCatalogResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WindowsSettingsCatalogProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(data.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.ReadPermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
