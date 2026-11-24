// Package graphBetaIdentityAndAccessNetworkRemoteNetwork manages Global Secure Access remote networks.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0,
//     networkaccess.ConnectivityRemoteNetworksItemForwardingProfilesRequestBuilder has Get and Post but no Patch, so the
//     $delta PATCH that replaces a remote network's forwarding profile associations cannot be sent through the SDK.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/networkaccess-remotenetwork-update?view=graph-rest-beta
//     documents the entity PATCH only; the forwardingProfiles $delta PATCH is not documented.
//   - Observed: when the resource was added (#4862) the entity PATCH rejected forwardingProfiles and the profile-side
//     associations overwrote other remote networks, so associations are replaced with
//     PATCH /networkaccess/connectivity/remoteNetworks/{id}/forwardingProfiles carrying {"@context":"#$delta","value":[...]}.
//     Checked against the SDK source on 2026-11-24; no live call, because the test tenant has no Global Secure Access
//     licence. The entity POST, GET, PATCH and DELETE use the same custom_requests raw JSON transport so their request
//     bodies keep the live-verified shape (the generated models would add @odata.type annotations to the entity and to
//     each forwarding profile reference). APIVersion is left empty so the requests use the beta client's own base URL,
//     as the previous transport did.
//
// Removal condition: the beta SDK ships a Patch on ConnectivityRemoteNetworksItemForwardingProfilesRequestBuilder that
// accepts a $delta collection of forwarding profile references.
package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// remoteNetworksEndpoint is the remote networks collection.
const remoteNetworksEndpoint = "/networkaccess/connectivity/remoteNetworks"

var (
	errMissingID        = stderrors.New("graph did not return an ID for the created remote network")
	errInvalidResponse  = stderrors.New("graph omitted required remote network fields")
	errDeletionTimedOut = stderrors.New("stopped waiting for remote network deletion to finish")
)

// deletionPollInterval is a variable so unit tests can shorten it.
var deletionPollInterval = 5 * time.Second

// readbackPaths are the attributes the read-back after a write waits for. MaxRetries is raised so the operation
// timeout, not the retry count, bounds the wait.
var readbackPaths = crud.RootPaths("id", "name", "region", "forwarding_profile_ids")

