package graphBetaNetworkDeviceForwardingProfile

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// Create handles the Create operation for network device forwarding profile resources.
//
// Operation: Creates a new network device forwarding profile
// API Calls:
//   - POST /networkaccess/deviceForwardingProfiles
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDeviceForwardingProfileResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan NetworkDeviceForwardingProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	body, diags := requestBody(ctx, plan, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// POST can return 504 after creating the profile, and a retried POST then fails
	// as a duplicate. List first so only a profile created by this request is adopted.
	var beforeCreate struct {
		Value []sharedmodels.ForwardingProfileWire `json:"value"`
	}
	listErr := customrequests.JSONRequest(ctx, r.client.GetAdapter(), abstractions.GET, r.ResourcePath, nil, nil, &beforeCreate)
	if listErr != nil {
		listErr = fmt.Errorf("list forwarding profiles: %w", listErr)
	}
	existing := make(map[string]bool, len(beforeCreate.Value))
	for _, profile := range beforeCreate.Value {
		if profile.ID != nil && profile.IsCustomProfile != nil && *profile.IsCustomProfile {
			existing[*profile.ID] = true
		}
	}
	if listErr != nil {
		tflog.Warn(ctx, "Could not list forwarding profiles before create", map[string]any{
			"error": listErr.Error(),
		})
	}
	var createResponse sharedmodels.ForwardingProfileWire
	createErr := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.POST,
		r.ResourcePath,
		nil,
		body,
		&createResponse,
	)
	remote, err := decodeCreateResponse(createResponse, createErr)
	// Without the list a match may be pre-existing, so report the error.
	if err != nil && listErr == nil {
		info := errors.GraphError(ctx, err)
		var created *sharedmodels.ForwardingProfileWire
		attempts := createdLookupAttemptsFor(info)
		attempt := 0
		_ = crud.PollUntil(ctx, createdLookupInterval, func(ctx context.Context) (bool, error) {
			attempt++
			var recoveryProfiles struct {
				Value []sharedmodels.ForwardingProfileWire `json:"value"`
			}
			lookupErr := customrequests.JSONRequest(ctx, r.client.GetAdapter(), abstractions.GET, r.ResourcePath, nil, nil, &recoveryProfiles)
			if lookupErr == nil {
				for i, profile := range recoveryProfiles.Value {
					if profile.ID != nil && !existing[*profile.ID] && profile.IsCustomProfile != nil && *profile.IsCustomProfile &&
						profile.Name != nil && *profile.Name == plan.Name.ValueString() && profile.TrafficForwardingType != nil &&
						*profile.TrafficForwardingType == plan.TrafficForwardingType.ValueString() && profile.Priority != nil &&
						*profile.Priority == plan.Priority.ValueInt32() && tryClaim(*profile.ID) {
						created = &recoveryProfiles.Value[i]
						return true, nil
					}
				}
				lookupErr = errProfileNotCreated
			} else {
				lookupErr = fmt.Errorf("list forwarding profiles: %w", lookupErr)
			}
			if !stderrors.Is(lookupErr, errProfileNotCreated) {
				lookupInfo := errors.GraphError(ctx, lookupErr)
				if lookupInfo.StatusCode >= 400 && lookupInfo.StatusCode < 500 && lookupInfo.StatusCode != 429 && lookupInfo.ErrorCode != "URLError" {
					return false, &crud.FatalPollError{Err: lookupErr}
				}
			}
			if attempt >= attempts {
				return false, &crud.FatalPollError{Err: lookupErr}
			}
			return false, lookupErr
		})
		if created != nil {
			tflog.Warn(
				ctx,
				"Forwarding profile create returned an error but the profile exists; adopting it",
				map[string]any{
					"status": info.StatusCode,
					"id":     *created.ID,
				},
			)
			remote, err = created, nil
		}
	}
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
	claimProfile(*remote.ID)
	// Persist the ID before reading back so a failed read still allows cleanup.
	plan.ID = types.StringPointerValue(remote.ID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Assignments apply asynchronously; wait until GET reflects the plan.
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		readbackOptions(plan, constants.TfOperationCreate),
	); err != nil {
		resp.Diagnostics.AddError("Error reading forwarding profile after create", err.Error())
	}
}

