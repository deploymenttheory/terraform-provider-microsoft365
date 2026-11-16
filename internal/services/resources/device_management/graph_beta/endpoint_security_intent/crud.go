package graphBetaEndpointSecurityIntent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for endpoint security intent resources.
//
// Operation: Creates a new intent from a security template and assigns it
// API Calls:
//   - POST /deviceManagement/templates/{templateId}/createInstance
//   - POST /deviceManagement/intents/{id}/assign (only when assignments are configured)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-deviceintent-devicemanagementtemplate-createinstance?view=graph-rest-beta
func (r *EndpointSecurityIntentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object EndpointSecurityIntentResourceModel

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

	assignments, err := constructAssignments(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignments",
			fmt.Sprintf("Could not construct assignments for resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		DeviceManagement().
		Templates().
		ByDeviceManagementTemplateId(object.TemplateId.ValueString()).
		CreateInstance().
		Post(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if createdResource == nil || createdResource.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error extracting resource ID",
			fmt.Sprintf("Could not extract ID from created resource: %s. The API may not return the full resource on creation.", ResourceName),
		)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(assignments.GetAssignments()) > 0 {
		err = r.client.
			DeviceManagement().
			Intents().
			ByDeviceManagementIntentId(object.ID.ValueString()).
			Assign().
			Post(ctx, assignments, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
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

// Read handles the Read operation for endpoint security intent resources.
//
// Operation: Retrieves an intent with its settings and assignments
// API Calls:
//   - GET /deviceManagement/intents/{id}
//   - GET /deviceManagement/intents/{id}/settings
//   - GET /deviceManagement/intents/{id}/assignments
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-deviceintent-devicemanagementintent-get?view=graph-rest-beta
func (r *EndpointSecurityIntentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object EndpointSecurityIntentResourceModel
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

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	intentRequest := r.client.
		DeviceManagement().
		Intents().
		ByDeviceManagementIntentId(object.ID.ValueString())

	remoteResource, err := intentRequest.Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	remoteSettings, err := intentRequest.Settings().Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	remoteAssignments, err := intentRequest.Assignments().Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	object = MapRemoteResourceStateToTerraform(ctx, object, remoteResource, remoteSettings.GetValue(), remoteAssignments.GetValue())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for endpoint security intent resources.
//
// Operation: Updates the intent properties, its settings and its assignments
// API Calls:
//   - PATCH /deviceManagement/intents/{id}
//   - POST /deviceManagement/intents/{id}/updateSettings (only when the configured settings changed)
//   - POST /deviceManagement/intents/{id}/assign (only when the assignments changed)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-deviceintent-devicemanagementintent-updatesettings?view=graph-rest-beta
func (r *EndpointSecurityIntentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EndpointSecurityIntentResourceModel
	var state EndpointSecurityIntentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

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

	requestBody, err := constructUpdateResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	intentRequest := r.client.
		DeviceManagement().
		Intents().
		ByDeviceManagementIntentId(state.ID.ValueString())

	_, err = intentRequest.Patch(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// Removed settings are not reset: Graph has no way to restore a template default, so they keep their current value.
	if !plan.Settings.IsNull() && !plan.Settings.Equal(state.Settings) {
		settingsBody, err := constructUpdateSettings(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing settings",
				fmt.Sprintf("Could not construct settings for resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = intentRequest.UpdateSettings().Post(ctx, settingsBody, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	if !plan.Assignments.Equal(state.Assignments) {
		assignments, err := constructAssignments(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Could not construct assignments for resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = intentRequest.Assign().Post(ctx, assignments, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	// resp.State still holds the prior state. Read maps only the settings listed in state, and maps no assignments to null only when
	// state has none, so it needs the planned lists to know what to track. Every other attribute is read back from Graph.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("settings"), plan.Settings)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("assignments"), plan.Assignments)...)
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

// Delete handles the Delete operation for endpoint security intent resources.
//
// Operation: Deletes an intent
// API Calls:
//   - DELETE /deviceManagement/intents/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-deviceintent-devicemanagementintent-delete?view=graph-rest-beta
func (r *EndpointSecurityIntentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object EndpointSecurityIntentResourceModel

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

	intentId := object.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting %s with ID: %s", ResourceName, intentId))

	err := r.client.
		DeviceManagement().
		Intents().
		ByDeviceManagementIntentId(intentId).
		Delete(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
