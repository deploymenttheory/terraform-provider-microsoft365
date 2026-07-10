package graphBetaGroupLicenseAssignment

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
)

// groupLicenseLocks serializes assignLicense operations per group.
//
// The assignLicense endpoint mutates the group's whole license set and is asynchronous:
// it returns 202 Accepted and the actual add/remove is applied later by the backend
// licensing service (see the Response section of
// https://learn.microsoft.com/en-us/graph/api/group-assignlicense?view=graph-rest-beta).
// Microsoft documents that group license processing takes time to complete and can fail
// after the request was accepted:
//   - https://learn.microsoft.com/en-us/entra/identity/users/licensing-group-advanced
//     ("How long does it take for licenses to be modified after group changes")
//   - https://learn.microsoft.com/en-us/entra/fundamentals/licensing-groups-resolve-problems
//
// There is no official documentation of the concurrency semantics of overlapping
// assignLicense calls against the same group. Because each call rewrites the same license
// set asynchronously, concurrent POSTs from multiple Terraform resources targeting the same
// group (Terraform applies resources in parallel by default) have been observed in the field
// to silently lose adds or removes. Serializing per group is therefore a defensive measure
// consistent with Microsoft's eventual-consistency guidance for Entra
// (https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/):
// each Create/Update/Delete takes the group's lock and holds it until the change is
// confirmed by a read-back.
var groupLicenseLocks sync.Map // group id -> *sync.Mutex

// lockGroupLicense acquires the per-group license lock and returns the unlock func.
// Callers acquire it before starting their operation timeout so time spent waiting on
// sibling license assignments does not consume the resource's own timeout budget.
func lockGroupLicense(ctx context.Context, groupID string) func() {
	m, _ := groupLicenseLocks.LoadOrStore(groupID, &sync.Mutex{})
	mu := m.(*sync.Mutex)
	start := time.Now()
	mu.Lock()
	if waited := time.Since(start); waited > time.Second {
		tflog.Debug(ctx, fmt.Sprintf("Waited %s for the license assignment lock on group %s", waited, groupID))
	}
	return mu.Unlock
}

