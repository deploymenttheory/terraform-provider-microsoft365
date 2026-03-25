package graphBetaAdministrativeUnitMembership

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for administrative unit membership resources.
//
// Operation: Adds members to an administrative unit
// API Calls:
//   - POST /directory/administrativeUnits/{id}/members/$ref (for each member)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/administrativeunit-post-members?view=graph-rest-beta
func (r *AdministrativeUnitMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AdministrativeUnitMembershipResourceModel

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

	administrativeUnitID := object.AdministrativeUnitID.ValueString()
	object.ID = types.StringValue(administrativeUnitID)

	memberIDs := extractMemberIDsFromSet(object.Members)

	if !validateRequest(ctx, r.client, memberIDs, &resp.Diagnostics) {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Adding %d members to administrative unit %s", len(memberIDs), administrativeUnitID))

	for _, memberID := range memberIDs {
		referenceBody := createReferenceRequest(memberID)

		err := r.client.
			Directory().
			AdministrativeUnits().
			ByAdministrativeUnitId(administrativeUnitID).
			Members().
			Ref().
			Post(ctx, referenceBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully added member %s to administrative unit %s", memberID, administrativeUnitID))
	}

	if len(memberIDs) > 0 {
		tflog.Debug(ctx, "Waiting 20 seconds for member additions to propagate")
		time.Sleep(20 * time.Second)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
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

// Read handles the Read operation for administrative unit membership resources.
//
// Operation: Retrieves members of an administrative unit
// API Calls:
//   - GET /directory/administrativeUnits/{id}/members/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/administrativeunit-list-members?view=graph-rest-beta
func (r *AdministrativeUnitMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AdministrativeUnitMembershipResourceModel
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

	administrativeUnitID := object.AdministrativeUnitID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s for administrative unit: %s", ResourceName, administrativeUnitID))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Retry logic to handle eventual consistency
	var memberIDs []string
	maxRetries := 3
	retryDelay := 5 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Retry attempt %d/%d after %v delay", attempt+1, maxRetries, retryDelay))
			time.Sleep(retryDelay)
		}

		tflog.Debug(ctx, fmt.Sprintf("Calling GET /directory/administrativeUnits/%s/members (attempt %d)", administrativeUnitID, attempt+1))

		membersResp, err := r.client.
			Directory().
			AdministrativeUnits().
			ByAdministrativeUnitId(administrativeUnitID).
			Members().
			Get(ctx, nil)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error calling GET /members: %v", err))
			errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
			return
		}

		tflog.Debug(ctx, "GET /members request succeeded")

		memberIDs = make([]string, 0)
		if membersResp != nil {
			tflog.Debug(ctx, "Members response is not nil")
			valueSlice := membersResp.GetValue()
			if valueSlice != nil {
				tflog.Debug(ctx, fmt.Sprintf("Members response has %d items", len(valueSlice)))
				for i, member := range valueSlice {
					tflog.Trace(ctx, fmt.Sprintf("Processing member %d: %v", i, member))
					if member != nil {
						memberIDPtr := member.GetId()
						if memberIDPtr != nil {
							memberID := *memberIDPtr
							memberIDs = append(memberIDs, memberID)
							tflog.Trace(ctx, fmt.Sprintf("Found member ID: %s", memberID))
						} else {
							tflog.Warn(ctx, fmt.Sprintf("Member %d has nil ID", i))
						}
					} else {
						tflog.Warn(ctx, fmt.Sprintf("Member %d is nil", i))
					}
				}
			} else {
				tflog.Warn(ctx, "Members response GetValue() returned nil")
			}
		} else {
			tflog.Warn(ctx, "Members response is nil")
		}

		tflog.Debug(ctx, fmt.Sprintf("Extracted %d member IDs from API response", len(memberIDs)))

		// If we got members or this is the last attempt, break
		if len(memberIDs) > 0 || attempt == maxRetries-1 {
			break
		}

		tflog.Debug(ctx, "No members found, will retry")
	}

	MapRemoteStateToTerraform(ctx, &object, memberIDs)

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

// Update handles the Update operation for administrative unit membership resources.
//
// Operation: Updates members of an administrative unit by adding/removing as needed
// API Calls:
//   - POST /directory/administrativeUnits/{id}/members/$ref (for each member to add)
//   - DELETE /directory/administrativeUnits/{id}/members/{memberId}/$ref (for each member to remove)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/administrativeunit-post-members?view=graph-rest-beta
// Reference: https://learn.microsoft.com/en-us/graph/api/administrativeunit-delete-members?view=graph-rest-beta
func (r *AdministrativeUnitMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AdministrativeUnitMembershipResourceModel
	var state AdministrativeUnitMembershipResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	administrativeUnitID := state.AdministrativeUnitID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Updating %s for administrative unit: %s", ResourceName, administrativeUnitID))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	toAdd, toRemove := calculateMembershipChanges(plan.Members, state.Members)

	if !validateRequest(ctx, r.client, toAdd, &resp.Diagnostics) {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Membership changes: %d to add, %d to remove", len(toAdd), len(toRemove)))

	for _, memberID := range toAdd {
		referenceBody := createReferenceRequest(memberID)

		err := r.client.
			Directory().
			AdministrativeUnits().
			ByAdministrativeUnitId(administrativeUnitID).
			Members().
			Ref().
			Post(ctx, referenceBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully added member %s", memberID))
	}

	for _, memberID := range toRemove {
		err := r.client.
			Directory().
			AdministrativeUnits().
			ByAdministrativeUnitId(administrativeUnitID).
			Members().
			ByDirectoryObjectId(memberID).
			Ref().
			Delete(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully removed member %s", memberID))
	}

	if len(toAdd) > 0 || len(toRemove) > 0 {
		tflog.Debug(ctx, "Waiting 20 seconds for membership changes to propagate")
		time.Sleep(20 * time.Second)
	}

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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for administrative unit membership resources.
//
// Operation: Removes all members from the administrative unit
// API Calls:
//   - DELETE /directory/administrativeUnits/{id}/members/{memberId}/$ref (for each member)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/administrativeunit-delete-members?view=graph-rest-beta
func (r *AdministrativeUnitMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AdministrativeUnitMembershipResourceModel

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

	administrativeUnitID := object.AdministrativeUnitID.ValueString()
	memberIDs := extractMemberIDsFromSet(object.Members)

	tflog.Debug(ctx, fmt.Sprintf("Removing %d members from administrative unit %s", len(memberIDs), administrativeUnitID))

	for _, memberID := range memberIDs {
		err := r.client.
			Directory().
			AdministrativeUnits().
			ByAdministrativeUnitId(administrativeUnitID).
			Members().
			ByDirectoryObjectId(memberID).
			Ref().
			Delete(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully removed member %s from administrative unit %s", memberID, administrativeUnitID))
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
