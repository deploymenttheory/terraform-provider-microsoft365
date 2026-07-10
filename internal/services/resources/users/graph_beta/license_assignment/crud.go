package graphBetaUserLicenseAssignment

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
)

// userLicenseLocks serializes assignLicense operations per user.
//
// The assignLicense endpoint mutates the user's whole license set
// (https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta),
// and Microsoft Entra is an eventually consistent, multi-replica system
// (https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/).
// There is no official documentation of the concurrency semantics of overlapping
// assignLicense calls against the same user; because each call rewrites the same license
// set, concurrent POSTs from multiple Terraform resources targeting the same user
// (Terraform applies resources in parallel by default) risk losing adds or removes, the
// same failure mode observed in the field for the group variant of this resource.
// Each Create/Update/Delete takes the user's lock and holds it until the change is
// confirmed by a read-back.
var userLicenseLocks sync.Map // user id -> *sync.Mutex

// lockUserLicense acquires the per-user license lock and returns the unlock func.
// Callers acquire it before starting their operation timeout so time spent waiting on
// sibling license assignments does not consume the resource's own timeout budget.
func lockUserLicense(ctx context.Context, userID string) func() {
	m, _ := userLicenseLocks.LoadOrStore(userID, &sync.Mutex{})
	mu := m.(*sync.Mutex)
	start := time.Now()
	mu.Lock()
	if waited := time.Since(start); waited > time.Second {
		tflog.Debug(ctx, fmt.Sprintf("Waited %s for the license assignment lock on user %s", waited, userID))
	}
	return mu.Unlock
}

