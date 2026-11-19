package graphBetaNetworkMCPPolicyRule

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

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
	createAdapter := r.client.GetAdapter()
	createRequest, err := newMCPPolicyRuleRequestInformation(
		ctx,
		createAdapter,
		abstractions.POST,
		object.MCPPolicyID.ValueString(),
		"",
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
	createResult, err := createAdapter.SendPrimitive(ctx, createRequest, "[]byte", mcpPolicyRuleErrorMapping)
	created, err := decodeSendMCPPolicyRuleResponse(createResult, err)
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
	readbackAdapter := r.client.GetAdapter()
	readbackRequest, err := newMCPPolicyRuleRequestInformation(
		ctx,
		readbackAdapter,
		abstractions.GET,
		object.MCPPolicyID.ValueString(),
		object.ID.ValueString(),
		nil,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	readbackResult, err := readbackAdapter.SendPrimitive(ctx, readbackRequest, "[]byte", mcpPolicyRuleErrorMapping)
	remote, err := decodeSendMCPPolicyRuleResponse(readbackResult, err)
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
	readAdapter := r.client.GetAdapter()
	readRequest, err := newMCPPolicyRuleRequestInformation(
		ctx,
		readAdapter,
		abstractions.GET,
		object.MCPPolicyID.ValueString(),
		object.ID.ValueString(),
		nil,
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
	readResult, err := readAdapter.SendPrimitive(ctx, readRequest, "[]byte", mcpPolicyRuleErrorMapping)
	remote, err := decodeSendMCPPolicyRuleResponse(readResult, err)
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
		updateAdapter := r.client.GetAdapter()
		updateRequest, err := newMCPPolicyRuleRequestInformation(
			ctx,
			updateAdapter,
			abstractions.PATCH,
			state.MCPPolicyID.ValueString(),
			state.ID.ValueString(),
			body,
		)
		if err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationUpdate,
				r.WritePermissions,
			)
			return
		}
		updateResult, err := updateAdapter.SendPrimitive(ctx, updateRequest, "[]byte", mcpPolicyRuleErrorMapping)
		if err := decodeUpdateMCPPolicyRuleResponse(updateResult, err); err != nil {
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
	readbackAdapter := r.client.GetAdapter()
	readbackRequest, err := newMCPPolicyRuleRequestInformation(
		ctx,
		readbackAdapter,
		abstractions.GET,
		object.MCPPolicyID.ValueString(),
		object.ID.ValueString(),
		nil,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}
	readbackResult, err := readbackAdapter.SendPrimitive(ctx, readbackRequest, "[]byte", mcpPolicyRuleErrorMapping)
	remote, err := decodeSendMCPPolicyRuleResponse(readbackResult, err)
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
	deleteAdapter := r.client.GetAdapter()
	deleteRequest, err := newMCPPolicyRuleRequestInformation(
		ctx,
		deleteAdapter,
		abstractions.DELETE,
		object.MCPPolicyID.ValueString(),
		object.ID.ValueString(),
		nil,
	)
	if err != nil {
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
	if err := deleteAdapter.SendNoContent(ctx, deleteRequest, mcpPolicyRuleErrorMapping); err != nil {
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

var mcpPolicyRuleErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}
