// Package graphBetaApplicationsOnPremisesIpApplicationSegment manages IP application segments of applications
// published through application proxy.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no request builder for
//     /applications/{id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments.
//     The applications package stops at the application item, and onpremisespublishingprofiles only builds the
//     tenant-level /onPremisesPublishingProfiles/{id}/applicationSegments collection.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/onpremisespublishingprofile-post-applicationsegments?view=graph-rest-beta
//   - Observed: read from the SDK source; the application-scoped path is the one Learn documents for create, update and
//     delete, and the one the resource has always called.
//
// Removal condition: the beta SDK ships a request builder for the application-scoped
// segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments collection and item.
//
// SDK deviation: unknown enum value
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, models.IpApplicationSegment.destinationType uses
//     models.PrivateNetworkDestinationType, whose parser knows ipAddress, ipRange, ipRangeCidr, fqdn, dnsSuffix,
//     unknownFutureValue and servicePrincipalName, but not "ip".
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/ipapplicationsegment?view=graph-rest-beta
//   - Observed: 2026-07, POST of a single IP segment with destinationType "ipAddress" returned
//     400 InvalidJson_BadRequest, while "ip" succeeded and GET returned "destinationType": "ip". The request and
//     response bodies are therefore plain structs (see model.go) sent through custom_requests in raw JSON mode.
//
// Removal condition: models.PrivateNetworkDestinationType includes "ip", or the endpoint accepts and returns
// "ipAddress" as Learn documents, and the missing type deviation above is removed.
package graphBetaApplicationsOnPremisesIpApplicationSegment

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// segmentsEndpointSuffix follows the application ID in every segment request path.
const segmentsEndpointSuffix = "/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments"

// Create handles the Create operation for IP application segment resources.
//
// Operation: Creates a new IP application segment for application proxy
// API Calls:
//   - POST /applications/{id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments
//
// Reference: https://learn.microsoft.com/en-us/graph/api/onpremisespublishingprofile-post-applicationsegments?view=graph-rest-beta
func (r *OnPremisesIpApplicationSegmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object OnPremisesIpApplicationSegmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Concurrent writes to the same application can return a segment id that then 404s on read, so sibling
	// segments are written one at a time, each with its read-back. The lock is taken before the operation timeout
	// starts, so time spent waiting on sibling segments does not consume this resource's own timeout.
	unlock, err := crud.AcquireKeyedLock(ctx, ResourceName, object.ApplicationObjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for ip application segment lock", err.Error())
		return
	}
	defer unlock()

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing ip application segment", fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()))
		return
	}

	// Retry explicit throttling only: ambiguous server failures can leave an untracked duplicate.
	createRetryOptions := r.retryOptions
	createRetryOptions.ShouldRetry = func(_ time.Duration, _ int, _ *http.Request, response *http.Response) bool {
		return response != nil && response.StatusCode == http.StatusTooManyRequests
	}

	var created ipApplicationSegmentResponse
	createConfig := customrequests.PostRequestConfig{
		APIVersion:     customrequests.GraphAPIBeta,
		Endpoint:       "applications/" + object.ApplicationObjectID.ValueString() + segmentsEndpointSuffix,
		RawBody:        requestBody,
		RawResponse:    &created,
		RequestOptions: []abstractions.RequestOption{&createRetryOptions},
	}
	if _, err := customrequests.PostRequest(ctx, r.client.GetAdapter(), createConfig, nil, nil); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if created.ID == nil {
		resp.Diagnostics.AddError("Error creating ip application segment", "The API returned an invalid response without an id.")
		return
	}

	object.ID = types.StringValue(*created.ID)

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

// Read handles the Read operation for IP application segment resources.
//
// Operation: Retrieves an IP application segment by ID
// API Calls:
//   - GET /applications/{id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/{segmentId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/ipapplicationsegment-get?view=graph-rest-beta
func (r *OnPremisesIpApplicationSegmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object OnPremisesIpApplicationSegmentResourceModel
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

	// A response with no body leaves ipApplicationSegment nil.
	var ipApplicationSegment *ipApplicationSegmentResponse
	getConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "applications/" + object.ApplicationObjectID.ValueString() + segmentsEndpointSuffix,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
		RawResponse:       &ipApplicationSegment,
	}
	if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if ipApplicationSegment == nil {
		resp.Diagnostics.AddError("Error reading ip application segment", fmt.Sprintf("Received nil ip application segment for ID: %s", object.ID.ValueString()))
		return
	}

	MapRemoteStateToTerraform(ctx, &object, ipApplicationSegment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for IP application segment resources.
//
// Operation: Updates an existing IP application segment
// API Calls:
//   - PATCH /applications/{id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/{segmentId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/ipapplicationsegment-update?view=graph-rest-beta
func (r *OnPremisesIpApplicationSegmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OnPremisesIpApplicationSegmentResourceModel
	var state OnPremisesIpApplicationSegmentResourceModel

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

	// Serialize with sibling segments on the same application (see Create).
	unlock, err := crud.AcquireKeyedLock(ctx, ResourceName, state.ApplicationObjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for ip application segment lock", err.Error())
		return
	}
	defer unlock()

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing ip application segment", fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()))
		return
	}

	patchConfig := customrequests.PatchRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "applications/" + state.ApplicationObjectID.ValueString() + segmentsEndpointSuffix,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID.ValueString(),
		RawBody:           requestBody,
	}
	if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.WritePermissions)
		return
	}

	// Keep the prior known state if the read-back fails after a successful update.
	object := state
	object.Timeouts = plan.Timeouts

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

// Delete handles the Delete operation for IP application segment resources.
//
// Operation: Deletes an IP application segment
// API Calls:
//   - DELETE /applications/{id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments/{segmentId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/onpremisespublishingprofile-delete-applicationsegments?view=graph-rest-beta
func (r *OnPremisesIpApplicationSegmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OnPremisesIpApplicationSegmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Serialize with sibling segments on the same application (see Create).
	unlock, err := crud.AcquireKeyedLock(ctx, ResourceName, data.ApplicationObjectID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for ip application segment lock", err.Error())
		return
	}
	defer unlock()

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deleteConfig := customrequests.DeleteRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "applications/" + data.ApplicationObjectID.ValueString() + segmentsEndpointSuffix,
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
