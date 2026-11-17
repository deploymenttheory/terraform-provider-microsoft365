package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for B2B Management Policy resources.
//
// Operation: Creates a new B2B management policy
// API Calls:
//   - POST /policies/b2bManagementPolicies
//
// Reference: https://learn.microsoft.com/en-us/graph/api/policyroot-post-b2bmanagementpolicies?view=graph-rest-beta
func (r *B2bManagementPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object B2bManagementPolicyResourceModel

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

	createdPolicy, err := r.client.
		Policies().
		B2bManagementPolicies().
		Post(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if createdPolicy == nil || createdPolicy.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Microsoft Graph did not return an ID for the created %s", ResourceName),
		)
		return
	}

	object.ID = types.StringValue(*createdPolicy.GetId())

	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.MaxRetries = 60
	opts.RetryInterval = policyRetryInterval
	opts.ConsistencyPredicate = stableConsistencyPredicate(&object)

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

// Read handles the Read operation for B2B Management Policy resources.
//
// Operation: Retrieves a B2B management policy by ID
// API Calls:
//   - GET /policies/b2bManagementPolicies/{b2bManagementPolicyId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-get?view=graph-rest-beta
func (r *B2bManagementPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object B2bManagementPolicyResourceModel
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

	policy, err := r.getPolicy(ctx, object.ID.ValueString(), operation == constants.TfOperationRead)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, policy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for B2B Management Policy resources.
//
// Operation: Updates an existing B2B management policy
// API Calls:
//   - PATCH /policies/b2bManagementPolicies/{b2bManagementPolicyId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-update?view=graph-rest-beta
func (r *B2bManagementPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state B2bManagementPolicyResourceModel

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

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Policies().
		B2bManagementPolicies().
		ByB2bManagementPolicyId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)
	if err != nil {
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

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.MaxRetries = 60
	opts.RetryInterval = policyRetryInterval
	opts.ConsistencyPredicate = stableConsistencyPredicate(&plan)

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

// Delete handles the Delete operation for B2B Management Policy resources.
//
// Operation: Deletes a B2B management policy
// API Calls:
//   - DELETE /policies/b2bManagementPolicies/{b2bManagementPolicyId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/policyroot-delete-b2bmanagementpolicies?view=graph-rest-beta
func (r *B2bManagementPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object B2bManagementPolicyResourceModel

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

	err := r.deletePolicy(ctx, object.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// deletePolicy confirms a DELETE 404 before treating the policy as absent. A stale replica
// can reject DELETE while another still returns the policy; retry instead of leaking it.
func (r *B2bManagementPolicyResource) deletePolicy(ctx context.Context, id string) error {
	err := crud.PollUntil(ctx, policyRetryInterval, func(ctx context.Context) (bool, error) {
		if err := ctx.Err(); err != nil {
			return false, &crud.FatalPollError{Err: err}
		}
		err := r.client.
			Policies().
			B2bManagementPolicies().
			ByB2bManagementPolicyId(id).
			Delete(ctx, nil)
		if err == nil {
			return true, nil
		}
		if errors.GraphError(ctx, err).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: err}
		}

		_, getErr := r.getPolicy(ctx, id, true)
		if getErr == nil {
			return false, fmt.Errorf("B2B management policy %s is still readable after DELETE returned 404: %w", id, err)
		}
		if errors.GraphError(ctx, getErr).StatusCode == 404 {
			return true, nil
		}
		return false, &crud.FatalPollError{Err: getErr}
	})
	// PollUntil wraps the last API error on timeout. Returning that 404 would make the
	// shared delete error handler drop state even though deletion was never confirmed.
	if ctx.Err() != nil {
		return fmt.Errorf("deleting B2B management policy %s: %w", id, ctx.Err())
	}
	if err != nil {
		return fmt.Errorf("deleting B2B management policy %s: %w", id, err)
	}
	return nil
}
