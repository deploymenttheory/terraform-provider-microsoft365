package graphBetaNetworkMCPPolicy

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
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

func (r *NetworkMCPPolicyResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var object NetworkMCPPolicyResourceModel
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
	createRequest, err := newMCPPolicyRequestInformation(
		ctx,
		createAdapter,
		abstractions.POST,
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
	createResult, err := createAdapter.SendPrimitive(ctx, createRequest, "[]byte", mcpPolicyErrorMapping)
	created, err := decodeCreateMCPPolicyResponse(createResult, err)
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
	object.Version = types.StringNull()
	object.LastModifiedDateTime = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{
			ID: object.ID.ValueString(),
		})...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	desired := object
	readbackAdapter := r.client.GetAdapter()
	readbackRequest, err := newMCPPolicyRequestInformation(
		ctx,
		readbackAdapter,
		abstractions.GET,
		object.ID.ValueString(),
		nil,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	readbackResult, err := readbackAdapter.SendPrimitive(ctx, readbackRequest, "[]byte", mcpPolicyErrorMapping)
	remote, err := decodeSendMCPPolicyResponse(readbackResult, err)
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

func (r *NetworkMCPPolicyResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var object NetworkMCPPolicyResourceModel
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
	readRequest, err := newMCPPolicyRequestInformation(
		ctx,
		readAdapter,
		abstractions.GET,
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
	readResult, err := readAdapter.SendPrimitive(ctx, readRequest, "[]byte", mcpPolicyErrorMapping)
	remote, err := decodeSendMCPPolicyResponse(readResult, err)
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
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{
			ID: object.ID.ValueString(),
		})...)
	}
}

func (r *NetworkMCPPolicyResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.State = req.State
	var plan, state NetworkMCPPolicyResourceModel
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
		updateRequest, err := newMCPPolicyRequestInformation(
			ctx,
			updateAdapter,
			abstractions.PATCH,
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
		updateResult, err := updateAdapter.SendPrimitive(ctx, updateRequest, "[]byte", mcpPolicyErrorMapping)
		if err := decodeUpdateMCPPolicyResponse(updateResult, err); err != nil {
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
	readbackRequest, err := newMCPPolicyRequestInformation(
		ctx,
		readbackAdapter,
		abstractions.GET,
		object.ID.ValueString(),
		nil,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}
	readbackResult, err := readbackAdapter.SendPrimitive(ctx, readbackRequest, "[]byte", mcpPolicyErrorMapping)
	remote, err := decodeSendMCPPolicyResponse(readbackResult, err)
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
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{
			ID: object.ID.ValueString(),
		})...)
	}
}

func (r *NetworkMCPPolicyResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var object NetworkMCPPolicyResourceModel
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
	deleteRequest, err := newMCPPolicyRequestInformation(
		ctx,
		deleteAdapter,
		abstractions.DELETE,
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
	if err := deleteAdapter.SendNoContent(ctx, deleteRequest, mcpPolicyErrorMapping); err != nil {
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

var mcpPolicyErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}
