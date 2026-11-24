// Package graphBetaApplicationsOnPremisesConnectorGroup manages application proxy connector groups.
//
// SDK deviation: unknown enum value
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, models.ConnectorGroupRegion.
//     ParseConnectorGroupRegion knows nam, eur, aus, asia, ind and unknownFutureValue, and
//     returns nil for anything else.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
//     documents connectorGroupType values that are not in the SDK enum either.
//   - Observed: 2026-07-05, GET /beta/onPremisesPublishingProfiles/applicationProxy/connectorGroups
//     returned the default connector group with "region": "japan", which the SDK model drops.
//     The same check showed create returning HTTP 200 instead of the documented 201, read-only
//     values such as isDefault being ignored on create, and PATCH returning 204 with no body.
//     Every request therefore goes through custom_requests in raw JSON mode, with region and
//     connectorGroupType decoded as strings (see connectorGroupResponse in model.go).
//
// Removal condition: models.ConnectorGroupRegion includes every region Graph returns
// (at least japan), and models.ConnectorGroupType includes the values Learn documents.
package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for connector group resources.
//
// Operation: Creates a new application proxy connector group
// API Calls:
//   - POST /onPremisesPublishingProfiles/applicationProxy/connectorGroups
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connectorgroup-post?view=graph-rest-beta
func (r *OnPremisesConnectorGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object OnPremisesConnectorGroupResourceModel

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

	var created connectorGroupResponse
	createConfig := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    r.ResourcePath,
		RawBody:     constructResource(ctx, &object, nil),
		RawResponse: &created,
	}
	if _, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), createConfig, nil, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if created.ID == nil {
		resp.Diagnostics.AddError("Error creating connector group", "The API returned an invalid response without an id.")
		return
	}

	object.ID = types.StringValue(*created.ID)
	MapRemoteStateToTerraform(ctx, &object, &created)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for connector group resources.
//
// Operation: Retrieves a connector group by ID
// API Calls:
//   - GET /onPremisesPublishingProfiles/applicationProxy/connectorGroups/{connectorGroupId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connectorgroup-get?view=graph-rest-beta
func (r *OnPremisesConnectorGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object OnPremisesConnectorGroupResourceModel
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

	// A response with no body leaves connectorGroup nil.
	var connectorGroup *connectorGroupResponse
	getConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
		RawResponse:       &connectorGroup,
	}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if connectorGroup == nil {
		resp.Diagnostics.AddError("Error reading connector group", fmt.Sprintf("Received nil connector group for ID: %s", object.ID.ValueString()))
		return
	}

	MapRemoteStateToTerraform(ctx, &object, connectorGroup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for connector group resources.
//
// Operation: Updates the name, and the region when it changed, of a connector group
// API Calls:
//   - PATCH /onPremisesPublishingProfiles/applicationProxy/connectorGroups/{connectorGroupId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connectorgroup-update?view=graph-rest-beta
func (r *OnPremisesConnectorGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OnPremisesConnectorGroupResourceModel
	var state OnPremisesConnectorGroupResourceModel

	operation := constants.TfOperationUpdate
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting update of resource: %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// PATCH returns 204 with no body, which leaves updated nil.
	var updated *connectorGroupResponse
	patchConfig := customrequests.PatchRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
		RawBody:           constructResource(ctx, &plan, &state),
		RawResponse:       &updated,
	}
	if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.WritePermissions)
		return
	}

	// Keep the prior known state if the read-back fails after a successful update.
	object := state
	object.Timeouts = plan.Timeouts
	if updated != nil {
		MapRemoteStateToTerraform(ctx, &object, updated)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after update", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for connector group resources.
//
// Operation: Deletes a connector group, or only removes the tenant default connector group from state
// API Calls:
//   - DELETE /onPremisesPublishingProfiles/applicationProxy/connectorGroups/{connectorGroupId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connectorgroup-delete?view=graph-rest-beta
func (r *OnPremisesConnectorGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OnPremisesConnectorGroupResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if !data.IsDefault.IsNull() && !data.IsDefault.IsUnknown() && data.IsDefault.ValueBool() {
		// Direct API verification on 2026-07-05 observed the tenant default
		// connector group with isDefault=true. Because the default is
		// system-managed, attempting a real DELETE would be risky if Graph ever
		// accepted it. Treat destroy of an imported default connector group like
		// AWS default resources: remove it from Terraform state, but leave the
		// remote default in place.
		//
		// Learn documents isDefault as read-only:
		// https://learn.microsoft.com/en-us/graph/api/resources/connectorgroup?view=graph-rest-beta
		tflog.Warn(ctx, "Default connector group cannot be deleted in Microsoft Graph; removing from Terraform state only", map[string]any{
			"connector_group_id": data.ID.ValueString(),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	deleteConfig := customrequests.DeleteRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        data.ID.ValueString(),
	}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
