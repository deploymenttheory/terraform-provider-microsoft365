// Package graphBetaNetworkFilteringProfile manages Global Secure Access filtering profiles.
//
// SDK deviation: serializer sends read-only fields or nulls
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 networkaccess.FilteringProfile serializes the
//     read-only createdDateTime, lastModifiedDateTime and version fields of its Profile base, and NewFilteringProfile
//     always writes the @odata.type discriminator.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-update?view=graph-rest-beta
//     lists name, description, priority and state as the updatable properties; Learn does not document create.
//   - Observed: the previous implementation recorded that POST /networkAccess/filteringProfiles rejected the SDK body
//     with "400 Invalid request parameters" and accepted a body of only name, description, priority and state. Not
//     re-run on 2026-11-24, because the test tenant has no Global Secure Access licence. POST and PATCH go through
//     custom_requests in raw JSON mode with APIVersion left empty, so they use the beta client's own base URL; GET and
//     DELETE stay on the SDK.
//
// Removal condition: a live POST and PATCH with the SDK FilteringProfile model, setting only name, description,
// priority and state, succeed.
package graphBetaNetworkFilteringProfile

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
)

// filteringProfilesEndpoint is the filtering profile collection.
const filteringProfilesEndpoint = "/networkAccess/filteringProfiles"

// Create handles the Create operation for Filtering Profile resources.
//
// Operation: Creates a new network filtering profile for Global Secure Access
// API Calls:
//   - POST /networkAccess/filteringProfiles
//
// Microsoft Learn documents filteringProfile plus list/get/update operations,
// while live Graph probing and the beta SDK confirm POST on this collection.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringprofile?view=graph-rest-beta
func (r *NetworkFilteringProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkFilteringProfileResourceModel

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

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	postConfig := customrequests.PostRequestConfig{
		Endpoint: filteringProfilesEndpoint,
		RawBody:  requestBody,
	}
	baseResource, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), postConfig, models.CreateFilteringProfileFromDiscriminatorValue, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	profile, ok := baseResource.(models.FilteringProfileable)
	if !ok {
		resp.Diagnostics.AddError("Error creating filtering profile", fmt.Sprintf("The API returned %T, expected a filtering profile.", baseResource))
		return
	}

	createdID := profile.GetId()
	if createdID == nil {
		resp.Diagnostics.AddError(
			"Error creating filtering profile",
			"The API returned an invalid response without an id.",
		)
		return
	}

	object.ID = types.StringValue(*createdID)
	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, *createdID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName)

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for Filtering Profile resources.
//
// Operation: Retrieves a network filtering profile by ID
// API Calls:
//   - GET /networkAccess/filteringProfiles/{filteringProfileId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-get?view=graph-rest-beta
func (r *NetworkFilteringProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkFilteringProfileResourceModel
	var identity sharedmodels.ResourceIdentity

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

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

	policy, err := r.client.
		NetworkAccess().
		FilteringProfiles().
		ByFilteringProfileId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, policy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Filtering Profile resources.
//
// Operation: Updates an existing network filtering profile
// API Calls:
//   - PATCH /networkAccess/filteringProfiles/{filteringProfileId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-update?view=graph-rest-beta
func (r *NetworkFilteringProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkFilteringProfileResourceModel
	var state NetworkFilteringProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)   // desired state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...) // current state (for ID)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	profileID := state.ID.ValueString()
	patchConfig := customrequests.PatchRequestConfig{
		Endpoint:          filteringProfilesEndpoint,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        profileID,
		RawBody:           requestBody,
	}
	if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName)

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Filtering Profile resources.
//
// Operation: Deletes a network filtering profile
// API Calls:
//   - DELETE /networkAccess/filteringProfiles/{filteringProfileId}
//
// Microsoft Learn does not currently document delete for filteringProfile, but
// the beta SDK exposes the DELETE builder and live verification confirmed it.
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-filteringprofile?view=graph-rest-beta
func (r *NetworkFilteringProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkFilteringProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

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
		ByFilteringProfileId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