// Read handles the Read operation for network device forwarding profile resources.
//
// Operation: Retrieves the network device forwarding profile by ID
// API Calls:
//   - GET /networkaccess/deviceForwardingProfiles/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDeviceForwardingProfileResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state NetworkDeviceForwardingProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		state.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	var readResponse sharedmodels.ForwardingProfileWire
	readErr := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		r.ResourcePath+"/{profileId}",
		map[string]string{"profileId": state.ID.ValueString()},
		nil,
		&readResponse,
	)
	remote, err := &readResponse, readErr
	if err != nil {
		err = fmt.Errorf("get forwarding profile: %w", err)

		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, r.ReadPermissions)
		return
	}
	resp.Diagnostics.Append(mapRemoteStateToTerraform(ctx, &state, remote)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update handles the Update operation for network device forwarding profile resources.
//
// Operation: Updates the network device forwarding profile
// API Calls:
//   - PATCH /networkaccess/deviceForwardingProfiles/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDeviceForwardingProfileResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkDeviceForwardingProfileResourceModel
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

	body, diags := requestBody(ctx, plan, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateErr error
	if len(body) > 0 {
		updateErr = customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.PATCH,
			r.ResourcePath+"/{profileId}",
			map[string]string{"profileId": state.ID.ValueString()},
			body,
			nil,
		)
	}
	if err := updateErr; err != nil {
		err = fmt.Errorf("update forwarding profile: %w", err)

		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationUpdate,
			r.WritePermissions,
		)
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		readbackOptions(plan, constants.TfOperationUpdate),
	); err != nil {
		resp.Diagnostics.AddError("Error reading forwarding profile after update", err.Error())
	}
}

// Delete handles the Delete operation for network device forwarding profile resources.
//
// Operation: Deletes the network device forwarding profile
// API Calls:
//   - DELETE /networkaccess/deviceForwardingProfiles/{id}
//
// Reference: https://learn.microsoft.com/en-us/entra/global-secure-access/concept-traffic-forwarding
func (r *NetworkDeviceForwardingProfileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state NetworkDeviceForwardingProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		state.Timeouts.Delete,
		DeleteTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	deleteErr := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.DELETE,
		r.ResourcePath+"/{profileId}",
		map[string]string{"profileId": state.ID.ValueString()},
		nil,
		nil,
	)
	if err := deleteErr; err != nil {
		err = fmt.Errorf("delete forwarding profile: %w", err)

		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationDelete,
			r.WritePermissions,
		)
		return
	}
	// A stale list could still show it; never adopt it on a re-create.
	claimProfile(state.ID.ValueString())
	resp.State.RemoveResource(ctx)
}

// readbackOptions configures shared read retries to verify the preceding write.
func readbackOptions(
	plan NetworkDeviceForwardingProfileResourceModel,
	operation string,
) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = sharedmodels.RequireStableReads(func(ctx context.Context, s tfsdk.State) bool {
		var actual NetworkDeviceForwardingProfileResourceModel
		if s.Get(ctx, &actual).HasError() {
			return false
		}
		return matches(plan.Name, actual.Name) &&
			matches(plan.Description, actual.Description) &&
			matches(plan.Priority, actual.Priority) &&
			matches(plan.State, actual.State) &&
			matches(plan.ClientFallbackAction, actual.ClientFallbackAction) &&
			matches(plan.LocalNetworkSettings, actual.LocalNetworkSettings) &&
			matches(plan.AssignmentRules, actual.AssignmentRules)
	})
	return opts
}

// matches reports whether a configured value has reached Graph. Unknown planned
// values are computed by Graph and are not compared.
func matches(planned, actual attr.Value) bool {
	return planned.IsUnknown() || planned.Equal(actual)
}

var (
	errInvalidProfileResponse = stderrors.New("graph returned a forwarding profile without an ID")
	errProfileNotCreated      = stderrors.New("forwarding profile not found after failed create")
)

// claimedProfiles holds IDs already created or adopted in this process, so two
// resources with the same name and priority cannot both adopt one profile.
var claimedProfiles sync.Map

// Variables so unit tests can shorten the wait.
var (
	createdLookupAttempts   = 12
	duplicateLookupAttempts = 3
	createdLookupInterval   = 5 * time.Second
)

// createdLookupAttemptsFor polls longest after errors that can hide a create
// (transport errors, which GraphError reports as 400 "URLError", 503 or 504;
// gateway failures; non-Graph 500s). 400/409 get a few checks (a retried POST
// fails as a duplicate); anything else gets one.
func createdLookupAttemptsFor(info errors.GraphErrorInfo) int {
	switch {
	case info.ErrorCode == "URLError",
		info.StatusCode == 502, info.StatusCode == 503, info.StatusCode == 504,
		info.StatusCode == 500 && !info.IsODataError:
		return createdLookupAttempts
	case info.StatusCode == 400, info.StatusCode == 409:
		return duplicateLookupAttempts
	default:
		return 1
	}
}

// claimProfile records the profile ID as claimed by a resource in this provider process.
func claimProfile(id string) { claimedProfiles.Store(strings.ToLower(id), struct{}{}) }

// tryClaim claims id atomically and reports whether this caller got it.
func tryClaim(id string) bool {
	_, loaded := claimedProfiles.LoadOrStore(strings.ToLower(id), struct{}{})
	return !loaded
}
