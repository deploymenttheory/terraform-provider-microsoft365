// Package graphBetaNetworkWebFilteringPolicy manages Global Secure Access web filtering policies.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no webFilteringPolicy model in models/networkaccess
//     and no webFilteringPolicies request builder in networkaccess.
//   - Learn: none. Learn documents only the generic networkaccess.filteringPolicy surface
//     (https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringpolicy?view=graph-rest-beta), not
//     /networkaccess/webFilteringPolicies, and the endpoint is not in the beta $metadata.
//   - Observed: the endpoint and the settings.defaultAction and empty policyRules shapes follow the Entra admin centre's
//     Global Secure Access web content filtering requests to /networkaccess/webFilteringPolicies. Checked against the
//     SDK source on 2026-11-24; no live call, because the test tenant has no Global Secure Access licence. Requests go
//     through custom_requests in raw JSON mode with APIVersion left empty, so they use the beta client's own base URL
//     as the previous transport did.
//
// Removal condition: the beta SDK ships a networkaccess webFilteringPolicy model and a /networkAccess/webFilteringPolicies
// request builder.
package graphBetaNetworkWebFilteringPolicy

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
)

// webFilteringPoliciesEndpoint is the portal-observed collection; note the lower-case "networkaccess" segment.
const webFilteringPoliciesEndpoint = "/networkaccess/webFilteringPolicies"

// Create handles the Create operation for network web filtering policy resources.
//
// Operation: Creates a new web filtering policy
// API Calls:
//   - POST /networkaccess/webFilteringPolicies
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-web-content-filtering
func (r *NetworkWebFilteringPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkWebFilteringPolicyResourceModel

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

	var created webFilteringPolicy
	postConfig := customrequests.PostRequestConfig{
		Endpoint:    webFilteringPoliciesEndpoint,
		RawBody:     constructResource(ctx, &object, true),
		RawResponse: &created,
	}
	if _, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, nil, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.ID == nil {
		resp.Diagnostics.AddError("Error creating web filtering policy", "The API returned an invalid response without an id.")
		return
	}

	object.ID = types.StringValue(*created.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

// Read handles the Read operation for network web filtering policy resources.
//
// Operation: Retrieves the web filtering policy by ID
// API Calls:
//   - GET /networkaccess/webFilteringPolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-web-content-filtering
func (r *NetworkWebFilteringPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkWebFilteringPolicyResourceModel
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

	identity.ID = object.ID.ValueString()
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var policy webFilteringPolicy
	policyID := object.ID.ValueString()
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          webFilteringPoliciesEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        policyID,
		RawResponse:       &policy,
	}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, &policy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

// Update handles the Update operation for network web filtering policy resources.
//
// Operation: Updates the web filtering policy
// API Calls:
//   - PATCH /networkaccess/webFilteringPolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-web-content-filtering
func (r *NetworkWebFilteringPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkWebFilteringPolicyResourceModel
	var state NetworkWebFilteringPolicyResourceModel

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

	patchConfig := customrequests.PatchRequestConfig{
		Endpoint:          webFilteringPoliciesEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
		RawBody:           constructResource(ctx, &plan, false),
	}
	if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after update", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
	}
}

// Delete handles the Delete operation for network web filtering policy resources.
//
// Operation: Deletes the web filtering policy
// API Calls:
//   - DELETE /networkaccess/webFilteringPolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/how-to-configure-web-content-filtering
func (r *NetworkWebFilteringPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkWebFilteringPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deleteConfig := customrequests.DeleteRequestConfig{
		Endpoint:          webFilteringPoliciesEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
	}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)
}
