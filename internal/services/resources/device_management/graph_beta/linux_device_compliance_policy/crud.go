package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequest "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// Create handles the Create operation for Settings Catalog resources.
//
//   - Retrieves the planned configuration from the create request
//   - Constructs the resource request body from the plan
//   - Sends POST request to create the base resource and settings
//   - Captures the new resource ID from the response
//   - Constructs and sends assignment configuration if specified with retry
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API with retry
//   - Updates the final state with the fresh data from the API
//
// The function ensures that both the settings catalog profile and its assignments
// (if specified) are created properly. The settings must be defined during creation
// as they are required for a successful deployment, while assignments are optional.
func (r *LinuxDeviceCompliancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object LinuxDeviceCompliancePolicyResourceModel

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

	baseResource, err := r.client.
		DeviceManagement().
		CompliancePolicies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	// Handle scheduleActionsForRules if present
	if !object.ScheduledActions.IsNull() && !object.ScheduledActions.IsUnknown() {
		var scheduledActionsModels []ScheduledActionForRuleModel
		diags := object.ScheduledActions.ElementsAs(ctx, &scheduledActionsModels, false)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error parsing scheduled actions for rules",
				fmt.Sprintf("Could not parse scheduled actions list: %v", diags.Errors()),
			)
			return
		}

		for _, scheduledAction := range scheduledActionsModels {
			scheduleRequestBody, err := constructScheduledActions(ctx, scheduledAction)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing scheduled actions for rules",
					fmt.Sprintf("Could not construct scheduled actions request: %s", err.Error()),
				)
				return
			}

			_, err = r.client.
				DeviceManagement().
				CompliancePolicies().
				ByDeviceManagementCompliancePolicyId(object.ID.ValueString()).
				SetScheduledActions().
				Post(ctx, scheduleRequestBody, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully scheduled actions for rule '%s' for policy ID: %s",
				scheduledAction.RuleName.ValueString(), object.ID.ValueString()))
		}
	}

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		CompliancePolicies().
		ByDeviceManagementCompliancePolicyId(object.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
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
//   - Gets the settings configuration from the API using @odata.nextLink
//   - Maps the settings configuration to Terraform state
//   - Gets the assignments configuration from the API
//   - Maps the assignments configuration to Terraform state
//
// The function ensures that all components (base resource, settings, and assignments)
// are properly read and mapped into the Terraform state, providing a complete view
// of the resource's current configuration on the server.
func (r *LinuxDeviceCompliancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object LinuxDeviceCompliancePolicyResourceModel
	var respResource models.DeviceManagementCompliancePolicyable

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

	// Create a copy of the current state to use as "plan" for secret value preservation
	currentState := object

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	respResource, err := r.client.
		DeviceManagement().
		CompliancePolicies().
		ByDeviceManagementCompliancePolicyId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.CompliancePoliciesDeviceManagementCompliancePolicyItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.CompliancePoliciesDeviceManagementCompliancePolicyItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments", "scheduledActionsForRule"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, respResource)

	// Use PageIterator from graph core as default response returns only the first 25 settings.
	tflog.Debug(ctx, "Using Microsoft Graph SDK PageIterator for settings")

	allSettings, err := r.getAllPolicySettingsWithPageIterator(ctx, object.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	combinedSettingsResponse := models.NewDeviceManagementConfigurationSettingCollectionResponse()
	combinedSettingsResponse.SetValue(allSettings)

	err = StateConfigurationPolicySettings(ctx, &object, combinedSettingsResponse, &currentState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error mapping settings state",
			fmt.Sprintf("Could not map settings to Terraform state: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// getAllPolicySettingsWithPageIterator retrieves all settings for a given policy using a page iterator util
// from the SDK.
func (r *LinuxDeviceCompliancePolicyResource) getAllPolicySettingsWithPageIterator(ctx context.Context, policyId string) ([]models.DeviceManagementConfigurationSettingable, error) {
	var allSettings []models.DeviceManagementConfigurationSettingable

	settingsResponse, err := r.client.
		DeviceManagement().
		CompliancePolicies().
		ByDeviceManagementCompliancePolicyId(policyId).
		Settings().
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[models.DeviceManagementConfigurationSettingable](
		settingsResponse,
		r.client.GetAdapter(),
		models.CreateDeviceManagementConfigurationSettingCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item models.DeviceManagementConfigurationSettingable) bool {
		if item != nil {
			allSettings = append(allSettings, item)

			if len(allSettings)%25 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d settings (estimated page %d)", len(allSettings), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total settings", len(allSettings)))

	return allSettings, nil
}

// Update handles the Update operation for Settings Catalog resources.
//
//   - Retrieves the planned changes from the update request
//   - Constructs the resource request body from the plan
//   - Sends PUT request to update the base resource and settings
//   - Constructs the assignment request body from the plan
//   - Sends POST request to update the assignments
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API with retry
//   - Updates the final state with the fresh data from the API
//
// The function ensures that both the settings and assignments are updated atomically,
// and the final state reflects the actual state of the resource on the server.
func (r *LinuxDeviceCompliancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LinuxDeviceCompliancePolicyResourceModel
	var state LinuxDeviceCompliancePolicyResourceModel

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
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.ReadPermissions)
		return
	}

	// Handle scheduleActionsForRules if present
	if !plan.ScheduledActions.IsNull() && !plan.ScheduledActions.IsUnknown() {
		var scheduledActionsModels []ScheduledActionForRuleModel
		diags := plan.ScheduledActions.ElementsAs(ctx, &scheduledActionsModels, false)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error parsing scheduled actions for rules",
				fmt.Sprintf("Could not parse scheduled actions list: %v", diags.Errors()),
			)
			return
		}

		for _, scheduledAction := range scheduledActionsModels {
			scheduleRequestBody, err := constructScheduledActions(ctx, scheduledAction)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing scheduled actions for rules",
					fmt.Sprintf("Could not construct scheduled actions request: %s", err.Error()),
				)
				return
			}

			_, err = r.client.
				DeviceManagement().
				CompliancePolicies().
				ByDeviceManagementCompliancePolicyId(plan.ID.ValueString()).
				SetScheduledActions().
				Post(ctx, scheduleRequestBody, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully scheduled actions for rule '%s' for policy ID: %s",
				scheduledAction.RuleName.ValueString(), plan.ID.ValueString()))
		}
	}

	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		CompliancePolicies().
		ByDeviceManagementCompliancePolicyId(state.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
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

// Delete handles the Delete operation for Settings Catalog resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *LinuxDeviceCompliancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object LinuxDeviceCompliancePolicyResourceModel

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
		CompliancePolicies().
		ByDeviceManagementCompliancePolicyId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
