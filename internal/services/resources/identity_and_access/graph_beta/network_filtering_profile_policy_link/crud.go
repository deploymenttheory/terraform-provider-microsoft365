// Package graphBetaNetworkFilteringProfilePolicyLink manages the policy links of Global Secure Access filtering profiles.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no webFilteringPolicyLink, promptPolicyLink,
//     filePolicyLink or securityProviderPolicyLink model in models/networkaccess, and CreatePolicyLinkFromDiscriminatorValue
//     does not know those discriminators, so the SDK cannot send these link types.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-policylink?view=graph-rest-beta lists
//     only the SDK's derived types; the portal-first link types are not in the beta $metadata.
//   - Observed: the Entra admin centre links V2 web filtering policies by POSTing a webFilteringPolicyLink body to
//     /networkAccess/filteringProfiles/{id}/policies, and its content policy option sends filePolicyLink/filePolicy.
//     Checked against the SDK source on 2026-11-24; no live call, because the test tenant has no Global Secure Access
//     licence. POST and PATCH go through custom_requests in raw JSON mode with APIVersion left empty, so they use the
//     beta client's own base URL; the GET with $expand and the DELETE stay on the SDK.
//
// Removal condition: the beta SDK ships the webFilteringPolicyLink, promptPolicyLink, filePolicyLink and
// securityProviderPolicyLink models in the policyLink discriminator.
package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
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

	requestBody, err := constructResource(ctx, &object, true)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing policy link", err.Error())
		return
	}

	postConfig := customrequests.PostRequestConfig{
		Endpoint: "/networkAccess/filteringProfiles/" + object.FilteringProfileID.ValueString() + "/policies",
		RawBody:  requestBody,
	}
	result, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, models.CreatePolicyLinkFromDiscriminatorValue, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	link, ok := result.(models.PolicyLinkable)
	if !ok || link.GetId() == nil {
		resp.Diagnostics.AddError("Error creating policy link", "The API returned an invalid response without an id.")
		return
	}

	object.PolicyLinkID = types.StringValue(*link.GetId())
	object.ID = types.StringValue(object.FilteringProfileID.ValueString() + "/" + *link.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName)

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

	// Match the link by its id, or by the linked policy id when the link id is not yet known.
	var link models.PolicyLinkable
	if profile != nil {
		for _, candidate := range profile.GetPolicies() {
			if candidate == nil {
				continue
			}
			if policyLinkID != "" {
				if id := candidate.GetId(); id != nil && *id == policyLinkID {
					link = candidate
					break
				}
				continue
			}
			if policy := candidate.GetPolicy(); policyID != "" && policy != nil {
				if id := policy.GetId(); id != nil && *id == policyID {
					link = candidate
					break
				}
			}
		}
	}
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

	requestBody, err := constructResource(ctx, &plan, false)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing policy link update", err.Error())
		return
	}

	linksEndpoint := "/networkAccess/filteringProfiles/" + state.FilteringProfileID.ValueString() + "/policies"
	patchConfig := customrequests.PatchRequestConfig{
		Endpoint:          linksEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.PolicyLinkID.ValueString(),
		RawBody:           requestBody,
	}
	if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
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
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName)

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
