package graphBetaDeviceAndAppManagementAndroidManagedMobileApp

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

	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
)

// NOTE ON THIS FILE'S APPROACH:
//
// The androidManagedAppProtections/{id}/apps navigation collection does not
// support POST, PATCH, or GET-by-key on this tenant's backend - confirmed
// directly against the live API ("No OData route exists ... with http verb
// POST/GET for request .../apps" and ".../apps/{id}"). Only two operations
// on this collection actually work: GET the full collection (no key), and
// the targetApps action, which REPLACES the entire targeted-app list for
// the policy in one call - it is not additive.
//
// So every write here follows the same pattern: read the full current app
// list, modify it in memory (add/remove/update one entry), then POST the
// complete resulting list via deviceAppManagement/managedAppPolicies/{id}/
// targetApps (the generic entity set - the Android-specific SDK builder for
// targetApps does not exist in this SDK version, but the generic one does
// and works identically since androidManagedAppProtection is a
// managedAppPolicy under the hood).
//
// IMPORTANT: because targetApps replaces the whole list, applying multiple
// android_managed_mobile_app resources against the SAME policy in parallel
// is a race - two simultaneous Creates can each read the list before the
// other's addition lands, and whichever targetApps call finishes last wins,
// silently dropping the other app. Apply resources sharing a policy one at
// a time (e.g. -target, or terraform apply -parallelism=1) until this
// resource is reworked to serialize writes per policy.

// Create handles the Create operation for Android Managed Mobile App resources.
//
// Operation: Adds a managed mobile app to an Android app protection policy
// API Calls:
//   - GET  /deviceAppManagement/androidManagedAppProtections/{androidManagedAppProtectionId}/apps
//   - POST /deviceAppManagement/managedAppPolicies/{managedAppPolicyId}/targetApps
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-managedmobileapp-create?view=graph-rest-beta
func (r *AndroidManagedMobileAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AndroidManagedMobileAppResourceModel

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

	newApp, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	policyId := object.ManagedAppProtectionId.ValueString()

	currentApps, err := r.client.
		DeviceAppManagement().
		AndroidManagedAppProtections().
		ByAndroidManagedAppProtectionId(policyId).
		Apps().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}

	existingIds := make(map[string]bool)
	for _, app := range currentApps.GetValue() {
		if app.GetId() != nil {
			existingIds[*app.GetId()] = true
		}
	}

	updatedApps := append(currentApps.GetValue(), newApp)

	targetBody := deviceappmanagement.NewManagedAppPoliciesItemTargetAppsPostRequestBody()
	targetBody.SetApps(updatedApps)

	err = r.client.
		DeviceAppManagement().
		ManagedAppPolicies().
		ByManagedAppPolicyId(policyId).
		TargetApps().
		Post(ctx, targetBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// targetApps returns no content, so re-read the collection and diff
	// against the pre-call list to find the server-assigned id for the app
	// we just added.
	refreshedApps, err := r.client.
		DeviceAppManagement().
		AndroidManagedAppProtections().
		ByAndroidManagedAppProtectionId(policyId).
		Apps().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}

	var createdId string
	for _, app := range refreshedApps.GetValue() {
		if app.GetId() != nil && !existingIds[*app.GetId()] {
			createdId = *app.GetId()
			break
		}
	}

	if createdId == "" {
		resp.Diagnostics.AddError(
			"Error locating created resource",
			fmt.Sprintf("targetApps call succeeded but could not find a new app entry in policy %s afterward", policyId),
		)
		return
	}

	object.ID = types.StringValue(createdId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
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

// Read handles the Read operation for Android Managed Mobile App resources.
//
// Operation: Retrieves a managed mobile app from an Android app protection policy by ID
// API Calls:
//   - GET /deviceAppManagement/androidManagedAppProtections/{androidManagedAppProtectionId}/apps
//
// Per-item GET (.../apps/{managedMobileAppId}) is not supported by this
// tenant's backend, so the full collection is read and filtered client-side.
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-managedmobileapp-get?view=graph-rest-beta
func (r *AndroidManagedMobileAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AndroidManagedMobileAppResourceModel
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

	allApps, err := r.client.
		DeviceAppManagement().
		AndroidManagedAppProtections().
		ByAndroidManagedAppProtectionId(object.ManagedAppProtectionId.ValueString()).
		Apps().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	var found bool
	for _, app := range allApps.GetValue() {
		if app.GetId() != nil && *app.GetId() == object.ID.ValueString() {
			MapRemoteStateToTerraform(ctx, &object, app)
			found = true
			break
		}
	}

	if !found {
		tflog.Debug(ctx, fmt.Sprintf("%s with ID %s no longer present in policy apps list, removing from state", ResourceName, object.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Android Managed Mobile App resources.
//
// Operation: Updates a managed mobile app entry in an Android app protection policy
// API Calls:
//   - GET  /deviceAppManagement/androidManagedAppProtections/{androidManagedAppProtectionId}/apps
//   - POST /deviceAppManagement/managedAppPolicies/{managedAppPolicyId}/targetApps
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-managedmobileapp-update?view=graph-rest-beta
func (r *AndroidManagedMobileAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AndroidManagedMobileAppResourceModel
	var state AndroidManagedMobileAppResourceModel

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

	policyId := state.ManagedAppProtectionId.ValueString()

	updatedApp, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	updatedApp.SetId(state.ID.ValueStringPointer())

	currentApps, err := r.client.
		DeviceAppManagement().
		AndroidManagedAppProtections().
		ByAndroidManagedAppProtectionId(policyId).
		Apps().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}

	apps := currentApps.GetValue()
	merged := apps[:0]
	replaced := false
	for _, app := range apps {
		if app.GetId() != nil && *app.GetId() == state.ID.ValueString() {
			merged = append(merged, updatedApp)
			replaced = true
		} else {
			merged = append(merged, app)
		}
	}
	if !replaced {
		merged = append(merged, updatedApp)
	}

	targetBody := deviceappmanagement.NewManagedAppPoliciesItemTargetAppsPostRequestBody()
	targetBody.SetApps(merged)

	err = r.client.
		DeviceAppManagement().
		ManagedAppPolicies().
		ByManagedAppPolicyId(policyId).
		TargetApps().
		Post(ctx, targetBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
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

// Delete handles the Delete operation for Android Managed Mobile App resources.
//
// Operation: Removes a managed mobile app entry from an Android app protection policy
// API Calls:
//   - GET  /deviceAppManagement/androidManagedAppProtections/{androidManagedAppProtectionId}/apps
//   - POST /deviceAppManagement/managedAppPolicies/{managedAppPolicyId}/targetApps
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-managedmobileapp-delete?view=graph-rest-beta
func (r *AndroidManagedMobileAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AndroidManagedMobileAppResourceModel

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

	policyId := object.ManagedAppProtectionId.ValueString()

	currentApps, err := r.client.
		DeviceAppManagement().
		AndroidManagedAppProtections().
		ByAndroidManagedAppProtectionId(policyId).
		Apps().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		return
	}

	apps := currentApps.GetValue()
	remaining := apps[:0]
	for _, app := range apps {
		if app.GetId() == nil || *app.GetId() != object.ID.ValueString() {
			remaining = append(remaining, app)
		}
	}

	targetBody := deviceappmanagement.NewManagedAppPoliciesItemTargetAppsPostRequestBody()
	targetBody.SetApps(remaining)

	err = r.client.
		DeviceAppManagement().
		ManagedAppPolicies().
		ByManagedAppPolicyId(policyId).
		TargetApps().
		Post(ctx, targetBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
