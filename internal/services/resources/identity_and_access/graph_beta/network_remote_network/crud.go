package graphBetaIdentityAndAccessNetworkRemoteNetwork

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

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

	profileIDs, diags := setStrings(ctx, plan.ForwardingProfileIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createBody, err := constructResource(plan.Name.ValueString(), plan.Region.ValueString(), profileIDs)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	var createResponse remoteNetworkResponse
	createErr := customrequests.JSONRequest(ctx, r.client.GetAdapter(), abstractions.POST, collectionPath, nil, createBody, &createResponse)
	created, err := decodeCreateResponse(createResponse, createErr)
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

	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		readbackOptions(plan, constants.TfOperationCreate),
	); err != nil {
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

	var readResponse remoteNetworkResponse
	readErr := customrequests.JSONRequest(
		ctx, r.client.GetAdapter(), abstractions.GET, itemPath, map[string]string{"remoteNetworkId": state.ID.ValueString()}, nil, &readResponse,
	)
	remote, err := &readResponse, readErr
	if err != nil {
		err = fmt.Errorf("get remote network: %w", err)

		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if err := mapResponse(&state, remote); err != nil {
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

		updateErr := customrequests.JSONRequest(
			ctx, r.client.GetAdapter(), abstractions.PATCH, itemPath, map[string]string{"remoteNetworkId": id}, constructNameUpdate(plan.Name.ValueString()), nil,
		)
		if err := updateErr; err != nil {
			err = fmt.Errorf("update remote network name: %w", err)

			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	if !plan.ForwardingProfileIDs.Equal(state.ForwardingProfileIDs) {
		profileIDs, diags := setStrings(ctx, plan.ForwardingProfileIDs)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		profileDelta, err := constructForwardingProfileDelta(profileIDs)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
		profileUpdateErr := customrequests.JSONRequest(
			ctx, r.client.GetAdapter(), abstractions.PATCH, forwardingProfilesPath, map[string]string{"remoteNetworkId": id}, profileDelta, nil,
		)
		if err := profileUpdateErr; err != nil {
			err = fmt.Errorf("update remote network forwarding profiles: %w", err)

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

	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		readbackOptions(plan, constants.TfOperationUpdate),
	); err != nil {
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

	deleteErr := customrequests.JSONRequest(
		ctx, r.client.GetAdapter(), abstractions.DELETE, itemPath, map[string]string{"remoteNetworkId": state.ID.ValueString()}, nil, nil,
	)
	if err := deleteErr; err != nil {
		err = fmt.Errorf("delete remote network: %w", err)

		// An earlier destroy may have started the deletion; wait for it anyway.
		if errors.GraphError(ctx, err).StatusCode != 404 {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}
	}

	if err := crud.PollUntil(ctx, deletionPollInterval, func(ctx context.Context) (bool, error) {
		var deletingNetwork remoteNetworkResponse
		readErr := customrequests.JSONRequest(
			ctx,
			r.client.GetAdapter(),
			abstractions.GET,
			itemPath,
			map[string]string{"remoteNetworkId": state.ID.ValueString()},
			nil,
			&deletingNetwork,
		)
		if ctx.Err() != nil {
			return false, &crud.FatalPollError{Err: fmt.Errorf("%w: %w", errDeletionTimedOut, ctx.Err())}
		}
		if readErr != nil {
			readErr = fmt.Errorf("get remote network: %w", readErr)
			info := errors.GraphError(ctx, readErr)
			switch info.StatusCode {
			case 404:
				if !strings.Contains(strings.ToLower(info.ErrorMessage), "deletion process") {
					return true, nil
				}
			case 401, 403:
				return false, &crud.FatalPollError{Err: readErr}
			}
		}
		return false, nil
	}); err != nil {
		if ctx.Err() != nil {
			err = fmt.Errorf("%w: %w", errDeletionTimedOut, ctx.Err())
		}

		resp.Diagnostics.AddError("Error waiting for remote network deletion", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

// readbackOptions configures shared read retries to verify the preceding write.
func readbackOptions(plan NetworkRemoteNetworkResourceModel, operation string) crud.ReadWithRetryOptions {
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = operation
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = stableConsistencyPredicate(plan)
	opts.MaxRetries = 1000 // Bounded by the operation timeout instead.
	return opts
}

func setStrings(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	var values []string
	if set.IsNull() || set.IsUnknown() {
		return values, nil
	}
	diags := set.ElementsAs(ctx, &values, false)
	return values, diags
}

const (
	collectionPath         = "/networkaccess/connectivity/remoteNetworks"
	itemPath               = collectionPath + "/{remoteNetworkId}"
	forwardingProfilesPath = itemPath + "/forwardingProfiles"
)

var (
	errMissingID            = stderrors.New("graph did not return an ID for the created remote network")
	errInvalidResponse      = stderrors.New("graph omitted required remote network fields")
	errNoForwardingProfiles = stderrors.New("at least one forwarding profile ID is required")
	errDeletionTimedOut     = stderrors.New("stopped waiting for remote network deletion to finish")
)

// deletionPollInterval is a variable so unit tests can shorten it.
var deletionPollInterval = 5 * time.Second
