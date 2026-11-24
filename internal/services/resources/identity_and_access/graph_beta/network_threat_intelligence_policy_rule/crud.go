package graphBetaNetworkThreatIntelligencePolicyRule

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	kiotahttp "github.com/microsoft/kiota-http-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func (r *NetworkThreatIntelligencePolicyRuleResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var object NetworkThreatIntelligencePolicyRuleResourceModel
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
	body, err := constructResource(ctx, &object, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing threat intelligence resource", err.Error())
		return
	}
	// Kiota 1.5.6 can retain a gzip header while replaying the original JSON body, so retries send it uncompressed.
	options := []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)}
	createConfig := &networkaccess.ThreatIntelligencePoliciesItemPolicyRulesRequestBuilderPostRequestConfiguration{Options: options}
	rules := r.client.NetworkAccess().ThreatIntelligencePolicies().ByThreatIntelligencePolicyId(object.ThreatIntelligencePolicyID.ValueString()).PolicyRules()
	created, err := rules.Post(ctx, body, createConfig)
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
	if created == nil || created.GetId() == nil || *created.GetId() == "" {
		resp.Diagnostics.AddError(
			"Invalid create response",
			"The API returned no resource id. Check Graph before retrying creation.",
		)
		return
	}
	// Persist the acquired identity before readback so a failed GET cannot cause duplicate creation.
	object.ID = types.StringValue(*created.GetId())
	object.Status = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				ThreatIntelligencePolicyRuleIdentity{
					ID:                         object.ID.ValueString(),
					ThreatIntelligencePolicyID: object.ThreatIntelligencePolicyID.ValueString(),
				},
			)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := rules.ByPolicyRuleId(object.ID.ValueString()).Get(ctx, nil)
	// This is a create readback, not a refresh: even a 404 must retain the created id.
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		if stderrors.Is(err, errUnknownRuleStatus) {
			// MapRemoteStateToTerraform records the observed status before reporting it as unsupported.
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("status"), object.Status)...)
		}
		resp.Diagnostics.AddError("Invalid threat intelligence resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkThreatIntelligencePolicyRuleResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var object NetworkThreatIntelligencePolicyRuleResourceModel
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
	rules := r.client.NetworkAccess().ThreatIntelligencePolicies().ByThreatIntelligencePolicyId(object.ThreatIntelligencePolicyID.ValueString()).PolicyRules()
	remote, err := rules.ByPolicyRuleId(object.ID.ValueString()).Get(ctx, nil)
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
		if stderrors.Is(err, errUnknownRuleStatus) {
			// MapRemoteStateToTerraform records the observed status before reporting it as unsupported.
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("status"), object.Status)...)
		}
		resp.Diagnostics.AddError("Invalid threat intelligence resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				ThreatIntelligencePolicyRuleIdentity{
					ID:                         object.ID.ValueString(),
					ThreatIntelligencePolicyID: object.ThreatIntelligencePolicyID.ValueString(),
				},
			)...)
	}
}

func (r *NetworkThreatIntelligencePolicyRuleResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkThreatIntelligencePolicyRuleResourceModel
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
	body, err := constructResource(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing threat intelligence update", err.Error())
		return
	}
	changed := !plan.Name.Equal(state.Name) || !plan.Description.Equal(state.Description) || !plan.Priority.Equal(state.Priority) ||
		!plan.Action.Equal(state.Action) || !plan.Enabled.Equal(state.Enabled) || !plan.Severity.Equal(state.Severity) ||
		!plan.Destinations.Equal(state.Destinations)
	if changed {
		// Kiota 1.5.6 can retain a gzip header while replaying the original JSON body, so retries send it uncompressed.
		options := []abstractions.RequestOption{kiotahttp.NewCompressionOptionsReference(false)}
		patchConfig := &networkaccess.ThreatIntelligencePoliciesItemPolicyRulesPolicyRuleItemRequestBuilderPatchRequestConfiguration{Options: options}
		rules := r.client.NetworkAccess().ThreatIntelligencePolicies().ByThreatIntelligencePolicyId(state.ThreatIntelligencePolicyID.ValueString()).PolicyRules()
		if _, err := rules.ByPolicyRuleId(state.ID.ValueString()).Patch(ctx, body, patchConfig); err != nil {
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
	rules := r.client.NetworkAccess().ThreatIntelligencePolicies().ByThreatIntelligencePolicyId(object.ThreatIntelligencePolicyID.ValueString()).PolicyRules()
	remote, err := rules.ByPolicyRuleId(object.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		if stderrors.Is(err, errUnknownRuleStatus) {
			// MapRemoteStateToTerraform records the observed status before reporting it as unsupported.
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("status"), object.Status)...)
		}
		resp.Diagnostics.AddError("Invalid threat intelligence resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				ThreatIntelligencePolicyRuleIdentity{
					ID:                         object.ID.ValueString(),
					ThreatIntelligencePolicyID: object.ThreatIntelligencePolicyID.ValueString(),
				},
			)...)
	}
}

func (r *NetworkThreatIntelligencePolicyRuleResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var object NetworkThreatIntelligencePolicyRuleResourceModel
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
	rules := r.client.NetworkAccess().ThreatIntelligencePolicies().ByThreatIntelligencePolicyId(object.ThreatIntelligencePolicyID.ValueString()).PolicyRules()
	if err := rules.ByPolicyRuleId(object.ID.ValueString()).Delete(ctx, nil); err != nil {
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
