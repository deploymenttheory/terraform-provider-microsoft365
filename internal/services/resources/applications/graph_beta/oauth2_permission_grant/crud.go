package graphBetaOAuth2PermissionGrant

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

const grantPollInterval = 2 * time.Second

var errGrantResponse = errors.New("delegated permission grant or grant ID missing from Microsoft Graph response")

// Create posts an owned consent grant to /beta/oauth2PermissionGrants.
func (r *OAuth2PermissionGrantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object OAuth2PermissionGrantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	body, err := constructResource(ctx, &object, true)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing delegated permission grant", err.Error())
		return
	}
	var created graphmodels.OAuth2PermissionGrantable
	err = crud.PollUntil(ctx, grantPollInterval, func(ctx context.Context) (bool, error) {
		var postErr error
		created, postErr = r.client.Oauth2PermissionGrants().Post(ctx, body, nil)
		if postErr == nil {
			return true, nil
		}
		info := grapherrors.GraphError(ctx, postErr)
		message := strings.ToLower(info.ErrorMessage)
		if info.StatusCode == 404 || (info.StatusCode == 400 && (strings.Contains(message, "does not exist") || strings.Contains(message, "not found"))) {
			return false, fmt.Errorf("waiting for referenced service principals: %w", postErr)
		}
		return false, &crud.FatalPollError{Err: postErr}
	})
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created == nil || created.GetId() == nil || *created.GetId() == "" {
		resp.Diagnostics.AddError("Error extracting delegated permission grant ID", errGrantResponse.Error())
		return
	}
	object.ID = types.StringValue(*created.GetId())
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	container := &crud.CreateResponseContainer{CreateResponse: resp}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = grantConsistencyPredicate(&object)
	if err = crud.ReadWithRetry(ctx, r.Read, readReq, container, opts); err != nil {
		resp.Diagnostics.AddError("Error reading delegated permission grant after create", err.Error())
	}
}

// Read refreshes a grant by its opaque ID. Repeated 404s confirm deletion during
// ordinary refresh; a post-write 404 remains an error for shared read retries.
func (r *OAuth2PermissionGrantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object OAuth2PermissionGrantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	operation := constants.TfOperationRead
	if op, ok := ctx.Value("retry_operation").(string); ok {
		operation = op
	}
	var remote graphmodels.OAuth2PermissionGrantable
	missing := 0
	err := crud.PollUntil(ctx, grantPollInterval, func(ctx context.Context) (bool, error) {
		var getErr error
		remote, getErr = r.client.Oauth2PermissionGrants().ByOAuth2PermissionGrantId(object.ID.ValueString()).Get(ctx, nil)
		if getErr == nil {
			return true, nil
		}
		info := grapherrors.GraphError(ctx, getErr)
		if info.StatusCode == 400 && strings.Contains(info.ErrorMessage, "Filter is invalid. principalId") {
			return false, fmt.Errorf("waiting for user consent replication: %w", getErr)
		}
		if grapherrors.GraphError(ctx, getErr).StatusCode == 404 && operation == constants.TfOperationRead {
			missing++
			return missing >= 3, nil
		}
		return false, &crud.FatalPollError{Err: getErr}
	})
	if err != nil {
		grapherrors.HandleKiotaGraphErrorWithOptions(
			ctx,
			err,
			resp,
			operation,
			r.ReadPermissions,
			grapherrors.GraphErrorOptions{PreserveStateOnReadBadRequest: true, PreserveStateOnReadNotFound: true},
		)
		return
	}
	if missing >= 3 {
		resp.State.RemoveResource(ctx)
		return
	}
	if remote == nil || remote.GetId() == nil {
		resp.Diagnostics.AddError("Invalid delegated permission grant response", errGrantResponse.Error())
		return
	}
	MapRemoteStateToTerraform(ctx, &object, remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
	}
}

// Update changes the complete scope set with a scope-only beta PATCH.
func (r *OAuth2PermissionGrantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state OAuth2PermissionGrantResourceModel
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
	body, err := constructResource(ctx, &plan, false)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing delegated permission grant", err.Error())
		return
	}
	err = crud.PollUntil(ctx, grantPollInterval, func(ctx context.Context) (bool, error) {
		_, patchErr := r.client.Oauth2PermissionGrants().ByOAuth2PermissionGrantId(state.ID.ValueString()).Patch(ctx, body, nil)
		if patchErr == nil {
			return true, nil
		}
		info := grapherrors.GraphError(ctx, patchErr)
		if info.StatusCode == 404 ||
			(info.StatusCode == 400 && (strings.Contains(info.ErrorMessage, "Permission being updated or deleted is not found") || strings.Contains(info.ErrorMessage, "Filter is invalid. principalId"))) {
			return false, fmt.Errorf("waiting for delegated permission grant update: %w", patchErr)
		}
		return false, &crud.FatalPollError{Err: patchErr}
	})
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}
	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	container := &crud.UpdateResponseContainer{UpdateResponse: resp}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = grantConsistencyPredicate(&plan)
	if err = crud.ReadWithRetry(ctx, r.Read, readReq, container, opts); err != nil {
		resp.Diagnostics.AddError("Error reading delegated permission grant after update", err.Error())
	}
}

// Delete removes the grant itself, then verifies its disappearance before
// removing state. A missing-object write response alone can be a stale replica.
func (r *OAuth2PermissionGrantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object OAuth2PermissionGrantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	missing := 0
	err := crud.PollUntil(ctx, grantPollInterval, func(ctx context.Context) (bool, error) {
		deleteErr := r.client.Oauth2PermissionGrants().ByOAuth2PermissionGrantId(object.ID.ValueString()).Delete(ctx, nil)
		if deleteErr != nil {
			info := grapherrors.GraphError(ctx, deleteErr)
			missingWrite := info.StatusCode == 404 ||
				(info.StatusCode == 400 && (strings.Contains(info.ErrorMessage, "Permission being updated or deleted is not found") || strings.Contains(info.ErrorMessage, "Filter is invalid. principalId")))
			if !missingWrite {
				return false, &crud.FatalPollError{Err: deleteErr}
			}
		}
		remote, getErr := r.client.Oauth2PermissionGrants().ByOAuth2PermissionGrantId(object.ID.ValueString()).Get(ctx, nil)
		if getErr != nil {
			info := grapherrors.GraphError(ctx, getErr)
			if info.StatusCode == 400 && strings.Contains(info.ErrorMessage, "Filter is invalid. principalId") {
				missing = 0
				return false, fmt.Errorf("waiting for user consent deletion replication: %w", getErr)
			}
			if grapherrors.GraphError(ctx, getErr).StatusCode == 404 {
				missing++
				return missing >= 3, nil
			}
			return false, &crud.FatalPollError{Err: getErr}
		}
		missing = 0
		if remote == nil {
			return false, &crud.FatalPollError{Err: errGrantResponse}
		}
		tflog.Debug(ctx, "Delegated permission grant is still visible after delete", map[string]any{"id": object.ID.ValueString()})
		return false, nil
	})
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, fmt.Errorf("confirming delegated permission grant deletion: %w", err), resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