// Create handles the Create operation for Group License Assignment resources.
//
// Operation: Assigns a license to a group (group-based licensing)
// API Calls:
//   - POST /groups/{groupId}/assignLicense
//
// Reference: https://learn.microsoft.com/en-us/graph/api/group-assignlicense?view=graph-rest-beta
func (r *GroupLicenseAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize with sibling license assignments targeting the same group (see groupLicenseLocks).
	unlock := lockGroupLicense(ctx, object.GroupId.ValueString())
	defer unlock()

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Create composite ID: group_id_sku_id
	object.ID = types.StringValue(fmt.Sprintf("%s_%s", object.GroupId.ValueString(), object.SkuId.ValueString()))

	// Ensure disabled_plans is set to empty set if not provided (can't be unknown)
	if object.DisabledPlans.IsNull() || object.DisabledPlans.IsUnknown() {
		object.DisabledPlans = types.SetValueMust(types.StringType, []attr.Value{})
	}

	requestBody, err := constructAddLicensesRequest(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing license assignment request",
			fmt.Sprintf("Could not construct license assignment request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Groups().
		ByGroupId(object.GroupId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned license to group: %s", object.GroupId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = groupLicenseAssignmentConsistencyPredicate(&object)

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

// Read handles the Read operation for Group License Assignment resources.
//
// Operation: Retrieves group license assignments
// API Calls:
//   - GET /groups/{groupId}?$select=id,displayName,assignedLicenses
//
// Reference: https://learn.microsoft.com/en-us/graph/api/group-get?view=graph-rest-beta
func (r *GroupLicenseAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupLicenseAssignmentResourceModel
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

	tflog.Debug(ctx, fmt.Sprintf("Reading group license assignments for group ID: %s", object.GroupId.ValueString()))

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

	requestParameters := &groups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "displayName", "assignedLicenses"},
		},
	}

	group, err := r.client.
		Groups().
		ByGroupId(object.GroupId.ValueString()).
		Get(ctx, requestParameters)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	// The group existing is not enough — this resource represents a single license (SKU)
	// on the group. A single read that lacks the SKU is not authoritative, however:
	// group license writes are asynchronous and reads can be served by lagging replicas.
	// Before removing state, re-check until the read timeout expires so a transient stale
	// read does not make Terraform forget a live license assignment and later try to delete
	// the still-licensed group.
	if found := MapRemoteResourceStateToTerraform(ctx, &object, group); !found {
		if err := r.waitForLicensePresence(ctx, object.GroupId.ValueString(), object.SkuId.ValueString()); err == nil {
			group, err = r.client.
				Groups().
				ByGroupId(object.GroupId.ValueString()).
				Get(ctx, requestParameters)
			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
				return
			}
			if found := MapRemoteResourceStateToTerraform(ctx, &object, group); found {
				resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
				return
			}
		}

		tflog.Warn(ctx, fmt.Sprintf("License SKU %s is not assigned to group %s, removing %s from state",
			object.SkuId.ValueString(), object.GroupId.ValueString(), ResourceName))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Group License Assignment resources.
//
// Operation: Updates license assignment (primarily disabled plans, since sku_id changes require replacement)
// API Calls:
//   - POST /groups/{groupId}/assignLicense
//
// Reference: https://learn.microsoft.com/en-us/graph/api/group-assignlicense?view=graph-rest-beta
// Note: Same API endpoint as Create; sku_id is marked as RequiresReplace, so updates mainly modify disabled_plans
func (r *GroupLicenseAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize with sibling license assignments targeting the same group (see groupLicenseLocks).
	unlock := lockGroupLicense(ctx, plan.GroupId.ValueString())
	defer unlock()

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Update the license (mainly disabled_plans since sku_id has RequiresReplace)
	requestBody, err := constructUpdateLicenseRequest(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing license assignment request for update",
			fmt.Sprintf("Could not construct license assignment request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Groups().
		ByGroupId(plan.GroupId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully updated license for group: %s", plan.GroupId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = groupLicenseAssignmentConsistencyPredicate(&plan)

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, plan.ID.ValueString()))
}

// Delete handles the Delete operation for Group License Assignment resources.
//
// Operation: Removes a license from a group
// API Calls:
//   - POST /groups/{groupId}/assignLicense (with removeLicenses parameter)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/group-assignlicense?view=graph-rest-beta
// Note: Uses same assignLicense endpoint with removeLicenses to remove the license.
// assignLicense is asynchronous — it returns 202 Accepted and the removal is processed later
// by the backend licensing service (https://learn.microsoft.com/en-us/graph/api/group-assignlicense?view=graph-rest-beta).
// Microsoft documents that a group cannot be deleted until license removal processing has
// actually completed ("A group with active licenses assigned cannot be deleted"), and that
// removals can fail after being accepted, e.g. when removed SKUs contain service plans that
// other assigned SKUs depend on:
//   - https://learn.microsoft.com/en-us/entra/fundamentals/licensing-groups-resolve-problems
//   - https://learn.microsoft.com/en-us/answers/questions/2006508/how-to-remove-multiple-licenses-from-a-security-gr
//
// Delete therefore polls the group's assignedLicenses until the SKU is actually gone before
// removing the resource from state. Returning early on the 202 would let a dependent group
// deletion run while the group still has active licenses, and a silently failed removal would
// be unrecoverable via Terraform because the resource would already be gone from state.
func (r *GroupLicenseAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize with sibling license assignments targeting the same group (see groupLicenseLocks).
	unlock := lockGroupLicense(ctx, object.GroupId.ValueString())
	defer unlock()

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Remove the single license managed by this resource
	requestBody, err := constructRemoveLicenseRequest(ctx, object.SkuId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing license removal request",
			fmt.Sprintf("Could not construct license removal request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Groups().
		ByGroupId(object.GroupId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("License removal request accepted for license %s on group: %s, waiting for removal to complete",
		object.SkuId.ValueString(), object.GroupId.ValueString()))

	if err := r.waitForLicenseRemoval(ctx, object.GroupId.ValueString(), object.SkuId.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"License removal did not complete",
			fmt.Sprintf("The removal of license %s from group %s was accepted by the API but the license "+
				"is still assigned after waiting. Group license processing is asynchronous and can fail "+
				"silently (e.g. when removed SKUs contain service plans that other assigned SKUs depend on) "+
				"or become stuck. Check the group's license processing state in the Microsoft Entra portal "+
				"and reprocess if needed, then retry. The resource has been kept in Terraform state so the "+
				"removal can be retried. Error: %s",
				object.SkuId.ValueString(), object.GroupId.ValueString(), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed license %s from group: %s", object.SkuId.ValueString(), object.GroupId.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// waitForLicenseRemoval polls the group's assignedLicenses until the given SKU is no longer
// present, the group itself is gone (404), or the context deadline is reached. This closes the
// gap between the asynchronous 202 Accepted returned by assignLicense and the license actually
// being removed by the backend licensing service.
func (r *GroupLicenseAssignmentResource) waitForLicenseRemoval(ctx context.Context, groupID, skuID string) error {
	const pollInterval = 5 * time.Second
	const requiredAbsentReads = 2
	absentReads := 0

	requestParameters := &groups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "assignedLicenses"},
		},
	}

	return crud.PollUntil(ctx, pollInterval, func(ctx context.Context) (bool, error) {
		group, err := r.client.
			Groups().
			ByGroupId(groupID).
			Get(ctx, requestParameters)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 {
				tflog.Debug(ctx, fmt.Sprintf("Group %s no longer exists, treating license removal as complete", groupID))
				return true, nil
			}
			// Permanent errors (authorization, bad request, ...) will not resolve by
			// polling — fail fast instead of consuming the whole delete timeout.
			if errors.IsNonRetryableReadError(&errorInfo) {
				return false, &crud.FatalPollError{Err: err}
			}
			// Transient read errors should not fail the delete outright — keep polling
			// until the deadline and surface the last error if the wait times out.
			tflog.Debug(ctx, fmt.Sprintf("Error reading group %s while waiting for license removal, will retry: %s", groupID, err.Error()))
			return false, err
		}

		found := false
		for _, license := range group.GetAssignedLicenses() {
			if license == nil || license.GetSkuId() == nil {
				continue
			}
			// The API returns SKU ids in canonical lowercase form while the configured
			// sku_id may use any casing, so compare case-insensitively.
			if strings.EqualFold(license.GetSkuId().String(), skuID) {
				found = true
				absentReads = 0
				tflog.Debug(ctx, fmt.Sprintf("License %s still assigned to group %s, waiting %s before re-checking", skuID, groupID, pollInterval))
				return false, fmt.Errorf("license %s is still assigned to group %s", skuID, groupID)
			}
		}
		if !found {
			absentReads++
		}
		if absentReads < requiredAbsentReads {
			return false, fmt.Errorf("license %s was absent from group %s on read %d/%d; confirming removal",
				skuID, groupID, absentReads, requiredAbsentReads)
		}
		return true, nil
	})
}

// waitForLicensePresence polls until the given SKU is visible in assignedLicenses. It is
// used by Read before removing state so one stale read cannot turn an eventually-consistent
// live assignment into Terraform drift.
func (r *GroupLicenseAssignmentResource) waitForLicensePresence(ctx context.Context, groupID, skuID string) error {
	const pollInterval = 5 * time.Second

	requestParameters := &groups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "assignedLicenses"},
		},
	}

	return crud.PollUntil(ctx, pollInterval, func(ctx context.Context) (bool, error) {
		group, err := r.client.
			Groups().
			ByGroupId(groupID).
			Get(ctx, requestParameters)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errors.IsNonRetryableReadError(&errorInfo) {
				return false, &crud.FatalPollError{Err: err}
			}
			return false, err
		}

		for _, license := range group.GetAssignedLicenses() {
			if license == nil || license.GetSkuId() == nil {
				continue
			}
			if strings.EqualFold(license.GetSkuId().String(), skuID) {
				return true, nil
			}
		}
		return false, fmt.Errorf("license %s is not yet visible on group %s", skuID, groupID)
	})
}
