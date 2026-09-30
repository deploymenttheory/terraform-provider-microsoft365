package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for device configuration template JSON resources.
//
//   - Retrieves the planned configuration and validates the create timeout
//   - Constructs the profile metadata, JSON settings, and assignment request bodies
//   - Sends POST request to create the profile and captures its ID
//   - Saves the initial state before writing relationships and assignments
//   - Reconciles certificate relationships and sends assignments if specified
//   - Calls Read with retry to populate the final state from the API
//
// Saving the profile ID before subsequent writes allows Terraform to manage or
// destroy the created profile if a relationship or assignment operation fails.
func (r *DeviceConfigurationTemplatesJsonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeviceConfigurationTemplatesJsonResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	request, err := constructResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid device configuration profile", err.Error())
		return
	}

	assignments, err := constructAssignment(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid profile assignments", err.Error())
		return
	}

	// iOS wired profiles reject embedded certificate bindings. Establish these
	// through the typed $ref endpoints after saving the created profile ID.
	bindings := request.GetAdditionalData()
	if *request.GetOdataType() == "#microsoft.graph.iosWiredNetworkConfiguration" {
		settings := maps.Clone(bindings)
		for name := range deviceConfigurationRelationships[*request.GetOdataType()] {
			delete(settings, name+"@odata.bind")
		}
		request.SetAdditionalData(settings)
	}

	remote, err := r.client.
		DeviceManagement().
		DeviceConfigurations().
		Post(ctx, request, nil)

	request.SetAdditionalData(bindings)
	if err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" {
		resp.Diagnostics.AddError(
			"Invalid create response",
			"Graph did not return the created profile ID.",
		)
		return
	}

	data.ID = types.StringPointerValue(remote.GetId())
	if data.Description.IsUnknown() {
		data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	}

	if data.RoleScopeTagIds.IsUnknown() {
		data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remote.GetRoleScopeTagIds())
	}

	// Preserve the remote ID before the second API write so an assignment failure cannot orphan the profile.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.createRelationships(ctx, data.ID.ValueString(), request); err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if len(assignments.GetAssignments()) > 0 {
		_, err = r.client.
			DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(data.ID.ValueString()).
			Assign().
			Post(ctx, assignments, nil)

		if err != nil {
			kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	readReq := resource.ReadRequest{State: resp.State}
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

// Read handles the Read operation for device configuration template JSON resources.
//
//   - Retrieves the current state and validates the read timeout
//   - Gets the complete profile JSON and maps its metadata to Terraform state
//   - Reads certificate relationships and recovers encrypted OMA values
//   - Maps complete writable settings while preserving equivalent configured JSON
//   - Gets all assignment pages and maps the targets to Terraform state
//   - Saves the refreshed state
//
// A missing base profile removes the resource from state. Errors reading settings
// or assignments preserve the profile so a failed child request cannot orphan it.
func (r *DeviceConfigurationTemplatesJsonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object DeviceConfigurationTemplatesJsonResourceModel
	var identity sharedmodels.ResourceIdentity
	var respResource json.RawMessage

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

	err := customrequests.JSONRequest(ctx, r.client.GetAdapter(), abstractions.GET,
		r.ResourcePath+"/{id}", map[string]string{"id": object.ID.ValueString()}, nil, &respResource)
	if err != nil {
		kiotaerrors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions,
			kiotaerrors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}

	if err := MapRemoteResourceStateToTerraform(ctx, &object, respResource); err != nil {
		resp.Diagnostics.AddError("Cannot read profile metadata", err.Error())
		return
	}

	if err := r.MapRemoteSettingsStateToTerraform(ctx, &object, respResource); err != nil {
		kiotaerrors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions,
			kiotaerrors.GraphErrorOptions{
				PreserveStateOnReadBadRequest: true,
				PreserveStateOnReadNotFound:   true,
			},
		)
		return
	}

	var assignments []graphmodels.DeviceConfigurationAssignmentable

	builder := r.client.
		DeviceManagement().
		DeviceConfigurations().
		ByDeviceConfigurationId(object.ID.ValueString()).
		Assignments()
	for {
		page, err := builder.Get(ctx, nil)
		if err != nil {
			kiotaerrors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions,
				kiotaerrors.GraphErrorOptions{
					PreserveStateOnReadBadRequest: true,
					PreserveStateOnReadNotFound:   true,
				},
			)
			return
		}

		if page == nil {
			resp.Diagnostics.AddError(
				"Invalid assignment response",
				"Graph returned an empty assignment response.",
			)
			return
		}

		assignments = append(assignments, page.GetValue()...)
		next := page.GetOdataNextLink()
		if next == nil || *next == "" {
			break
		}

		builder = builder.WithUrl(*next)
	}

	resp.Diagnostics.Append(MapAssignmentsToTerraform(ctx, &object, assignments)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for device configuration template JSON resources.
//
//   - Retrieves the plan and prior state and validates the update timeout
//   - Retains the profile ID and constructs the desired and previous request bodies
//   - Rejects changes to the root profile type before making API writes
//   - Sends PATCH request when profile metadata or ordinary settings change
//   - Reconciles certificate relationships through their separate $ref endpoints
//   - Sends the assignment request when targets change, including an empty set
//   - Saves the planned state and calls Read with retry to refresh it from the API
//
// Profile settings, relationships, and assignments use separate API operations.
// Assignment-only changes retain the existing profile without rewriting its settings.
func (r *DeviceConfigurationTemplatesJsonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceConfigurationTemplatesJsonResourceModel
	var state DeviceConfigurationTemplatesJsonResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	plan.ID = state.ID
	if plan.Description.IsUnknown() {
		plan.Description = state.Description
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

	previousRequestBody, err := constructResource(ctx, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing previous resource for update method",
			fmt.Sprintf("Could not construct previous resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	if *requestBody.GetOdataType() != *previousRequestBody.GetOdataType() {
		resp.Diagnostics.AddError(
			"Profile type cannot be changed",
			"Recreate the profile with Terraform's -replace option to change the root @odata.type.",
		)
		return
	}

	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for update method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Relationship bindings are updated through $ref: PATCH cannot remove them.
	settings := maps.Clone(requestBody.GetAdditionalData())
	previousSettings := maps.Clone(previousRequestBody.GetAdditionalData())
	for name := range deviceConfigurationRelationships[*requestBody.GetOdataType()] {
		delete(settings, name+"@odata.bind")
		delete(previousSettings, name+"@odata.bind")
	}

	if !reflect.DeepEqual(settings, previousSettings) ||
		!plan.DisplayName.Equal(state.DisplayName) ||
		!plan.Description.Equal(state.Description) ||
		!plan.RoleScopeTagIds.Equal(state.RoleScopeTagIds) {
		bindings := requestBody.GetAdditionalData()
		requestBody.SetAdditionalData(settings)
		_, err = r.client.
			DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(state.ID.ValueString()).
			Patch(ctx, requestBody, nil)

		requestBody.SetAdditionalData(bindings)
		if err != nil {
			kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	err = r.updateRelationships(ctx, state.ID.ValueString(), previousRequestBody, requestBody)
	if err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	if !plan.Assignments.Equal(state.Assignments) {
		_, err = r.client.
			DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(state.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
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

// Delete handles the Delete operation for device configuration template JSON resources.
//
//   - Retrieves the current state and validates the delete timeout
//   - Sends DELETE request to remove the profile from the API
//   - Removes the resource from Terraform state after successful deletion
//
// Graph removes the profile's settings and assignments with the profile. Referenced
// certificate profiles and assignment groups remain independently managed resources.
func (r *DeviceConfigurationTemplatesJsonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object DeviceConfigurationTemplatesJsonResourceModel

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
		DeviceConfigurations().
		ByDeviceConfigurationId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
