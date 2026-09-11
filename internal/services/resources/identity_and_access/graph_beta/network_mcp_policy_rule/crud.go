package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func (r *NetworkMCPPolicyRuleResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var object NetworkMCPPolicyRuleResourceModel
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
	body, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing MCP resource", err.Error())
		return
	}
	created, err := r.createMCPPolicyRule(
		ctx,
		object.MCPPolicyID.ValueString(),
		body,
	)
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	if created == nil || created.ID == nil || *created.ID == "" {
		resp.Diagnostics.AddError(
			"Invalid create response",
			"The API returned no resource id. Check Graph before retrying creation.",
		)
		return
	}
	// Persist the response body identity before readback. The observed Location header is only the Graph root.
	object.ID = types.StringValue(*created.ID)
	object.Status = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				MCPPolicyRuleIdentity{
					ID:          object.ID.ValueString(),
					MCPPolicyID: object.MCPPolicyID.ValueString(),
				},
			)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	desired := object
	remote, err := r.getMCPPolicyRule(
		ctx,
		object.MCPPolicyID.ValueString(),
		object.ID.ValueString(),
	)
	// This is a create readback, not a refresh: even a 404 must retain the created id.
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid MCP resource response", err.Error())
		return
	}
	if diff, err := constructUpdateResource(
		ctx,
		&desired,
		&object,
	); err != nil ||
		diff.hasChanges() {
		resp.Diagnostics.AddError(
			"MCP readback does not match configuration",
			"The API returned success but did not persist all requested attributes. Check the resource before retrying.",
		)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkMCPPolicyRuleResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var object NetworkMCPPolicyRuleResourceModel
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
	remote, err := r.getMCPPolicyRule(
		ctx,
		object.MCPPolicyID.ValueString(),
		object.ID.ValueString(),
	)
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(
			ctx,
			err,
			resp,
			constants.TfOperationRead,
			r.ReadPermissions,
			errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true},
		)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid MCP resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				MCPPolicyRuleIdentity{
					ID:          object.ID.ValueString(),
					MCPPolicyID: object.MCPPolicyID.ValueString(),
				},
			)...)
	}
}

func (r *NetworkMCPPolicyRuleResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.State = req.State
	var plan, state NetworkMCPPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructUpdateResource(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing MCP update", err.Error())
		return
	}
	if body.hasChanges() {
		if err := r.updateMCPPolicyRule(
			ctx,
			state.MCPPolicyID.ValueString(),
			state.ID.ValueString(),
			body,
		); err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationUpdate,
				r.WritePermissions,
			)
			return
		}
	}
	// Keep the previous known state if readback fails after a successful PATCH.
	resp.State = req.State
	object := plan
	object.ID = state.ID
	remote, err := r.getMCPPolicyRule(
		ctx,
		object.MCPPolicyID.ValueString(),
		object.ID.ValueString(),
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid MCP resource response", err.Error())
		return
	}
	if diff, err := constructUpdateResource(ctx, &plan, &object); err != nil || diff.hasChanges() {
		resp.Diagnostics.AddError(
			"MCP readback does not match configuration",
			"The API returned success but did not persist all requested attributes. Check the resource before retrying.",
		)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				MCPPolicyRuleIdentity{
					ID:          object.ID.ValueString(),
					MCPPolicyID: object.MCPPolicyID.ValueString(),
				},
			)...)
	}
}

func (r *NetworkMCPPolicyRuleResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var object NetworkMCPPolicyRuleResourceModel
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
	if err := r.deleteMCPPolicyRule(
		ctx,
		object.MCPPolicyID.ValueString(),
		object.ID.ValueString(),
	); err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationDelete,
			r.WritePermissions,
		)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.State.RemoveResource(ctx)
}
