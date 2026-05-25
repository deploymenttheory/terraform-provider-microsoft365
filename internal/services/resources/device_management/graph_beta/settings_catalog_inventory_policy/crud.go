package graphBetaSettingsCatalogInventoryPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	configPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (r *InventoryPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object InventoryPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if object.ConfigurationPolicy != nil && len(object.ConfigurationPolicy.Settings) > 0 {
		configContent := fmt.Sprintf("%+v", object.ConfigurationPolicy)
		configPolicy.ResolveDCV2ConfigurationDepth(object.Name.ValueString(), configContent)
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

	config := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    "deviceManagement/inventoryPolicies",
		RequestBody: requestBody,
	}

	createdResource, err := customrequests.PostRequest(
		ctx,
		r.client.GetAdapter(),
		config,
		models.CreateDeviceManagementConfigurationPolicyFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	baseResource, ok := createdResource.(models.DeviceManagementConfigurationPolicyable)
	if !ok || baseResource.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error reading created resource ID",
			fmt.Sprintf("Could not extract ID from created resource: %s", ResourceName),
		)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, object.ID.ValueString()))

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	assignConfig := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    fmt.Sprintf("deviceManagement/inventoryPolicies('%s')/assign", object.ID.ValueString()),
		RequestBody: requestAssignment,
	}

	err = customrequests.PostRequestNoContent(ctx, r.client.GetAdapter(), assignConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
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

func (r *InventoryPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object InventoryPolicyResourceModel
	var identity sharedmodels.ResourceIdentity

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := constants.TfOperationRead
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

	currentState := object

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	getConfig := customrequests.GetRequestConfig{
		APIVersion: customrequests.GraphAPIBeta,
		Endpoint:   fmt.Sprintf("deviceManagement/inventoryPolicies('%s')", object.ID.ValueString()),
	}

	baseResult, err := customrequests.GetRequest(
		ctx,
		r.client.GetAdapter(),
		getConfig,
		models.CreateDeviceManagementConfigurationPolicyFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	baseResource, ok := baseResult.(models.DeviceManagementConfigurationPolicyable)
	if !ok {
		resp.Diagnostics.AddError(
			"Error reading resource",
			fmt.Sprintf("Could not cast response to DeviceManagementConfigurationPolicyable: %s", ResourceName),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, baseResource)

	settingsConfig := customrequests.GetRequestConfig{
		APIVersion: customrequests.GraphAPIBeta,
		Endpoint:   fmt.Sprintf("deviceManagement/inventoryPolicies('%s')/settings", object.ID.ValueString()),
	}

	settingsResult, err := customrequests.GetRequest(
		ctx,
		r.client.GetAdapter(),
		settingsConfig,
		models.CreateDeviceManagementConfigurationSettingCollectionResponseFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if settingsResult != nil {
		settingsResponse, ok := settingsResult.(models.DeviceManagementConfigurationSettingCollectionResponseable)
		if ok {
			err = StateInventoryPolicySettings(ctx, &object, settingsResponse, &currentState)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error mapping settings state",
					fmt.Sprintf("Could not map settings to Terraform state: %s", err.Error()),
				)
				return
			}
		}
	}

	assignmentsConfig := customrequests.GetRequestConfig{
		APIVersion: customrequests.GraphAPIBeta,
		Endpoint:   fmt.Sprintf("deviceManagement/inventoryPolicies('%s')/assignments", object.ID.ValueString()),
	}

	assignmentsResult, err := customrequests.GetRequest(
		ctx,
		r.client.GetAdapter(),
		assignmentsConfig,
		models.CreateDeviceManagementConfigurationPolicyAssignmentCollectionResponseFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if assignmentsResult != nil {
		if assignmentsResponse, ok := assignmentsResult.(models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable); ok && assignmentsResponse != nil {
			MapAssignmentsToTerraform(ctx, &object, assignmentsResponse.GetValue())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

func (r *InventoryPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan InventoryPolicyResourceModel
	var state InventoryPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.ConfigurationPolicy != nil && len(plan.ConfigurationPolicy.Settings) > 0 {
		configContent := fmt.Sprintf("%+v", plan.ConfigurationPolicy)
		configPolicy.ResolveDCV2ConfigurationDepth(plan.Name.ValueString(), configContent)
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

	putRequest := customrequests.PutRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    r.ResourcePath,
		ResourceID:  state.ID.ValueString(),
		RequestBody: requestBody,
	}

	err = customrequests.PutRequestByResourceId(ctx, r.client.GetAdapter(), putRequest)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Update Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	assignConfig := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    fmt.Sprintf("deviceManagement/inventoryPolicies('%s')/assign", state.ID.ValueString()),
		RequestBody: requestAssignment,
	}

	err = customrequests.PostRequestNoContent(ctx, r.client.GetAdapter(), assignConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

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

func (r *InventoryPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object InventoryPolicyResourceModel

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

	deleteConfig := customrequests.DeleteRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "deviceManagement/inventoryPolicies",
		ResourceIDPattern: "('id')",
		ResourceID:        object.ID.ValueString(),
	}

	err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
