package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// readbackMatch is what the read-back after a write waits for: a populated id and the written policy.
var readbackMatch = crud.AttributeMatchOptions{Paths: crud.RootPaths("display_name", "definition", "is_organization_default"), RequireKnown: crud.RootPaths("id")}

// policyRetryInterval paces the read-back. The four matching reads it waits for span at least 15 seconds;
// curl probes showed two current reads followed by old values again.
var policyRetryInterval = 5 * time.Second

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
	opts.ConsistencyPredicate = crud.RequireConsecutiveReads(4, crud.AttributesMatchPredicate(resp.State, readbackMatch))

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

	policy, err := r.client.
		Policies().
		B2bManagementPolicies().
		ByB2bManagementPolicyId(object.ID.ValueString()).
		Get(ctx, nil)
	if err != nil && operation == constants.TfOperationRead && errors.GraphError(ctx, err).StatusCode == 404 {
		policy, err = confirmPolicyNotFound(ctx, object.ID.ValueString(), err, func(ctx context.Context) (graphmodels.B2bManagementPolicyable, error) {
			policy, readErr := r.client.
				Policies().
				B2bManagementPolicies().
				ByB2bManagementPolicyId(object.ID.ValueString()).
				Get(ctx, nil)
			if readErr != nil {
				return nil, fmt.Errorf("reading B2B management policy %s: %w", object.ID.ValueString(), readErr)
			}
			return policy, nil
		})
	}
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
	opts.ConsistencyPredicate = crud.RequireConsecutiveReads(4, crud.AttributesMatchPredicate(resp.State, readbackMatch))

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

	id := object.ID.ValueString()
	// Keep the Graph calls in Delete; a replica 404 needs confirmation before state is removed.
	getPolicy := func(ctx context.Context) (graphmodels.B2bManagementPolicyable, error) {
		policy, err := r.client.
			Policies().
			B2bManagementPolicies().
			ByB2bManagementPolicyId(id).
			Get(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("reading B2B management policy %s: %w", id, err)
		}
		return policy, nil
	}

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

		_, getErr := getPolicy(ctx)
		if getErr != nil && errors.GraphError(ctx, getErr).StatusCode == 404 {
			_, getErr = confirmPolicyNotFound(ctx, id, getErr, getPolicy)
		}
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
		err = ctx.Err()
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
