package graphBetaApplicationsServicePrincipalTokenLifetimePolicyAssignment

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
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Service Principal Token Lifetime Policy Assignment resources.
//
// Operation: Assigns a token lifetime policy to a service principal
// API Calls:
//   - POST /servicePrincipals/{servicePrincipalId}/tokenLifetimePolicies/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-tokenlifetimepolicies?view=graph-rest-beta
func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ServicePrincipalTokenLifetimePolicyAssignmentResourceModel

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

	spID := object.ServicePrincipalID.ValueString()
	policyID := object.TokenLifetimePolicyID.ValueString()

	refBody := graphmodels.NewReferenceCreate()
	odataID := fmt.Sprintf("https://graph.microsoft.com/beta/policies/tokenLifetimePolicies/%s", policyID)
	refBody.SetOdataId(&odataID)

	tflog.Debug(ctx, fmt.Sprintf("Assigning token lifetime policy %s to service principal %s", policyID, spID))

	// A token lifetime policy created moments earlier may not have propagated across
	// Microsoft Entra replicas yet, in which case the $ref POST fails with 404
	// ("Unable to read the company information from the directory"). The POST is not
	// idempotent, so it is never retried; instead, wait until the referenced policy is
	// readable (an idempotent GET) before attempting the assignment exactly once.
	if err := r.waitForTokenLifetimePolicyPropagation(ctx, policyID); err != nil {
		resp.Diagnostics.AddError(
			"Error verifying token lifetime policy before assignment",
			fmt.Sprintf("Token lifetime policy %s could not be read prior to assigning it to service principal %s: %s", policyID, spID, err.Error()),
		)
		return
	}

	err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(spID).
		TokenLifetimePolicies().
		Ref().
		Post(ctx, refBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(spID + "/" + policyID)

	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned token lifetime policy %s to service principal %s", policyID, spID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = servicePrincipalTokenLifetimePolicyAssignmentConsistencyPredicate(&object)

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

// waitForTokenLifetimePolicyPropagation polls GET /policies/tokenLifetimePolicies/{id} until
// the referenced policy is visible, treating 404 as Microsoft Entra replication lag. Reads are
// idempotent, so polling is safe; the subsequent non-idempotent $ref POST is performed exactly
// once. A successful read does not strictly guarantee that the replica resolving the $ref POST
// has the policy — Graph exposes no propagation-complete signal — but reading immediately
// before the write is the closest achievable precondition. Polling is bounded by the create
// timeout on ctx; the same permission that authorizes the $ref POST
// (Policy.ReadWrite.ApplicationConfiguration) authorizes this read.
//
// See: https://devblogs.microsoft.com/identity/designing-for-eventual-consistency-for-microsoft-entra/
func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) waitForTokenLifetimePolicyPropagation(ctx context.Context, policyID string) error {
	const pollInterval = 2 * time.Second

	deadline, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		return fmt.Errorf("context must have a deadline")
	}

	for attempt := 1; ; attempt++ {
		_, err := r.client.
			Policies().
			TokenLifetimePolicies().
			ByTokenLifetimePolicyId(policyID).
			Get(ctx, nil)

		if err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Token lifetime policy %s is visible in the directory (attempt %d)", policyID, attempt))
			return nil
		}

		errorInfo := errors.GraphError(ctx, err)
		if errorInfo.StatusCode != 404 {
			return err
		}

		if time.Until(deadline) < pollInterval+time.Second {
			return fmt.Errorf("policy was not visible in the directory before the create timeout (last error: %w)", err)
		}

		tflog.Debug(ctx, fmt.Sprintf("Token lifetime policy %s not visible yet (attempt %d), waiting %s for Entra propagation", policyID, attempt, pollInterval))

		select {
		case <-time.After(pollInterval):
		case <-ctx.Done():
			return fmt.Errorf("context cancelled while waiting for policy propagation: %w", ctx.Err())
		}
	}
}

// Read handles the Read operation for Service Principal Token Lifetime Policy Assignment resources.
//
// Operation: Verifies the token lifetime policy is assigned to the service principal
// API Calls:
//   - GET /servicePrincipals/{servicePrincipalId}/tokenLifetimePolicies
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-list-tokenlifetimepolicies?view=graph-rest-beta
func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ServicePrincipalTokenLifetimePolicyAssignmentResourceModel
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

	spID := object.ServicePrincipalID.ValueString()
	policyID := object.TokenLifetimePolicyID.ValueString()

	policies, err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(spID).
		TokenLifetimePolicies().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if policies == nil || policies.GetValue() == nil {
		tflog.Debug(ctx, "No token lifetime policies found for service principal, removing from state", map[string]any{
			"service_principal_id":     spID,
			"token_lifetime_policy_id": policyID,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	found := false
	for _, policy := range policies.GetValue() {
		if policy.GetId() != nil && *policy.GetId() == policyID {
			found = true
			break
		}
	}

	if !found {
		tflog.Debug(ctx, "Token lifetime policy assignment not found, removing from state", map[string]any{
			"service_principal_id":     spID,
			"token_lifetime_policy_id": policyID,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Service Principal Token Lifetime Policy Assignment resources.
//
// Operation: Since both fields have RequiresReplace, this is effectively a no-op (terraform will destroy and recreate)
// This Update implementation handles the edge case where only timeout changes occur.
func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServicePrincipalTokenLifetimePolicyAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles the Delete operation for Service Principal Token Lifetime Policy Assignment resources.
//
// Operation: Removes a token lifetime policy from a service principal
// API Calls:
//   - DELETE /servicePrincipals/{servicePrincipalId}/tokenLifetimePolicies/{tokenLifetimePolicyId}/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-delete-tokenlifetimepolicies?view=graph-rest-beta
func (r *ServicePrincipalTokenLifetimePolicyAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ServicePrincipalTokenLifetimePolicyAssignmentResourceModel

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

	spID := object.ServicePrincipalID.ValueString()
	policyID := object.TokenLifetimePolicyID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Removing token lifetime policy %s from service principal %s", policyID, spID))

	err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(spID).
		TokenLifetimePolicies().
		ByTokenLifetimePolicyId(policyID).
		Ref().
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
