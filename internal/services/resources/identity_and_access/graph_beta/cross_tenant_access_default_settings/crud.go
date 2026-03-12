package graphBetaCrossTenantAccessDefaultSettings

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
)

// Create configures the cross-tenant access default settings.
//
// The crossTenantAccessPolicyConfigurationDefault resource is a singleton that always exists in the tenant; there is no
// POST endpoint to create it. This operation therefore uses a PATCH request to apply the desired
// configuration, mirroring the Update behaviour.
//
// API Calls:
//   - PATCH /policies/crossTenantAccessPolicy/default
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationdefault-update?view=graph-rest-beta
func (r *CrossTenantAccessDefaultSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object CrossTenantAccessDefaultSettingsResourceModel

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
			"Validation failed for Create Method",
			fmt.Sprintf("Pre-request validation failed for %s: %s", ResourceName, err.Error()),
		)
		return
	}

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Policies().
		CrossTenantAccessPolicy().
		DefaultEscaped().
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// This is a singleton resource — set the static ID immediately so the subsequent Read
	// can locate the correct state entry.
	object.ID = types.StringValue(singletonID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// required for eventual consistency
	time.Sleep(20 * time.Second)

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

// Read retrieves the current state of the cross-tenant access default settings from the Graph API.
//
// API Calls:
//   - GET /policies/crossTenantAccessPolicy/default
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationdefault-get?view=graph-rest-beta
func (r *CrossTenantAccessDefaultSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object CrossTenantAccessDefaultSettingsResourceModel

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

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	remoteResource, err := r.client.
		Policies().
		CrossTenantAccessPolicy().
		DefaultEscaped().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, remoteResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update applies changes to the cross-tenant access default settings.
//
// API Calls:
//   - PATCH /policies/crossTenantAccessPolicy/default
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationdefault-update?view=graph-rest-beta
func (r *CrossTenantAccessDefaultSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CrossTenantAccessDefaultSettingsResourceModel
	var state CrossTenantAccessDefaultSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s", ResourceName))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if err := validateRequest(ctx, r.client, &plan); err != nil {
		resp.Diagnostics.AddError(
			"Validation failed for Update Method",
			fmt.Sprintf("Pre-request validation failed for %s: %s", ResourceName, err.Error()),
		)
		return
	}

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Policies().
		CrossTenantAccessPolicy().
		DefaultEscaped().
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// required for eventual consistency
	time.Sleep(20 * time.Second)

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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the removal of this resource from Terraform management.
//
// Because the crossTenantAccessPolicyConfigurationDefault is a singleton that always exists in the tenant and has no
// DELETE API endpoint, this operation has two behaviours controlled by restore_defaults_on_destroy:
//
//   - false (default): Terraform removes the resource from state only. The existing default configuration
//     in Microsoft Entra ID is left unchanged.
//
//   - true: Terraform issues a POST request to resetToSystemDefault to reset the default configuration to
//     system defaults, then verifies that is_service_default is true, and finally removes it from state.
//
// API Calls (when restore_defaults_on_destroy = true):
//   - POST /policies/crossTenantAccessPolicy/default/resetToSystemDefault
//   - GET /policies/crossTenantAccessPolicy/default (to verify reset)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationdefault-resettosystemdefault?view=graph-rest-beta
func (r *CrossTenantAccessDefaultSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CrossTenantAccessDefaultSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if state.RestoreDefaultsOnDestroy.ValueBool() {
		tflog.Debug(ctx, fmt.Sprintf("restore_defaults_on_destroy is true — resetting %s to system defaults", ResourceName))

		err := r.client.
			Policies().
			CrossTenantAccessPolicy().
			DefaultEscaped().
			ResetToSystemDefault().
			Post(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully called resetToSystemDefault for %s, now verifying is_service_default is true", ResourceName))

		// Wait for eventual consistency
		time.Sleep(20 * time.Second)

		// Verify that the reset was successful by checking is_service_default
		remoteResource, err := r.client.
			Policies().
			CrossTenantAccessPolicy().
			DefaultEscaped().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
			return
		}

		if remoteResource.GetIsServiceDefault() == nil || !*remoteResource.GetIsServiceDefault() {
			resp.Diagnostics.AddError(
				"Failed to verify cross tenant access policy configuration default was reset to system defaults",
				fmt.Sprintf("After calling resetToSystemDefault, is_service_default is not true for %s", ResourceName),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully verified %s is reset to system defaults (is_service_default = true)", ResourceName))
	} else {
		tflog.Debug(ctx, fmt.Sprintf("restore_defaults_on_destroy is false — removing %s from Terraform state only, default configuration unchanged", ResourceName))
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
