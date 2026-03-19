package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment

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
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// Create handles the Create operation for Windows Updates autopatch updatable asset group assignment resources.
//
// Operation: Adds device members to an updatable asset group
// API Calls:
//   - POST /admin/windows/updates/updatableAssets/{updatableAssetId}/microsoft.graph.windowsUpdates.addMembersById
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableassetgroup-addmembersbyid?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel

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

	object.ID = object.UpdatableAssetGroupId

	deviceIDs := extractDeviceIDs(&object)

	if len(deviceIDs) > 0 {
		requestBody, err := constructAddMembersRequest(ctx, deviceIDs)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing resource",
				fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			ByUpdatableAssetId(object.UpdatableAssetGroupId.ValueString()).
			MicrosoftGraphWindowsUpdatesAddMembersById().
			Post(ctx, requestBody, nil)

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

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for Windows Updates autopatch updatable asset group assignment resources.
//
// Operation: Retrieves members of an updatable asset group
// API Calls:
//   - GET /admin/windows/updates/updatableAssets/{updatableAssetId}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableassetgroup-list-members?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s for group ID: %s", ResourceName, object.UpdatableAssetGroupId.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Read members via the dedicated type-cast members collection endpoint:
	// GET /admin/windows/updates/updatableAssets/{id}/microsoft.graph.windowsUpdates.updatableAssetGroup/members
	// $expand=members on the base-type endpoint returns 400, and GET on the type-cast
	// item path returns 405. The members collection sub-path is the correct approach
	// per the Graph API documentation.
	groupID := object.UpdatableAssetGroupId.ValueString()
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

	MapRemoteStateToTerraform(ctx, &object, pageIterator)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Updates autopatch updatable asset group assignment resources.
//
// Operation: Adds new members and removes existing members from an updatable asset group
// API Calls:
//   - POST /admin/windows/updates/updatableAssets/{updatableAssetId}/microsoft.graph.windowsUpdates.addMembersById
//   - POST /admin/windows/updates/updatableAssets/{updatableAssetId}/microsoft.graph.windowsUpdates.removeMembersById
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableassetgroup-addmembersbyid?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel
	var state WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s for group ID: %s", ResourceName, state.UpdatableAssetGroupId.ValueString()))

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

	planDeviceIDs := make(map[string]bool)
	if !plan.EntraDeviceIds.IsNull() && !plan.EntraDeviceIds.IsUnknown() {
		for _, elem := range plan.EntraDeviceIds.Elements() {
			if strVal, ok := elem.(types.String); ok {
				planDeviceIDs[strVal.ValueString()] = true
			}
		}
	}

	stateDeviceIDs := make(map[string]bool)
	if !state.EntraDeviceIds.IsNull() && !state.EntraDeviceIds.IsUnknown() {
		for _, elem := range state.EntraDeviceIds.Elements() {
			if strVal, ok := elem.(types.String); ok {
				stateDeviceIDs[strVal.ValueString()] = true
			}
		}
	}

	var devicesToAdd []string
	for id := range planDeviceIDs {
		if !stateDeviceIDs[id] {
			devicesToAdd = append(devicesToAdd, id)
		}
	}

	var devicesToRemove []string
	for id := range stateDeviceIDs {
		if !planDeviceIDs[id] {
			devicesToRemove = append(devicesToRemove, id)
		}
	}

	if len(devicesToAdd) > 0 {
		addRequest, err := constructAddMembersRequest(ctx, devicesToAdd)
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
			ByUpdatableAssetId(state.UpdatableAssetGroupId.ValueString()).
			MicrosoftGraphWindowsUpdatesAddMembersById().
			Post(ctx, addRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	if len(devicesToRemove) > 0 {
		removeRequest, err := constructRemoveMembersRequest(ctx, devicesToRemove)
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
			ByUpdatableAssetId(state.UpdatableAssetGroupId.ValueString()).
			MicrosoftGraphWindowsUpdatesRemoveMembersById().
			Post(ctx, removeRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Windows Updates autopatch updatable asset group assignment resources.
//
// Operation: Removes all device members from an updatable asset group
// API Calls:
//   - POST /admin/windows/updates/updatableAssets/{updatableAssetId}/microsoft.graph.windowsUpdates.removeMembersById
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableassetgroup-removemembersbyid?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deviceIDs := extractDeviceIDs(&object)

	if len(deviceIDs) > 0 {
		requestBody, err := constructRemoveMembersRequest(ctx, deviceIDs)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing remove members request",
				fmt.Sprintf("Could not construct remove members request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.Admin().
			Windows().
			Updates().
			UpdatableAssets().
			ByUpdatableAssetId(object.UpdatableAssetGroupId.ValueString()).
			MicrosoftGraphWindowsUpdatesRemoveMembersById().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