// Create handles the Create operation for user license assignment resources.
//
// Operation: Assigns a license to a user
// API Calls:
//   - POST /users/{id}/assignLicense
//
// Reference: https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta
// Note: Composite ID (userId_skuId) is constructed as Graph API does not return unique assignment IDs
func (r *UserLicenseAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object UserLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize with sibling license assignments targeting the same user (see userLicenseLocks).
	unlock := lockUserLicense(ctx, object.UserId.ValueString())
	defer unlock()

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Create composite ID: user_id_sku_id
	object.ID = types.StringValue(fmt.Sprintf("%s_%s", object.UserId.ValueString(), object.SkuId.ValueString()))

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
		Users().
		ByUserId(object.UserId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned licenses to user: %s", object.UserId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = licenseAssignmentConsistencyPredicate(&object)

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

// Read handles the Read operation for user license assignment resources.
//
// Operation: Retrieves user's assigned licenses to verify assignment exists
// API Calls:
//   - GET /users/{id}?$select=id,userPrincipalName,assignedLicenses
//
// Reference: https://learn.microsoft.com/en-us/graph/api/user-get?view=graph-rest-beta
func (r *UserLicenseAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object UserLicenseAssignmentResourceModel
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

	tflog.Debug(ctx, fmt.Sprintf("Reading user license assignments for user ID: %s", object.UserId.ValueString()))

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

	requestParameters := &users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "userPrincipalName", "assignedLicenses"},
		},
	}

	user, err := r.client.
		Users().
		ByUserId(object.UserId.ValueString()).
		Get(ctx, requestParameters)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	// The user existing is not enough — this resource represents a single license (SKU)
	// on the user. When the SKU is absent from assignedLicenses, the license assignment
	// does not exist, so remove it from state. During post-write ReadWithRetry polling this
	// causes the consistency predicate to fail and the read to be retried, which is how a
	// not-yet-propagated (or failed) assignLicense is detected instead of being silently
	// reported as successful.
	if found := MapRemoteResourceStateToTerraform(ctx, &object, user); !found {
		tflog.Warn(ctx, fmt.Sprintf("License SKU %s is not assigned to user %s, removing %s from state",
			object.SkuId.ValueString(), object.UserId.ValueString(), ResourceName))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for user license assignment resources.
//
// Operation: Updates disabled plans for an existing license assignment
// API Calls:
//   - POST /users/{id}/assignLicense
//
// Reference: https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta
// Note: Only disabled_plans can be updated; sku_id changes trigger resource replacement
func (r *UserLicenseAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize with sibling license assignments targeting the same user (see userLicenseLocks).
	unlock := lockUserLicense(ctx, plan.UserId.ValueString())
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
		Users().
		ByUserId(plan.UserId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully updated license for user: %s", plan.UserId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = licenseAssignmentConsistencyPredicate(&plan)

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

// Delete handles the Delete operation for user license assignment resources.
//
// Operation: Removes a license from a user
// API Calls:
//   - POST /users/{id}/assignLicense
//
// Reference: https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta
// Note: License removal is performed by passing the skuId in removeLicenses array
func (r *UserLicenseAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object UserLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize with sibling license assignments targeting the same user (see userLicenseLocks).
	unlock := lockUserLicense(ctx, object.UserId.ValueString())
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
		Users().
		ByUserId(object.UserId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("License removal request accepted for license %s on user: %s, waiting for removal to be reflected",
		object.SkuId.ValueString(), object.UserId.ValueString()))

	if err := r.waitForLicenseRemoval(ctx, object.UserId.ValueString(), object.SkuId.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"License removal did not complete",
			fmt.Sprintf("The removal of license %s from user %s was accepted by the API but the license "+
				"is still assigned after waiting. If the license is inherited from a group the user is a "+
				"member of, it cannot be removed directly — remove it from the group instead. The resource "+
				"has been kept in Terraform state so the removal can be retried. Error: %s",
				object.SkuId.ValueString(), object.UserId.ValueString(), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed license %s from user: %s", object.SkuId.ValueString(), object.UserId.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// waitForLicenseRemoval polls the user's assignedLicenses until the given SKU is no longer
// present, the user itself is gone (404), or the context deadline is reached. Unlike the
// group variant, user assignLicense applies synchronously (200 OK), but reads may still be
// served from a replica that has not yet received the write
// (https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/),
// so removing the resource from state without confirming would let a dependent operation
// (e.g. deleting the user's last license before disabling the account) observe stale data.
func (r *UserLicenseAssignmentResource) waitForLicenseRemoval(ctx context.Context, userID, skuID string) error {
	const pollInterval = 5 * time.Second

	requestParameters := &users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "assignedLicenses"},
		},
	}

	return crud.PollUntil(ctx, pollInterval, func(ctx context.Context) (bool, error) {
		user, err := r.client.
			Users().
			ByUserId(userID).
			Get(ctx, requestParameters)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 {
				tflog.Debug(ctx, fmt.Sprintf("User %s no longer exists, treating license removal as complete", userID))
				return true, nil
			}
			// Permanent errors (authorization, bad request, ...) will not resolve by
			// polling — fail fast instead of consuming the whole delete timeout.
			if errors.IsNonRetryableReadError(&errorInfo) {
				return false, &crud.FatalPollError{Err: err}
			}
			// Transient read errors should not fail the delete outright — keep polling
			// until the deadline and surface the last error if the wait times out.
			tflog.Debug(ctx, fmt.Sprintf("Error reading user %s while waiting for license removal, will retry: %s", userID, err.Error()))
			return false, err
		}

		for _, license := range user.GetAssignedLicenses() {
			if license == nil || license.GetSkuId() == nil {
				continue
			}
			// The API returns SKU ids in canonical lowercase form while the configured
			// sku_id may use any casing, so compare case-insensitively.
			if strings.EqualFold(license.GetSkuId().String(), skuID) {
				tflog.Debug(ctx, fmt.Sprintf("License %s still assigned to user %s, waiting %s before re-checking", skuID, userID, pollInterval))
				return false, fmt.Errorf("license %s is still assigned to user %s", skuID, userID)
			}
		}
		return true, nil
	})
}
