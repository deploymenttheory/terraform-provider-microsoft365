package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/networkaccess"
)

func (r *NetworkFilteringProfilePolicyLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkFilteringProfilePolicyLinkResourceModel

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

	requestBody, err := constructCreateResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing policy link", err.Error())
		return
	}
	if err := populateComputedRequestFields(&object); err != nil {
		resp.Diagnostics.AddError("Error resolving policy link OData types", err.Error())
		return
	}

	link, err := r.createPolicyLink(ctx, object.FilteringProfileID.ValueString(), requestBody)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	createdID := link.GetId()
	if createdID == nil {
		resp.Diagnostics.AddError("Error creating policy link", "The API returned an invalid response without an id.")
		return
	}

	object.PolicyLinkID = types.StringValue(*createdID)
	object.ID = compositeID(object.FilteringProfileID.ValueString(), *createdID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName

	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading policy link state after create", err.Error())
		return
	}
}

func (r *NetworkFilteringProfilePolicyLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkFilteringProfilePolicyLinkResourceModel
	var identity sharedmodels.ResourceIdentity

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

	if object.FilteringProfileID.IsNull() || object.FilteringProfileID.IsUnknown() {
		resp.Diagnostics.AddError("Missing filtering_profile_id", "The filtering_profile_id attribute is required to read a policy link.")
		return
	}

	identity.ID = object.ID.ValueString()
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	profile, err := r.client.
		NetworkAccess().
		FilteringProfiles().
		ByFilteringProfileId(object.FilteringProfileID.ValueString()).
		Get(ctx, &networkaccess.FilteringProfilesFilteringProfileItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &networkaccess.FilteringProfilesFilteringProfileItemRequestBuilderGetQueryParameters{
				Expand: []string{"policies($expand=policy)"},
			},
		})
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	var policyLinkID string
	if !object.PolicyLinkID.IsNull() && !object.PolicyLinkID.IsUnknown() {
		policyLinkID = object.PolicyLinkID.ValueString()
	}
	var policyID string
	if !object.PolicyID.IsNull() && !object.PolicyID.IsUnknown() {
		policyID = object.PolicyID.ValueString()
	}

	link := findPolicyLink(profile, policyLinkID, policyID)
	if link == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, object.FilteringProfileID.ValueString(), link)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkFilteringProfilePolicyLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkFilteringProfilePolicyLinkResourceModel
	var state NetworkFilteringProfilePolicyLinkResourceModel

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

	requestBody, err := constructUpdateResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing policy link update", err.Error())
		return
	}
	if err := populateComputedRequestFields(&plan); err != nil {
		resp.Diagnostics.AddError("Error resolving policy link OData types", err.Error())
		return
	}

	if err := r.updatePolicyLink(ctx, state.FilteringProfileID.ValueString(), state.PolicyLinkID.ValueString(), requestBody); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = state.ID
	plan.PolicyLinkID = state.PolicyLinkID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading policy link state after update", err.Error())
		return
	}
}

func (r *NetworkFilteringProfilePolicyLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkFilteringProfilePolicyLinkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		NetworkAccess().
		FilteringProfiles().
		ByFilteringProfileId(object.FilteringProfileID.ValueString()).
		Policies().
		ByPolicyLinkId(object.PolicyLinkID.ValueString()).
		Delete(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		if !resp.Diagnostics.HasError() {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	resp.State.RemoveResource(ctx)
}
