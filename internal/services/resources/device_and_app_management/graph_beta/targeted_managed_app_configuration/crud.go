package graphBetaTargetedManagedAppConfigurations

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	deviceappmanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
)

// Create handles the Create operation for Targeted Managed App Configuration resources.
//
// Operation: Creates a new targeted managed app configuration policy
// API Calls:
//   - POST /deviceAppManagement/targetedManagedAppConfigurations
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-targetedmanagedappconfiguration-create?view=graph-rest-beta
func (r *TargetedManagedAppConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state TargetedManagedAppConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run request validation (caller-responder pattern in validate.go)
	if diags := validateRequest(ctx, &state); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceAppManagement().
		TargetedManagedAppConfigurations().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	state.ID = types.StringValue(*baseResource.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName

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

// Read handles the Read operation for Targeted Managed App Configuration resources.
//
// Operation: Retrieves a targeted managed app configuration policy by ID
// API Calls:
//   - GET /deviceAppManagement/targetedManagedAppConfigurations/{targetedManagedAppConfigurationId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-targetedmanagedappconfiguration-get?view=graph-rest-beta
func (r *TargetedManagedAppConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TargetedManagedAppConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := constants.TfOperationRead
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	result, err := r.client.
		DeviceAppManagement().
		TargetedManagedAppConfigurations().
		ByTargetedManagedAppConfigurationId(state.ID.ValueString()).
		Get(ctx, &deviceappmanagement.TargetedManagedAppConfigurationsTargetedManagedAppConfigurationItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.TargetedManagedAppConfigurationsTargetedManagedAppConfigurationItemRequestBuilderGetQueryParameters{
				Expand: []string{"apps", "assignments", "settings"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, result)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Targeted Managed App Configuration resources.
//
// Operation: Updates an existing targeted managed app configuration policy with multiple update paths
// API Calls:
//   - PATCH /deviceAppManagement/targetedManagedAppConfigurations/{targetedManagedAppConfigurationId}
//   - POST /deviceAppManagement/targetedManagedAppConfigurations/{targetedManagedAppConfigurationId}/targetApps
//   - POST /deviceAppManagement/targetedManagedAppConfigurations/{targetedManagedAppConfigurationId}/changeSettings
//   - POST /deviceAppManagement/targetedManagedAppConfigurations/{targetedManagedAppConfigurationId}/assign
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-targetedmanagedappconfiguration-update?view=graph-rest-beta
func (r *TargetedManagedAppConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TargetedManagedAppConfigurationResourceModel
	var state TargetedManagedAppConfigurationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Run request validation. Defined here because the crete and update req builders
	// are separate.
	if diags := validateRequest(ctx, &plan); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Step 1: PATCH base properties (display name, description, role scope tags)
	if !plan.DisplayName.Equal(state.DisplayName) ||
		!plan.Description.Equal(state.Description) ||
		!plan.RoleScopeTagIds.Equal(state.RoleScopeTagIds) {

		tflog.Debug(ctx, "Base properties changed, updating via PATCH")

		requestBody, err := constructBaseResourceUpdate(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing base resource update",
				fmt.Sprintf("Could not construct base resource update: %s", err.Error()),
			)
			return
		}

		_, err = r.client.
			DeviceAppManagement().
			TargetedManagedAppConfigurations().
			ByTargetedManagedAppConfigurationId(state.ID.ValueString()).
			Patch(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else {
		tflog.Debug(ctx, "Base properties unchanged, skipping base PATCH")
	}

	// Step 2: Update apps changes via /targetApps endpoint
	if !plan.Apps.Equal(state.Apps) || !plan.AppGroupType.Equal(state.AppGroupType) {
		tflog.Debug(ctx, "Apps or app group type changed, updating via /targetApps endpoint")

		appsRequest, err := constructTargetAppsUpdate(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing target apps update",
				fmt.Sprintf("Could not construct target apps update: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceAppManagement().
			TargetedManagedAppConfigurations().
			ByTargetedManagedAppConfigurationId(state.ID.ValueString()).
			TargetApps().
			Post(ctx, appsRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else {
		tflog.Debug(ctx, "Apps unchanged, skipping apps update")
	}

	// Step 3:Update custom settings changes via PATCH with @odata.type
	if !customSettingsEqual(plan.CustomSettings, state.CustomSettings) {
		tflog.Debug(ctx, "Custom settings changed, updating via PATCH with custom settings")

		requestBody, err := constructCustomSettingsUpdate(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing custom settings update",
				fmt.Sprintf("Could not construct custom settings update: %s", err.Error()),
			)
			return
		}

		_, err = r.client.
			DeviceAppManagement().
			TargetedManagedAppConfigurations().
			ByTargetedManagedAppConfigurationId(state.ID.ValueString()).
			Patch(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else {
		tflog.Debug(ctx, "Custom settings unchanged, skipping custom settings update")
	}

	// Step 4: Update settings catalog changes via /changeSettings endpoint
	if !settingsEqual(plan.SettingsCatalog, state.SettingsCatalog) {
		tflog.Debug(ctx, "Settings catalog changed, updating via /changeSettings endpoint")

		settingsRequest, err := constructSettingsUpdate(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing settings update",
				fmt.Sprintf("Could not construct settings update: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceAppManagement().
			TargetedManagedAppConfigurations().
			ByTargetedManagedAppConfigurationId(state.ID.ValueString()).
			ChangeSettings().
			Post(ctx, settingsRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else {
		tflog.Debug(ctx, "Settings unchanged, skipping settings update")
	}

	// Step 5: Update assignments with the /assign endpoint
	// Note: Req for update is the same as create
	if !plan.Assignments.Equal(state.Assignments) {
		tflog.Debug(ctx, "Assignments changed, updating via /assign endpoint")

		assignments, err := constructAssignments(ctx, plan.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Could not construct assignments: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		assignRequest := deviceappmanagement.NewTargetedManagedAppConfigurationsItemAssignPostRequestBody()
		assignRequest.SetAssignments(assignments)

		err = r.client.
			DeviceAppManagement().
			TargetedManagedAppConfigurations().
			ByTargetedManagedAppConfigurationId(state.ID.ValueString()).
			Assign().
			Post(ctx, assignRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else {
		tflog.Debug(ctx, "Assignments unchanged, skipping assignments update")
	}

	// Step 6: Read with retry to get complete updated state
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Targeted Managed App Configuration resources.
//
// Operation: Deletes a targeted managed app configuration policy
// API Calls:
//   - DELETE /deviceAppManagement/targetedManagedAppConfigurations/{targetedManagedAppConfigurationId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-targetedmanagedappconfiguration-delete?view=graph-rest-beta
func (r *TargetedManagedAppConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state TargetedManagedAppConfigurationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeviceAppManagement().TargetedManagedAppConfigurations().ByTargetedManagedAppConfigurationId(state.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// customSettingsEqual compares two custom settings slices for equality
func customSettingsEqual(a, b []KeyValuePairResourceModel) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}

	return true
}

// settingsEqual compares two settings catalog configurations for equality
func settingsEqual(a, b *DeviceConfigV2GraphServiceResourceModel) bool {
	// Both nil
	if a == nil && b == nil {
		return true
	}

	// One nil, one not
	if a == nil || b == nil {
		return false
	}

	// Different lengths
	if len(a.Settings) != len(b.Settings) {
		return false
	}

	// Deep comparison
	for i := range a.Settings {
		if !reflect.DeepEqual(a.Settings[i], b.Settings[i]) {
			return false
		}
	}

	return true
}
