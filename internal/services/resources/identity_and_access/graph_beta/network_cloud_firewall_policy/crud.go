package graphBetaNetworkCloudFirewallPolicy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graph "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

func (r *NetworkCloudFirewallPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkCloudFirewallPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	createBuilder := r.client.NetworkAccess().CloudFirewallPolicies()
	createRequest, err := createBuilder.ToPostRequestInformation(ctx, constructResource(ctx, &object), nil)
	if err != nil {
		err = fmt.Errorf("build cloud firewall create request: %w", err)
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	createRetryOptions := r.retryOptions
	createRetryOptions.ShouldRetry = func(_ time.Duration, _ int, _ *http.Request, response *http.Response) bool {
		return response.StatusCode == http.StatusTooManyRequests
	}
	createRequest.AddRequestOptions([]abstractions.RequestOption{&createRetryOptions})
	createResult, err := r.client.GetAdapter().
		Send(ctx, createRequest, graph.CreateCloudFirewallPolicyFromDiscriminatorValue, abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	created, err := decodeCreatePolicyResponse(createResult, err)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created == nil || created.GetId() == nil || *created.GetId() == "" {
		resp.Diagnostics.AddError("Invalid cloud firewall create response", "The API did not return a resource ID.")
		return
	}
	object.ID = types.StringValue(*created.GetId())
	object.Version = types.StringNull()
	object.LastModifiedDateTime = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	remote, err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(object.ID.ValueString()).Get(ctx, nil)
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

func (r *NetworkCloudFirewallPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkCloudFirewallPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	remote, err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(object.ID.ValueString()).Get(ctx, nil)
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
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
	}
}

func (r *NetworkCloudFirewallPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkCloudFirewallPolicyResourceModel
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
	body := constructPolicy(&plan, &state)
	if body.changed {
		if _, err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(state.ID.ValueString()).Patch(ctx, body, nil); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}
	// Keep the prior known state if readback fails after a successful update.
	object := state
	object.Timeouts = plan.Timeouts
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	remote, err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(object.ID.ValueString()).Get(ctx, nil)
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

func (r *NetworkCloudFirewallPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkCloudFirewallPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	if err := r.client.NetworkAccess().CloudFirewallPolicies().ByCloudFirewallPolicyId(object.ID.ValueString()).Delete(ctx, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
