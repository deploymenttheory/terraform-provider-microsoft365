package graphBetaDeviceAndAppManagementIOSManagedMobileApp

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create adds the app to the policy's targeted apps.
//
// API Calls:
//   - GET  /deviceAppManagement/iosManagedAppProtections/{id}/apps
//   - POST /deviceAppManagement/managedAppPolicies/{id}/targetApps
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-mam-targetedmanagedappprotection-targetapps?view=graph-rest-beta
func (r *IOSManagedMobileAppResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var object IOSManagedMobileAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
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
	pkg := bundleID(newApp)

	unlock, err := lockPolicy(ctx, policyID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
		return
	}
	defer unlock()

	currentApps, err := r.listPolicyApps(ctx, policyID)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}

	if existing := findByBundle(currentApps, pkg); existing != nil {
		existingID := pkg + ".ios"
		if existing.GetId() != nil {
			existingID = *existing.GetId()
		}
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf(
				"%s: bundle %q in policy %s. Import it with ID \"%s/%s\" instead.",
				errAppAlreadyTargeted, pkg, policyID, policyID, existingID,
			),
		)
		return
	}

	desired := make([]graphmodels.ManagedMobileAppable, 0, len(currentApps)+1)
	for _, app := range currentApps {
		desired = append(desired, writableApp(app))
	}
	desired = append(desired, newApp)

	if err := r.setPolicyApps(ctx, policyID, desired); err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}

	// targetApps returns no content; Intune assigns the id, so wait for the entry.
	apps, err := r.waitForApps(ctx, policyID, func(apps []graphmodels.ManagedMobileAppable) bool {
		created := findByBundle(apps, pkg)
		return created != nil && created.GetId() != nil
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not find bundle %q in policy %s: %s", pkg, policyID, err.Error()),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, findByBundle(apps, pkg))

	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...,
		)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
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
func (r *IOSManagedMobileAppResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var object IOSManagedMobileAppResourceModel
	var identity sharedmodels.ResourceIdentity

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	apps, err := r.listPolicyApps(ctx, object.ManagedAppProtectionId.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions)
		return
	}

	app := findByID(apps, object.ID.ValueString())
	if app == nil {
		tflog.Warn(ctx, fmt.Sprintf(
			"%s with ID %s is no longer targeted by policy %s, removing from state",
			ResourceName, object.ID.ValueString(), object.ManagedAppProtectionId.ValueString(),
		))
		resp.State.RemoveResource(ctx)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, app)

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
// to Graph forces replacement, and Intune assigns version itself.
func (r *IOSManagedMobileAppResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan IOSManagedMobileAppResourceModel
	var state IOSManagedMobileAppResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = state.ID
	plan.Version = state.Version

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete removes the app from the policy's targeted apps.
//
// API Calls:
//   - GET  /deviceAppManagement/iosManagedAppProtections/{id}/apps
//   - POST /deviceAppManagement/managedAppPolicies/{id}/targetApps
func (r *IOSManagedMobileAppResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var object IOSManagedMobileAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Delete,
		DeleteTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	policyID := object.ManagedAppProtectionId.ValueString()
	appID := object.ID.ValueString()

	unlock, err := lockPolicy(ctx, policyID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting resource", err.Error())
		return
	}
	defer unlock()

	currentApps, err := r.listPolicyApps(ctx, policyID)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		return
	}

	if findByID(currentApps, appID) != nil {
		remaining := make([]graphmodels.ManagedMobileAppable, 0, len(currentApps))
		for _, app := range currentApps {
			if app.GetId() == nil || *app.GetId() != appID {
				remaining = append(remaining, writableApp(app))
			}
		}

		if err := r.setPolicyApps(ctx, policyID, remaining); err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationDelete,
				r.WritePermissions,
			)
			return
		}

		_, err = r.waitForApps(ctx, policyID, func(apps []graphmodels.ManagedMobileAppable) bool {
			return findByID(apps, appID) == nil
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error verifying resource deletion",
				fmt.Sprintf(
					"App %s is still targeted by policy %s: %s",
					appID,
					policyID,
					err.Error(),
				),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
