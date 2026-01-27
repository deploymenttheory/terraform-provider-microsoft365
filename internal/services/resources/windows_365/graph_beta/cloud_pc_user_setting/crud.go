package graphBetaCloudPcUserSetting

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the Create operation for Cloud PC user setting resources.
//
// Operation: Creates a new Cloud PC user setting with optional assignments
// API Calls:
//   - POST /deviceManagement/virtualEndpoint/userSettings
//   - POST /deviceManagement/virtualEndpoint/userSettings/{id}/assign (if assignments are configured)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/virtualendpoint-post-usersettings?view=graph-rest-beta
func (r *CloudPcUserSettingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object CloudPcUserSettingResourceModel

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
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	cloudPcUserSetting, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		UserSettings().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*cloudPcUserSetting.GetId())

	// Handle assignments if present
	if !object.Assignments.IsNull() && !object.Assignments.IsUnknown() {
		assignBody, err := constructAssignmentsRequestBody(ctx, object.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments request body",
				fmt.Sprintf("Could not construct assignments request body: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			VirtualEndpoint().
			UserSettings().
			ByCloudPcUserSettingId(object.ID.ValueString()).
			Assign().
			Post(ctx, assignBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
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

// Read handles the Read operation for Cloud PC user setting resources.
//
// Operation: Retrieves a Cloud PC user setting including assignments
// API Calls:
//   - GET /deviceManagement/virtualEndpoint/userSettings/{id}?$expand=assignments
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcusersetting-get?view=graph-rest-beta
func (r *CloudPcUserSettingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object CloudPcUserSettingResourceModel
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

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	cloudPcUserSetting, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		UserSettings().
		ByCloudPcUserSettingId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.VirtualEndpointUserSettingsCloudPcUserSettingItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.VirtualEndpointUserSettingsCloudPcUserSettingItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
				Select: []string{"*"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, cloudPcUserSetting)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Cloud PC user setting resources.
//
// Operation: Updates a Cloud PC user setting with optional assignment updates
// API Calls:
//   - PATCH /deviceManagement/virtualEndpoint/userSettings/{id}
//   - POST /deviceManagement/virtualEndpoint/userSettings/{id}/assign (if assignments are configured)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcusersetting-update?view=graph-rest-beta
func (r *CloudPcUserSettingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudPcUserSettingResourceModel
	var state CloudPcUserSettingResourceModel

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
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	updated, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		UserSettings().
		ByCloudPcUserSettingId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, updated)

	// Handle assignments if present
	if !plan.Assignments.IsNull() && !plan.Assignments.IsUnknown() {
		assignBody, err := constructAssignmentsRequestBody(ctx, plan.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments request body",
				fmt.Sprintf("Could not construct assignments request body: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			VirtualEndpoint().
			UserSettings().
			ByCloudPcUserSettingId(state.ID.ValueString()).
			Assign().
			Post(ctx, assignBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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

// Delete handles the Delete operation for Cloud PC user setting resources.
//
// Operation: Deletes a Cloud PC user setting
// API Calls:
//   - DELETE /deviceManagement/virtualEndpoint/userSettings/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcusersetting-delete?view=graph-rest-beta
func (r *CloudPcUserSettingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object CloudPcUserSettingResourceModel

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
		VirtualEndpoint().
		UserSettings().
		ByCloudPcUserSettingId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
