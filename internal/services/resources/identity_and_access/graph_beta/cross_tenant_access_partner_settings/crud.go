package graphBetaCrossTenantAccessPartnerSettings

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// readbackAttributes are the attributes readbackMatch compares.
var readbackAttributes = []string{
	"tenant_id", "b2b_collaboration_inbound", "b2b_collaboration_outbound", "b2b_direct_connect_inbound", "b2b_direct_connect_outbound", "inbound_trust",
	"automatic_user_consent_settings", "tenant_restrictions",
}

// readbackMatch is what the read-back after a write waits for: a populated id, the partner tenant and the written settings. With enforce_default_inheritance, an unset setting must read back null.
var readbackMatch = crud.AttributeMatchOptions{
	Paths:        crud.RootPaths(readbackAttributes...),
	RequireKnown: crud.RootPaths("id"),
}

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

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.RetryInterval = 5 * time.Second
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	// Without enforce_default_inheritance an unset setting keeps its remote value.
	match := readbackMatch
	match.IgnoreNullExpected = !object.EnforceDefaultInheritance.ValueBool()
	opts.ConsistencyPredicate = crud.RequireStableReads(crud.AttributesMatchPredicate(resp.State, match))

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

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.RetryInterval = 5 * time.Second
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	// Without enforce_default_inheritance an unset setting keeps its remote value.
	match := readbackMatch
	match.IgnoreNullExpected = !plan.EnforceDefaultInheritance.ValueBool()
	opts.ConsistencyPredicate = crud.RequireStableReads(crud.AttributesMatchPredicate(resp.State, match))

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

// Delete soft-deletes the active partner and, when requested, permanently deletes
// it from the policy recycle bin. Directory deleted items are a different store.
//
// API Calls:
//   - DELETE /policies/crossTenantAccessPolicy/partners/{tenantId}
//   - GET /policies/crossTenantAccessPolicy/partners/{tenantId}
//   - GET /policies/deletedItems/crossTenantPartners/{tenantId}
//   - DELETE /policies/deletedItems/crossTenantPartners/{tenantId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/policydeletableitem-delete?view=graph-rest-beta
func (r *CrossTenantAccessPartnerSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CrossTenantAccessPartnerSettingsResourceModel
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
	err := r.client.Policies().CrossTenantAccessPolicy().Partners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).Delete(ctx, nil)
	softDeleted := err == nil
	if err != nil && errors.GraphError(ctx, err).StatusCode != 404 {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	// A retry can start with the partner already soft-deleted. Do not discard state
	// on that active-store 404 before checking the policy recycle bin.
	const deleteMaxRetries = 30
	const deleteRetryInterval = 5 * time.Second
	for attempt := 1; attempt <= deleteMaxRetries; attempt++ {
		_, getErr := r.client.Policies().CrossTenantAccessPolicy().Partners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).Get(ctx, nil)
		if getErr != nil {
			if errors.GraphError(ctx, getErr).StatusCode != 404 {
				errors.HandleKiotaGraphError(ctx, getErr, resp, constants.TfOperationDelete, r.ReadPermissions)
				return
			}
			if !hardDelete {
				resp.State.RemoveResource(ctx)
				return
			}
			_, deletedErr := r.client.Policies().DeletedItems().CrossTenantPartners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).Get(ctx, nil)
			if deletedErr == nil {
				break
			}
			if errors.GraphError(ctx, deletedErr).StatusCode != 404 {
				errors.HandleKiotaGraphError(ctx, deletedErr, resp, constants.TfOperationDelete, r.ReadPermissions)
				return
			}
			if !softDeleted {
				// Both stores confirm absence, including a repeated destroy after a prior purge.
				resp.State.RemoveResource(ctx)
				return
			}
			// The policy recycle bin may lag behind the active-store deletion. A 404
			// immediately after our soft delete is not evidence of permanent removal.
		}
		if attempt == deleteMaxRetries {
			resp.Diagnostics.AddError(
				"Partner deletion did not propagate",
				"Could not confirm the expected partner deletion state. Terraform state is retained so deletion can be retried.",
			)
			return
		}
		select {
		case <-time.After(deleteRetryInterval):
		case <-ctx.Done():
			resp.Diagnostics.AddError("Context cancelled while waiting for partner deletion", ctx.Err().Error())
			return
		}
	}
	err = r.client.Policies().DeletedItems().CrossTenantPartners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).Delete(ctx, nil)
	if err != nil && errors.GraphError(ctx, err).StatusCode != 404 {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	// Permanent deletion also propagates. Confirm absence in the same policy
	// store before forgetting the partner; a successful DELETE can precede its GET.
	absentReads := 0
	const requiredAbsentReads = 3
	for attempt := 1; attempt <= deleteMaxRetries; attempt++ {
		_, getErr := r.client.Policies().DeletedItems().CrossTenantPartners().ByCrossTenantAccessPolicyConfigurationPartnerTenantId(tenantID).Get(ctx, nil)
		if getErr != nil {
			if errors.GraphError(ctx, getErr).StatusCode != 404 {
				errors.HandleKiotaGraphError(ctx, getErr, resp, constants.TfOperationDelete, r.ReadPermissions)
				return
			}
			absentReads++
			if absentReads >= requiredAbsentReads {
				resp.State.RemoveResource(ctx)
				return
			}
		} else {
			absentReads = 0
		}
		if attempt == deleteMaxRetries {
			resp.Diagnostics.AddError(
				"Permanent partner deletion did not propagate",
				"The partner remains visible in policy deleted items. Terraform state is retained so deletion can be retried.",
			)
			return
		}
		select {
		case <-time.After(deleteRetryInterval):
		case <-ctx.Done():
			resp.Diagnostics.AddError("Context cancelled while waiting for permanent partner deletion", ctx.Err().Error())
			return
		}
	}
}
