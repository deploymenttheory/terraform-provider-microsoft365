package graphBetaCrossTenantAccessPartnerUserSyncSettings

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

// Create creates a new cross-tenant access partner user sync settings configuration.
// Since this is a PUT operation (not POST), it creates or replaces the configuration.
//
// API Calls:
//   - PUT /policies/crossTenantAccessPolicy/partners/{tenant-id}/identitySynchronization
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantaccesspolicyconfigurationpartner-put-identitysynchronization?view=graph-rest-beta
func (r *CrossTenantAccessPartnerUserSyncSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object CrossTenantAccessPartnerUserSyncSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Create Method: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tenantID := object.TenantID.ValueString()

	_, err = r.client.
		Policies().
		CrossTenantAccessPolicy().
		Partners().
		ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).
		IdentitySynchronization().
		Put(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(tenantID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

// Read retrieves the current state of the cross-tenant access partner user sync settings from the Graph API.
//
// API Calls:
//   - GET /policies/crossTenantAccessPolicy/partners/{tenant-id}/identitySynchronization
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantidentitysyncpolicypartner-get?view=graph-rest-beta
func (r *CrossTenantAccessPartnerUserSyncSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object CrossTenantAccessPartnerUserSyncSettingsResourceModel

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
		IdentitySynchronization().
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

// Update updates an existing cross-tenant access partner user sync settings configuration.
//
// API Calls:
//   - PATCH /policies/crossTenantAccessPolicy/partners/{tenant-id}/identitySynchronization
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantidentitysyncpolicypartner-update?view=graph-rest-beta
func (r *CrossTenantAccessPartnerUserSyncSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object CrossTenantAccessPartnerUserSyncSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update Method: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tenantID := object.TenantID.ValueString()

	_, err = r.client.
		Policies().
		CrossTenantAccessPolicy().
		Partners().
		ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).
		IdentitySynchronization().
		Put(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

// Delete removes the cross-tenant access partner user sync settings configuration.
//
// API Calls:
//   - DELETE /policies/crossTenantAccessPolicy/partners/{tenant-id}/identitySynchronization
//
// Reference: https://learn.microsoft.com/en-us/graph/api/crosstenantidentitysyncpolicypartner-delete?view=graph-rest-beta
func (r *CrossTenantAccessPartnerUserSyncSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object CrossTenantAccessPartnerUserSyncSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete Method: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	tenantID := object.TenantID.ValueString()

	err := r.client.
		Policies().
		CrossTenantAccessPolicy().
		Partners().
		ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).
		IdentitySynchronization().
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
