package graphBetaWindowsUpdatesAutopatchDeploymentAudience

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for Windows Update deployment audience resources.
//
// Operation: Creates a new deployment audience
// API Calls:
//   - POST /admin/windows/updates/deploymentAudiences
//
// Reference: https://learn.microsoft.com/en-us/graph/api/adminwindowsupdates-post-deploymentaudiences?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeploymentAudienceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsUpdatesAutopatchDeploymentAudienceResourceModel

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
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		Admin().
		Windows().
		Updates().
		DeploymentAudiences().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName

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

// Read handles the Read operation for Windows Update deployment audience resources.
//
// Operation: Retrieves a deployment audience by ID
// API Calls:
//   - GET /admin/windows/updates/deploymentAudiences/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-deploymentaudience-get?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeploymentAudienceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsUpdatesAutopatchDeploymentAudienceResourceModel
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

	respResource, err := r.client.
		Admin().
		Windows().
		Updates().
		DeploymentAudiences().
		ByDeploymentAudienceId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Update deployment audience resources.
//
// Operation: Updates an existing deployment audience
// API Calls:
//   - N/A
//
// Reference: N/A
func (r *WindowsUpdatesAutopatchDeploymentAudienceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsUpdatesAutopatchDeploymentAudienceResourceModel
	var state WindowsUpdatesAutopatchDeploymentAudienceResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// No update operations - this resource is just a container
	// Members and exclusions are managed by the separate audience_members resource

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Windows Update deployment audience resources.
//
// Operation: Deletes a deployment audience
// API Calls:
//   - DELETE /admin/windows/updates/deploymentAudiences/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-deploymentaudience-delete?view=graph-rest-beta
//
// Note: If the audience is still referenced by deployments (e.g. implicit deployments created
// by a content_approval that have not yet been cleaned up by the Graph API), the delete
// is retried with exponential backoff until the dependency clears or the context deadline
// is exceeded.
func (r *WindowsUpdatesAutopatchDeploymentAudienceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsUpdatesAutopatchDeploymentAudienceResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	const maxWait = 30 * time.Second
	wait := 2 * time.Second
	attempt := 0

	for {
		attempt++
		tflog.Debug(ctx, fmt.Sprintf("Delete attempt %d for audience %s", attempt, object.ID.ValueString()))

		err := r.client.
			Admin().
			Windows().
			Updates().
			DeploymentAudiences().
			ByDeploymentAudienceId(object.ID.ValueString()).
			Delete(ctx, nil)

		if err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Successfully deleted audience %s after %d attempts", object.ID.ValueString(), attempt))
			break
		}

		errorInfo := errors.GraphError(ctx, err)
		
		// Special case: Deployment audience deletion can fail with 400 or 409 when still referenced by deployments
		// This is retryable as the deployments may be in the process of being deleted
		isDeploymentReferenceError := (errorInfo.StatusCode == 400 || errorInfo.StatusCode == 409) &&
			(strings.Contains(errorInfo.ErrorMessage, "being used by deployments") ||
				strings.Contains(errorInfo.ErrorMessage, "referenced by") ||
				strings.Contains(errorInfo.ErrorMessage, "delete deployments before proceeding") ||
				strings.Contains(errorInfo.ErrorCode, "Conflict"))

		if isDeploymentReferenceError {
			tflog.Warn(ctx, fmt.Sprintf("Audience %s still referenced by deployments (attempt %d, status %d, code '%s', msg: '%s'), waiting %s before retry", 
				object.ID.ValueString(), attempt, errorInfo.StatusCode, errorInfo.ErrorCode, errorInfo.ErrorMessage, wait))
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				resp.Diagnostics.AddError(
					"Timeout waiting for deployments to clear",
					fmt.Sprintf("Audience %s is still referenced by deployments after %d attempts and %d seconds: %s\nLast error: %s", 
						object.ID.ValueString(), attempt, DeleteTimeout, ctx.Err(), errorInfo.ErrorMessage),
				)
				return
			}
			if wait*2 <= maxWait {
				wait *= 2
			}
			continue
		}

		if errors.IsRetryableDeleteError(&errorInfo) {
			tflog.Debug(ctx, fmt.Sprintf("Retryable delete error (attempt %d, status %d, code %s), waiting %s before retry",
				attempt, errorInfo.StatusCode, errorInfo.ErrorCode, wait))
			select {
			case <-time.After(wait):
			case <-ctx.Done():
				resp.Diagnostics.AddError(
					"Timeout during delete operation",
					fmt.Sprintf("Delete operation timed out after %d attempts: %s", attempt, ctx.Err()),
				)
				return
			}
			if wait*2 <= maxWait {
				wait *= 2
			}
			continue
		}

		// Non-retryable error
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
