package graphBetaNetworkCloudFirewallPolicyRule

import (
	"context"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *NetworkCloudFirewallPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkCloudFirewallPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	created, err := r.createRule(ctx, &object)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created == nil || created.GetId() == nil || *created.GetId() == "" {
		resp.Diagnostics.AddError("Invalid cloud firewall create response", "The API did not return a resource ID.")
		return
	}
	object.ID = types.StringValue(*created.GetId())
	object.Status = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, ruleIdentity{PolicyID: object.PolicyID.ValueString(), ID: object.ID.ValueString()})...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(object.PolicyID.ValueString()).PolicyRules().ByPolicyRuleId(object.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid cloud firewall readback", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkCloudFirewallPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkCloudFirewallPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	remote, err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(object.PolicyID.ValueString()).PolicyRules().ByPolicyRuleId(object.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions, errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid cloud firewall response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, ruleIdentity{PolicyID: object.PolicyID.ValueString(), ID: object.ID.ValueString()})...)
	}
}

func (r *NetworkCloudFirewallPolicyRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkCloudFirewallPolicyRuleResourceModel
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
	body, err := constructRule(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing cloud firewall rule", err.Error())
		return
	}
	if body.changed {
		if _, err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(state.PolicyID.ValueString()).PolicyRules().ByPolicyRuleId(state.ID.ValueString()).Patch(ctx, body, nil); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}
	// Keep the prior known state if readback fails after a successful update.
	object := state
	object.Timeouts = plan.Timeouts
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	remote, err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(object.PolicyID.ValueString()).PolicyRules().ByPolicyRuleId(object.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid cloud firewall readback", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkCloudFirewallPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkCloudFirewallPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	if err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(object.PolicyID.ValueString()).PolicyRules().ByPolicyRuleId(object.ID.ValueString()).Delete(ctx, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
