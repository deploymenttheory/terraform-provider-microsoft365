package graphBetaCrossTenantAccessPartnerSettings

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create creates a new partner-specific cross-tenant access configuration.
//
// API Calls:
//   - POST /policies/crossTenantAccessPolicy/partners
//   - GET /policies/crossTenantAccessPolicy/partners/{tenantId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicy-post-partners?view=graph-rest-beta
func (r *CrossTenantAccessPartnerSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object CrossTenantAccessPartnerSettingsResourceModel

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
		Partners().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = object.TenantID

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

// Read retrieves the current state of a partner-specific cross-tenant access configuration.
//
// API Calls:
//   - GET /policies/crossTenantAccessPolicy/partners/{tenantId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationpartner-get?view=graph-rest-beta
func (r *CrossTenantAccessPartnerSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object CrossTenantAccessPartnerSettingsResourceModel

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

	tenantID := object.TenantID.ValueString()

	remoteResource, err := r.client.
		Policies().
		CrossTenantAccessPolicy().
		Partners().
		ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).
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

// Update applies changes to a partner-specific cross-tenant access configuration.
//
// API Calls:
//   - PATCH /policies/crossTenantAccessPolicy/partners/{tenantId}
//   - GET /policies/crossTenantAccessPolicy/partners/{tenantId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationpartner-update?view=graph-rest-beta
func (r *CrossTenantAccessPartnerSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CrossTenantAccessPartnerSettingsResourceModel
	var state CrossTenantAccessPartnerSettingsResourceModel

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

	tenantID := plan.TenantID.ValueString()

	_, err = r.client.
		Policies().
		CrossTenantAccessPolicy().
		Partners().
		ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	time.Sleep(30 * time.Second)

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

// Delete removes a partner-specific cross-tenant access configuration.
//
// This operation supports both soft delete (default) and hard delete (permanent):
//   - false (default): Soft delete moves the configuration to deleted items (can be restored within 30 days)
//   - true: Hard delete permanently removes the configuration from deleted items
//
// API Calls:
//   - DELETE /policies/crossTenantAccessPolicy/partners/{tenantId} (soft delete)
//   - DELETE /directory/deletedItems/{tenantId} (if hard_delete is true)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationpartner-delete?view=graph-rest-beta
// Reference: https://learn.microsoft.com/en-us/graph/api/policydeletableitem-delete?view=graph-rest-beta
func (r *CrossTenantAccessPartnerSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CrossTenantAccessPartnerSettingsResourceModel

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

	tenantID := state.TenantID.ValueString()
	hardDelete := state.HardDelete.ValueBool()

	tflog.Info(ctx, fmt.Sprintf("Deleting partner configuration for tenant ID: %s (hard_delete: %t)", tenantID, hardDelete))

	err := r.client.
		Policies().
		CrossTenantAccessPolicy().
		Partners().
		ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	if !hardDelete {
		tflog.Info(ctx, fmt.Sprintf("Soft delete only - partner configuration %s moved to deleted items (can be restored within 30 days)", tenantID))
		resp.State.RemoveResource(ctx)
		tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Hard delete enabled - proceeding with permanent deletion of partner configuration %s", tenantID))

	err = r.client.
		Directory().
		DeletedItems().
		ByDirectoryObjectId(tenantID).
		Delete(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		if errorInfo.StatusCode == 404 || errorInfo.ErrorCode == "ResourceNotFound" || errorInfo.ErrorCode == "Request_ResourceNotFound" {
			tflog.Info(ctx, fmt.Sprintf("Partner configuration %s already permanently deleted (not found in deleted items)", tenantID))
		} else {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}
	} else {
		tflog.Info(ctx, fmt.Sprintf("Hard delete successful for partner configuration %s", tenantID))
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
