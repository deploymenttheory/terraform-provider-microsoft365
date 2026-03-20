package graphBetaWindowsUpdatesAutopatchUpdatePolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for Windows Update policy resources.
//
// Operation: Creates a new Windows Update policy
// API Calls:
//   - POST /admin/windows/updates/updatePolicies
//
// Reference: https://learn.microsoft.com/en-us/graph/api/adminwindowsupdates-post-updatepolicies?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsUpdatesAutopatchUpdatePolicyResourceModel

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

	requestBody, err := constructResource(ctx, &object, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatePolicies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

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

// Read handles the Read operation for Windows Update policy resources.
//
// Operation: Retrieves a Windows Update policy by ID
// API Calls:
//   - GET /admin/windows/updates/updatePolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatepolicy-get?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsUpdatesAutopatchUpdatePolicyResourceModel
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
		UpdatePolicies().
		ByUpdatePolicyId(object.ID.ValueString()).
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

// Update handles the Update operation for Windows Update policy resources.
//
// Operation: Updates an existing Windows Update policy
// API Calls:
//   - PATCH /admin/windows/updates/updatePolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatepolicy-update?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsUpdatesAutopatchUpdatePolicyResourceModel
	var state WindowsUpdatesAutopatchUpdatePolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Admin().
		Windows().
		Updates().
		UpdatePolicies().
		ByUpdatePolicyId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Windows Update policy resources.
//
// Operation: Deletes a Windows Update policy
// API Calls:
//   - DELETE /admin/windows/updates/updatePolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatepolicy-delete?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchUpdatePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsUpdatesAutopatchUpdatePolicyResourceModel
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
		tflog.Debug(ctx, fmt.Sprintf("Delete attempt %d for update policy %s", attempt, object.ID.ValueString()))

		err := r.client.
			Admin().
			Windows().
			Updates().
			UpdatePolicies().
			ByUpdatePolicyId(object.ID.ValueString()).
			Delete(ctx, nil)

		if err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Delete API call succeeded for update policy %s", object.ID.ValueString()))
			break
		}

		errorInfo := errors.GraphError(ctx, err)

		if errors.IsRetryableDeleteError(&errorInfo) {
			tflog.Debug(ctx, fmt.Sprintf("Retryable delete error (attempt %d, status %d), waiting %s before retry",
				attempt, errorInfo.StatusCode, wait))

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

		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	verifyAttempt := 0
	verifyWait := 2 * time.Second

verifyLoop:
	for {
		verifyAttempt++
		tflog.Debug(ctx, fmt.Sprintf("Verification attempt %d: checking if update policy %s is deleted", verifyAttempt, object.ID.ValueString()))

		_, getErr := r.client.
			Admin().
			Windows().
			Updates().
			UpdatePolicies().
			ByUpdatePolicyId(object.ID.ValueString()).
			Get(ctx, nil)

		if getErr != nil {
			errorInfo := errors.GraphError(ctx, getErr)
			if errorInfo.StatusCode == 404 {
				tflog.Debug(ctx, fmt.Sprintf("Update policy %s confirmed deleted (404)", object.ID.ValueString()))
				break
			}

			tflog.Debug(ctx, fmt.Sprintf("Error verifying deletion (attempt %d, status %d): %s", verifyAttempt, errorInfo.StatusCode, errorInfo.ErrorMessage))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Update policy %s still exists (attempt %d)", object.ID.ValueString(), verifyAttempt))
		}

		select {
		case <-time.After(verifyWait):
		case <-ctx.Done():
			tflog.Warn(ctx, fmt.Sprintf("Timeout waiting for update policy deletion confirmation after %d attempts: %s", verifyAttempt, ctx.Err()))
			break verifyLoop
		}

		if verifyWait*2 <= maxWait {
			verifyWait *= 2
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
