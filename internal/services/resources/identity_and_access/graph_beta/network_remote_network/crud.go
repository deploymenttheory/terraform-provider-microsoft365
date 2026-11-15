package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// Create handles the Create operation.
//
// API Calls:
//   - POST /networkaccess/connectivity/remoteNetworks
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-connectivity-post-remotenetworks?view=graph-rest-beta
func (r *NetworkRemoteNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkRemoteNetworkResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	profileIDs, diags := setStrings(ctx, plan.ForwardingProfileIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.create(ctx, plan.Name.ValueString(), plan.Region.ValueString(), profileIDs)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// Keep the ID in state even if the read-back fails, so the object can still be destroyed.
	plan.ID = types.StringValue(*created.ID)
	plan.Version = types.StringPointerValue(created.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		readbackOptions(plan, constants.TfOperationCreate),
	); err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation.
//
// API Calls:
//   - GET /networkaccess/connectivity/remoteNetworks/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-remotenetwork-get?view=graph-rest-beta
func (r *NetworkRemoteNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkRemoteNetworkResourceModel

	// During post-write read-back a 404 is replica lag and must be retried, not treated as a deletion.
	operation := constants.TfOperationRead
	if ctxOp, ok := ctx.Value("retry_operation").(string); ok && ctxOp != "" {
		operation = ctxOp
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	remote, err := r.get(ctx, state.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if err := mapResponse(&state, remote); err != nil {
		resp.Diagnostics.AddError("Invalid Graph response", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles the Update operation.
//
// API Calls:
//   - PATCH /networkaccess/connectivity/remoteNetworks/{id} (name)
//   - PATCH /networkaccess/connectivity/remoteNetworks/{id}/forwardingProfiles ($delta, full desired set)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-remotenetwork-update?view=graph-rest-beta
func (r *NetworkRemoteNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkRemoteNetworkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	id := state.ID.ValueString()

	if !plan.Name.Equal(state.Name) {
		if err := r.updateName(ctx, id, plan.Name.ValueString()); err != nil {
			r.handleUpdateError(ctx, req, resp, err)
			return
		}
	}

	if !plan.ForwardingProfileIDs.Equal(state.ForwardingProfileIDs) {
		profileIDs, diags := setStrings(ctx, plan.ForwardingProfileIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if err := r.updateForwardingProfiles(ctx, id, profileIDs); err != nil {
			r.handleUpdateError(ctx, req, resp, err)
			return
		}
	}

	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		readbackOptions(plan, constants.TfOperationUpdate),
	); err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, id))
}

// handleUpdateError reports a failed write and keeps the prior state, so the next plan retries
// the change instead of trusting the planned values that were never applied.
func (r *NetworkRemoteNetworkResource) handleUpdateError(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, err error) {
	errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
	resp.State.Raw = req.State.Raw.Copy()
}

// Delete handles the Delete operation.
//
// API Calls:
//   - DELETE /networkaccess/connectivity/remoteNetworks/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-connectivity-delete-remotenetworks?view=graph-rest-beta
func (r *NetworkRemoteNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkRemoteNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if err := r.delete(ctx, state.ID.ValueString()); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)
}

func readbackOptions(plan NetworkRemoteNetworkResourceModel, operation string) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = stableConsistencyPredicate(plan)
	return opts
}

func setStrings(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	var values []string
	if set.IsNull() || set.IsUnknown() {
		return values, nil
	}
	diags := set.ElementsAs(ctx, &values, false)
	return values, diags
}