// Create handles the Create operation for remote network resources.
//
// Operation: Creates a new remote network
// API Calls:
//   - POST /networkaccess/connectivity/remoteNetworks
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-connectivity-post-remotenetworks?view=graph-rest-beta
func (r *NetworkRemoteNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworkRemoteNetworkResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	body, err := constructResource(ctx, plan)
	var created *remoteNetworkResponse
	if err == nil {
		postConfig := customrequests.PostRequestConfig{
			Endpoint:    remoteNetworksEndpoint,
			RawBody:     body,
			RawResponse: &created,
		}
		_, err = customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, nil, nil)
	}
	if err == nil && (created == nil || created.ID == nil || *created.ID == "") {
		err = errMissingID
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// Keep the ID in state even if the read-back fails, so the object can still be destroyed.
	plan.ID = types.StringValue(*created.ID)
	plan.Version = types.StringPointerValue(created.Version)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	predicate := crud.RequireStableReads(crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{Paths: readbackPaths}))
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithMaxRetries(1000), crud.ReadWithConsistencyPredicate(predicate))
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for remote network resources.
//
// Operation: Retrieves a remote network by ID
// API Calls:
//   - GET /networkaccess/connectivity/remoteNetworks/{id}
//
// The item GET is used because the forwardingProfiles navigation GET returns an empty collection even when profiles
// are associated.
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-remotenetwork-get?view=graph-rest-beta
func (r *NetworkRemoteNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworkRemoteNetworkResourceModel

	// During read-back a 404 is replica lag, not a deletion.
	operation := constants.TfOperationRead
	if ctxOp, ok := ctx.Value("retry_operation").(string); ok && ctxOp != "" {
		operation = ctxOp
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// An empty or null body leaves remote nil, which MapRemoteStateToTerraform reports as an invalid response.
	var remote *remoteNetworkResponse
	getConfig := customrequests.GetRequestConfig{
		Endpoint:          remoteNetworksEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
		RawResponse:       &remote,
	}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if err := MapRemoteStateToTerraform(&state, remote); err != nil {
		resp.Diagnostics.AddError("Invalid Graph response", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles the Update operation for remote network resources.
//
// Operation: Updates the remote network name and forwarding profile membership
// API Calls:
//   - PATCH /networkaccess/connectivity/remoteNetworks/{id} (name)
//   - PATCH /networkaccess/connectivity/remoteNetworks/{id}/forwardingProfiles ($delta, full desired set)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-remotenetwork-update?view=graph-rest-beta
func (r *NetworkRemoteNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkRemoteNetworkResourceModel

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

	id := state.ID.ValueString()

	if !plan.Name.Equal(state.Name) {
		// The name PATCH leaves the forwarding profiles unchanged.
		patchConfig := customrequests.PatchRequestConfig{
			Endpoint:          remoteNetworksEndpoint,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        id,
			RawBody:           map[string]any{"name": plan.Name.ValueString()},
		}
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	if !plan.ForwardingProfileIDs.Equal(state.ForwardingProfileIDs) {
		var profileIDs []string
		if !plan.ForwardingProfileIDs.IsNull() && !plan.ForwardingProfileIDs.IsUnknown() {
			resp.Diagnostics.Append(plan.ForwardingProfileIDs.ElementsAs(ctx, &profileIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		// Send the full set as a $delta PATCH. The entity PATCH rejects forwardingProfiles, and the profile-side
		// associations would overwrite other remote networks.
		err := errNoForwardingProfiles
		if len(profileIDs) > 0 {
			refs := make([]map[string]string, 0, len(profileIDs))
			for _, profileID := range profileIDs {
				refs = append(refs, map[string]string{"id": profileID})
			}
			patchConfig := customrequests.PatchRequestConfig{
				Endpoint:          remoteNetworksEndpoint,
				ResourceIDPattern: customrequests.ResourceIDPatternSegment,
				ResourceID:        id,
				EndpointSuffix:    "/forwardingProfiles",
				RawBody:           map[string]any{"@context": "#$delta", "value": refs},
			}
			err = customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig)
		}
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	// Read keeps the model it is given, so it needs the planned timeouts.
	state.Timeouts = plan.Timeouts
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	predicate := crud.RequireStableReads(crud.AttributesMatchPredicate(req.Plan, crud.AttributeMatchOptions{Paths: readbackPaths}))
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName, crud.ReadWithMaxRetries(1000), crud.ReadWithConsistencyPredicate(predicate))
	readRequest := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	if err := crud.ReadWithRetry(ctx, r.Read, readRequest, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, id))
}

// Delete handles the Delete operation for remote network resources.
//
// Operation: Deletes the remote network and waits for deletion to finish
// API Calls:
//   - DELETE /networkaccess/connectivity/remoteNetworks/{id}
//   - GET /networkaccess/connectivity/remoteNetworks/{id} (until the deletion finishes)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-connectivity-delete-remotenetworks?view=graph-rest-beta
func (r *NetworkRemoteNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworkRemoteNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deleteConfig := customrequests.DeleteRequestConfig{
		Endpoint:          remoteNetworksEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
	}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		// An earlier destroy may have started the deletion; wait for it anyway.
		if errors.GraphError(ctx, err).StatusCode != 404 {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}
	}

	// Graph deletes asynchronously. Until it finishes, GET returns 404 "... is on deletion process." and the name
	// cannot be reused (409). The DELETE succeeded, so keep polling through transient errors.
	err := crud.PollUntil(ctx, deletionPollInterval, func(ctx context.Context) (bool, error) {
		getConfig := customrequests.GetRequestConfig{
			Endpoint:          remoteNetworksEndpoint,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        state.ID.ValueString(),
		}
		_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
		if err == nil {
			return false, nil
		}
		info := errors.GraphError(ctx, err)
		switch info.StatusCode {
		case 404:
			return !strings.Contains(strings.ToLower(info.ErrorMessage), "deletion process"), nil
		case 401, 403:
			return false, &crud.FatalPollError{Err: err}
		}
		return false, nil
	})
	if err != nil && ctx.Err() != nil {
		err = fmt.Errorf("%w: %w", errDeletionTimedOut, ctx.Err())
	}
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for remote network deletion", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}
