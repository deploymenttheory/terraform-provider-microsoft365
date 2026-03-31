package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// Create handles the Create operation for Windows Updates autopatch updatable asset group resources.
//
// Operation: Creates a new updatable asset group then adds any configured device members.
// API Calls:
//   - POST /admin/windows/updates/updatableAssets
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.addMembersById (if entra_device_object_ids is set)
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//
// References:
//   - https://learn.microsoft.com/en-us/graph/api/adminwindowsupdates-post-updatableassets-updatableassetgroup?view=graph-rest-beta
//   - https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableassetgroup-addmembersbyid?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatableAssetGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel

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

	if err := validateRequest(ctx, r.client, &object); err != nil {
		resp.Diagnostics.AddError(
			"Validation Failed",
			fmt.Sprintf("Pre-flight validation failed for %s: %s", ResourceName, err.Error()),
		)
		return
	}

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatableAssets().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = convert.GraphToFrameworkString(createdResource.GetId())

	deviceObjectIDs := extractDeviceObjectIDs(&object)
	if len(deviceObjectIDs) > 0 {
		addRequest, err := constructAddMembersRequest(ctx, deviceObjectIDs)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing add members request",
				fmt.Sprintf("Could not construct add members request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			ByUpdatableAssetId(object.ID.ValueString()).
			MicrosoftGraphWindowsUpdatesAddMembersById().
			Post(ctx, addRequest, nil)

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
	opts.ConsistencyPredicate = autopatchUpdatableAssetGroupConsistencyPredicate(&object)

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

// Read handles the Read operation for Windows Updates autopatch updatable asset group resources.
//
// Operation: Retrieves the group by ID and its current member list.
// API Calls:
//   - GET /admin/windows/updates/updatableAssets/{id}
//   - GET /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//
// References:
//   - https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableasset-get?view=graph-rest-beta
//   - https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableassetgroup-list-members?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatableAssetGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read of resource: %s", ResourceName))

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

	groupID := object.ID.ValueString()

	remoteResource, err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatableAssets().
		ByUpdatableAssetId(groupID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	// Read members via the type-cast members collection endpoint.
	// GET with $expand=members on the base-type endpoint returns 400, and GET on the
	// type-cast item path returns 405. The members collection sub-path is the correct
	// approach per the Graph API documentation.
	membersURL := fmt.Sprintf(
		"https://graph.microsoft.com/beta/admin/windows/updates/updatableAssets/%s/microsoft.graph.windowsUpdates.updatableAssetGroup/members",
		groupID,
	)

	membersResponse, err := r.client.Admin().Windows().Updates().UpdatableAssets().
		WithUrl(membersURL).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodelswindowsupdates.UpdatableAssetable](
		membersResponse,
		r.client.GetAdapter(),
		graphmodelswindowsupdates.CreateUpdatableAssetCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating page iterator for members",
			fmt.Sprintf("Could not create page iterator: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, remoteResource, pageIterator)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Updates autopatch updatable asset group resources.
//
// Operation: Applies a diff-based membership update — adds new device members and removes
// device members that are no longer in the plan.
// API Calls:
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.addMembersById    (if members to add)
//   - POST /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.removeMembersById (if members to remove)
//   - GET  /admin/windows/updates/updatableAssets/{id}
//   - GET  /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//
// References:
//   - https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableassetgroup-addmembersbyid?view=graph-rest-beta
//   - https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableassetgroup-removemembersbyid?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatableAssetGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel
	var state WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if err := validateRequest(ctx, r.client, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Validation Failed",
			fmt.Sprintf("Pre-flight validation failed for %s: %s", ResourceName, err.Error()),
		)
		return
	}

	planIDs := make(map[string]bool)
	if !plan.EntraDeviceObjectIds.IsNull() && !plan.EntraDeviceObjectIds.IsUnknown() {
		for _, elem := range plan.EntraDeviceObjectIds.Elements() {
			if strVal, ok := elem.(types.String); ok {
				planIDs[strVal.ValueString()] = true
			}
		}
	}

	stateIDs := make(map[string]bool)
	if !state.EntraDeviceObjectIds.IsNull() && !state.EntraDeviceObjectIds.IsUnknown() {
		for _, elem := range state.EntraDeviceObjectIds.Elements() {
			if strVal, ok := elem.(types.String); ok {
				stateIDs[strVal.ValueString()] = true
			}
		}
	}

	var toAdd []string
	for id := range planIDs {
		if !stateIDs[id] {
			toAdd = append(toAdd, id)
		}
	}

	var toRemove []string
	for id := range stateIDs {
		if !planIDs[id] {
			toRemove = append(toRemove, id)
		}
	}

	groupID := state.ID.ValueString()

	if len(toAdd) > 0 {
		addRequest, err := constructAddMembersRequest(ctx, toAdd)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing add members request",
				fmt.Sprintf("Could not construct add members request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			ByUpdatableAssetId(groupID).
			MicrosoftGraphWindowsUpdatesAddMembersById().
			Post(ctx, addRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	if len(toRemove) > 0 {
		removeRequest, err := constructRemoveMembersRequest(ctx, toRemove)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing remove members request",
				fmt.Sprintf("Could not construct remove members request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			ByUpdatableAssetId(groupID).
			MicrosoftGraphWindowsUpdatesRemoveMembersById().
			Post(ctx, removeRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	plan.ID = state.ID

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = autopatchUpdatableAssetGroupConsistencyPredicate(&plan)

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Windows Updates autopatch updatable asset group resources.
//
// Operation: Permanently deletes the updatable asset group. The API removes all memberships
// as part of the group deletion.
// API Calls:
//   - DELETE /admin/windows/updates/updatableAssets/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableasset-delete?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatableAssetGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel

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
		Admin().
		Windows().
		Updates().
		UpdatableAssets().
		ByUpdatableAssetId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
