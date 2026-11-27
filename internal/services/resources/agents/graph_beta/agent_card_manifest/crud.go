package graphBetaAgentCardManifest

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

var readbackMatch = crud.AttributeMatchOptions{
	Paths: crud.RootPaths(
		"display_name",
		"description",
		"icon_url",
		"protocol_version",
		"version",
		"documentation_url",
		"managed_by",
		"originating_store",
		"owner_ids",
		"default_input_modes",
		"default_output_modes",
		"supports_authenticated_extended_card",
		"agent_provider",
		"capabilities",
		"skills",
		"security_schemes_json",
		"security_json",
	),
	RequireKnown:       crud.RootPaths("id"),
	IgnoreNullExpected: true,
}

// Create adopts and patches an existing manifest. Graph has no standalone POST.
func (r *AgentCardManifestResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object, configured AgentCardManifestResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &configured)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	existing, err := r.client.AgentRegistry().AgentCardManifests().ByAgentCardManifestId(object.AgentCardManifestID.ValueString()).Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	if existing == nil || existing.GetId() == nil {
		resp.Diagnostics.AddError("Manifest not found", "Graph returned no existing manifest identity.")
		return
	}
	object.ID = types.StringValue(*existing.GetId())
	resp.Diagnostics.Append(crud.SetPartialState(ctx, &resp.State, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body, err := constructResource(ctx, &configured)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing manifest", err.Error())
		return
	}
	_, err = r.client.AgentRegistry().AgentCardManifests().ByAgentCardManifestId(object.ID.ValueString()).Patch(ctx, body, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	container := &crud.CreateResponseContainer{CreateResponse: resp}
	opts := crud.NewReadWithRetryOptions(
		constants.TfOperationCreate,
		ResourceName,
		crud.ReadWithConsistencyPredicate(crud.RequireStableReads(crud.AttributesMatchPredicate(req.Config, readbackMatch))),
	)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, container, opts); err != nil {
		resp.Diagnostics.AddError("Error reading adopted manifest", fmt.Sprintf("Read-back failed: %s", err))
	}
}

// Read reads the manifest by ID. A genuine 404 removes it from state.
func (r *AgentCardManifestResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AgentCardManifestResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	operation := constants.TfOperationRead
	if value, ok := ctx.Value("retry_operation").(string); ok {
		operation = value
	}
	remote, err := r.client.AgentRegistry().AgentCardManifests().ByAgentCardManifestId(object.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	if err := MapRemoteResourceStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Error mapping manifest", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
	}
}

// Update patches configured properties only; omitted values remain remote-owned.
func (r *AgentCardManifestResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state, configured AgentCardManifestResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &configured)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(ctx, &configured)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing manifest", err.Error())
		return
	}
	_, err = r.client.AgentRegistry().AgentCardManifests().ByAgentCardManifestId(state.ID.ValueString()).Patch(ctx, body, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(crud.SetPartialState(ctx, &resp.State, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	container := &crud.UpdateResponseContainer{UpdateResponse: resp}
	opts := crud.NewReadWithRetryOptions(
		constants.TfOperationUpdate,
		ResourceName,
		crud.ReadWithConsistencyPredicate(crud.RequireStableReads(crud.AttributesMatchPredicate(req.Config, readbackMatch))),
	)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, container, opts); err != nil {
		resp.Diagnostics.AddError("Error reading updated manifest", fmt.Sprintf("Read-back failed: %s", err))
	}
}

// Delete forgets management. This resource owns neither the manifest's lifetime
// nor its parent agent instance, and Graph exposes no standalone manifest DELETE.
func (r *AgentCardManifestResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

// ModifyPlan retains metadata only when no writable attribute changes.
func (r *AgentCardManifestResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	planmodifiers.PreserveComputedStringsWhenUnchanged(ctx, req, resp, "created_by", "created_date_time", "last_modified_date_time")
}
