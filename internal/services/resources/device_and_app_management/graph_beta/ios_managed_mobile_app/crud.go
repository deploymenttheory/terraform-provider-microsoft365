package graphBetaDeviceAndAppManagementIOSManagedMobileApp

import (
	"context"
	stderrors "errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// managedAppPolicyTargetAppsLockScope serializes targetApps writes per policy.
//
// The iosManagedAppProtections/{id}/apps collection supports only a GET of the
// whole list; per-app GET, POST and PATCH return "No OData route exists". Writes go
// through managedAppPolicies/{id}/targetApps, which REPLACES the policy's entire
// app list. Every write therefore reads the list, changes one entry and posts the
// full result while holding this lock for the policy, so resources that share a
// policy cannot overwrite each other's apps.
const managedAppPolicyTargetAppsLockScope = "managed_app_policy_target_apps"

var (
	errAppAlreadyTargeted = stderrors.New("the app is already targeted by this policy")
	errAppWaitTimeout     = stderrors.New("timed out waiting for the policy app list to reflect the change")

	// appPollInterval is how often the app list is re-read after a write.
	appPollInterval = 2 * time.Second
)

// Create adds the app to the policy's targeted apps.
//
// API Calls:
//   - GET  /deviceAppManagement/iosManagedAppProtections/{id}/apps
//   - POST /deviceAppManagement/managedAppPolicies/{id}/targetApps
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-targetedmanagedappprotection-targetapps?view=graph-rest-beta
func (r *IOSManagedMobileAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object IOSManagedMobileAppResourceModel

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

	policyID := object.ManagedAppProtectionId.ValueString()
	bundle := object.MobileAppIdentifier.BundleId.ValueString()

	// Serialize with sibling apps targeting the same policy (see managedAppPolicyTargetAppsLockScope).
	unlock, err := crud.AcquireKeyedLock(ctx, managedAppPolicyTargetAppsLockScope, policyID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
		return
	}
	defer unlock()

	appsResponse, err := r.client.DeviceAppManagement().IosManagedAppProtections().ByIosManagedAppProtectionId(policyID).Apps().Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, fmt.Errorf("list policy apps: %w", err), resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	currentApps := appsResponse.GetValue()

	if existing := findTargetedAppByBundleID(currentApps, bundle); existing != nil {
		existingID := bundle + ".ios"
		if existing.GetId() != nil {
			existingID = *existing.GetId()
		}
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("%s: bundle %q in policy %s. Import it with ID \"%s/%s\" instead.", errAppAlreadyTargeted, bundle, policyID, policyID, existingID),
		)
		return
	}

	desired := append(slices.Clone(currentApps), newApp)

	if err := r.client.DeviceAppManagement().ManagedAppPolicies().ByManagedAppPolicyId(policyID).TargetApps().Post(ctx, constructTargetAppsRequest(desired), nil); err != nil {
		errors.HandleKiotaGraphError(ctx, fmt.Errorf("set policy target apps: %w", err), resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// targetApps returns no content and Intune assigns the id, so poll the list until the entry has one.
	pollReadFailed := false
	err = crud.PollUntil(ctx, appPollInterval, func(ctx context.Context) (bool, error) {
		response, readErr := r.client.DeviceAppManagement().IosManagedAppProtections().ByIosManagedAppProtectionId(policyID).Apps().Get(ctx, nil)
		if readErr != nil {
			pollReadFailed = true
			return false, &crud.FatalPollError{Err: fmt.Errorf("list policy apps: %w", readErr)}
		}
		created := findTargetedAppByBundleID(response.GetValue(), bundle)
		if created == nil || created.GetId() == nil {
			return false, nil
		}
		object.ID = types.StringValue(*created.GetId())
		return true, nil
	})
	if err != nil && ctx.Err() != nil && !pollReadFailed {
		err = fmt.Errorf("%w: %w", errAppWaitTimeout, ctx.Err())
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not find bundle %q in policy %s: %s", bundle, policyID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read finds the app in the policy's targeted apps.
//
// Per-app GET (.../apps/{id}) is not supported, so the full collection is read.
//
// API Calls:
//   - GET /deviceAppManagement/iosManagedAppProtections/{id}/apps
func (r *IOSManagedMobileAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object IOSManagedMobileAppResourceModel
	var identity sharedmodels.ResourceIdentity

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	policyID := object.ManagedAppProtectionId.ValueString()
	appID := object.ID.ValueString()

	appsResponse, err := r.client.DeviceAppManagement().IosManagedAppProtections().ByIosManagedAppProtectionId(policyID).Apps().Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, fmt.Errorf("list policy apps: %w", err), resp, constants.TfOperationRead, r.ReadPermissions)
		return
	}
	apps := appsResponse.GetValue()

	index := slices.IndexFunc(apps, func(app graphmodels.ManagedMobileAppable) bool {
		return app != nil && app.GetId() != nil && strings.EqualFold(*app.GetId(), appID)
	})
	if index < 0 {
		tflog.Warn(ctx, fmt.Sprintf("%s with ID %s is no longer targeted by policy %s, removing from state", ResourceName, appID, policyID))
		resp.State.RemoveResource(ctx)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, apps[index])

	identity.ID = object.ID.ValueString()
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update only applies Terraform-side changes such as timeouts. Every attribute sent
// to Graph forces replacement and Intune assigns version itself, so the prior state is
// kept and only the planned timeouts are applied to it.
func (r *IOSManagedMobileAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IOSManagedMobileAppResourceModel
	var state IOSManagedMobileAppResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete removes the app from the policy's targeted apps.
//
// API Calls:
//   - GET  /deviceAppManagement/iosManagedAppProtections/{id}/apps
//   - POST /deviceAppManagement/managedAppPolicies/{id}/targetApps
func (r *IOSManagedMobileAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object IOSManagedMobileAppResourceModel

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

	policyID := object.ManagedAppProtectionId.ValueString()
	bundle := object.MobileAppIdentifier.BundleId.ValueString()

	// Serialize with sibling apps targeting the same policy (see managedAppPolicyTargetAppsLockScope).
	unlock, err := crud.AcquireKeyedLock(ctx, managedAppPolicyTargetAppsLockScope, policyID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting resource", err.Error())
		return
	}
	defer unlock()

	appsResponse, err := r.client.DeviceAppManagement().IosManagedAppProtections().ByIosManagedAppProtectionId(policyID).Apps().Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, fmt.Errorf("list policy apps: %w", err), resp, constants.TfOperationDelete, r.ReadPermissions)
		return
	}
	currentApps := appsResponse.GetValue()

	if target := findTargetedAppByBundleID(currentApps, bundle); target != nil {
		remaining := slices.DeleteFunc(slices.Clone(currentApps), func(app graphmodels.ManagedMobileAppable) bool { return app == target })

		requestBody := constructTargetAppsRequest(remaining)

		if err := r.client.DeviceAppManagement().ManagedAppPolicies().ByManagedAppPolicyId(policyID).TargetApps().Post(ctx, requestBody, nil); err != nil {
			errors.HandleKiotaGraphError(ctx, fmt.Errorf("set policy target apps: %w", err), resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}

		pollReadFailed := false
		err = crud.PollUntil(ctx, appPollInterval, func(ctx context.Context) (bool, error) {
			response, readErr := r.client.DeviceAppManagement().IosManagedAppProtections().ByIosManagedAppProtectionId(policyID).Apps().Get(ctx, nil)
			if readErr != nil {
				pollReadFailed = true
				return false, &crud.FatalPollError{Err: fmt.Errorf("list policy apps: %w", readErr)}
			}
			return findTargetedAppByBundleID(response.GetValue(), bundle) == nil, nil
		})
		if err != nil && ctx.Err() != nil && !pollReadFailed {
			err = fmt.Errorf("%w: %w", errAppWaitTimeout, ctx.Err())
		}
		if err != nil {
			resp.Diagnostics.AddError(
				"Error verifying resource deletion",
				fmt.Sprintf("Bundle %q is still targeted by policy %s: %s", bundle, policyID, err.Error()),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
