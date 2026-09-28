package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

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

// Create handles the Create operation for B2B Management Policy Assignment resources.
//
// Operation: Applies a B2B management policy to an application or service principal
// API Calls:
//   - GET /policies/b2bManagementPolicies/{b2bManagementPolicyId}
//   - GET /directoryObjects/{directoryObjectId}
//   - POST /{applications|servicePrincipals}/{directoryObjectId}/policies/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-post-appliesto?view=graph-rest-beta
func (r *B2bManagementPolicyAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object B2bManagementPolicyAssignmentResourceModel

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

	policyID := object.B2bManagementPolicyID.ValueString()
	directoryObjectID := object.DirectoryObjectID.ValueString()

	// A policy, application or service principal created moments earlier may not have propagated
	// across Microsoft Entra replicas yet. The $ref POST is not idempotent, so it is never retried;
	// instead, wait until both ends of the reference are readable (idempotent GETs) first.
	if err := r.waitForPolicyPropagation(ctx, policyID); err != nil {
		resp.Diagnostics.AddError(
			"Error verifying B2B management policy before assignment",
			fmt.Sprintf("B2B management policy %s could not be read prior to assigning it to directory object %s: %s", policyID, directoryObjectID, err.Error()),
		)
		return
	}

	directoryObjectType, err := r.waitForDirectoryObjectType(ctx, directoryObjectID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving directory object before assignment",
			fmt.Sprintf("Directory object %s could not be resolved to an application or service principal: %s", directoryObjectID, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Applying B2B management policy %s to %s %s", policyID, directoryObjectType, directoryObjectID))

	if err := r.addPolicyReference(ctx, directoryObjectType, directoryObjectID, policyID); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = compositeID(policyID, directoryObjectID)
	object.DirectoryObjectType = types.StringValue(directoryObjectType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = b2bManagementPolicyAssignmentConsistencyPredicate(&object)

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

// waitForPolicyPropagation polls GET /policies/b2bManagementPolicies/{id} until the policy is
// visible, treating 404 as Microsoft Entra replication lag.
func (r *B2bManagementPolicyAssignmentResource) waitForPolicyPropagation(ctx context.Context, policyID string) error {
	return crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		_, err := r.client.
			Policies().
			B2bManagementPolicies().
			ByB2bManagementPolicyId(policyID).
			Get(ctx, nil)

		if err == nil {
			return true, nil
		}

		if errors.GraphError(ctx, err).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: err}
		}

		tflog.Debug(ctx, fmt.Sprintf("B2B management policy %s not visible yet, awaiting Entra propagation", policyID))
		return false, fmt.Errorf("B2B management policy %s is not visible in the directory yet: %w", policyID, err)
	})
}

// waitForDirectoryObjectType polls GET /directoryObjects/{id} until the object is visible and
// returns whether it is an application or a service principal. 404 is treated as replication lag;
// any other type of directory object is a fatal error.
func (r *B2bManagementPolicyAssignmentResource) waitForDirectoryObjectType(ctx context.Context, directoryObjectID string) (string, error) {
	var directoryObjectType string

	err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		resolvedType, err := r.resolveDirectoryObjectType(ctx, directoryObjectID)
		if err == nil {
			directoryObjectType = resolvedType
			return true, nil
		}

		if errors.GraphError(ctx, err).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: err}
		}

		tflog.Debug(ctx, fmt.Sprintf("Directory object %s not visible yet, awaiting Entra propagation", directoryObjectID))
		return false, fmt.Errorf("directory object %s is not visible in the directory yet: %w", directoryObjectID, err)
	})

	return directoryObjectType, err
}

// Read handles the Read operation for B2B Management Policy Assignment resources.
//
// Operation: Verifies the B2B management policy applies to the directory object
// API Calls:
//   - GET /policies/b2bManagementPolicies/{b2bManagementPolicyId}/appliesTo
//
// Reference: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-list-appliesto?view=graph-rest-beta
func (r *B2bManagementPolicyAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object B2bManagementPolicyAssignmentResourceModel
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

	policyID := object.B2bManagementPolicyID.ValueString()
	directoryObjectID := object.DirectoryObjectID.ValueString()

	directoryObjects, err := r.listAppliesTo(ctx, policyID)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	for _, directoryObject := range directoryObjects {
		if directoryObject.GetId() == nil || *directoryObject.GetId() != directoryObjectID {
			continue
		}

		if directoryObjectType, ok := directoryObjectTypeFromODataType(directoryObject.GetOdataType()); ok {
			object.DirectoryObjectType = types.StringValue(directoryObjectType)
		}
		object.ID = compositeID(policyID, directoryObjectID)

		resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
		return
	}

	tflog.Debug(ctx, "B2B management policy assignment not found, removing from state", map[string]any{
		"b2b_management_policy_id": policyID,
		"directory_object_id":      directoryObjectID,
	})
	resp.State.RemoveResource(ctx)
}

// Update handles the Update operation for B2B Management Policy Assignment resources.
//
// Operation: Since both identifying fields have RequiresReplace, this is effectively a no-op
// (terraform will destroy and recreate). This Update implementation handles the edge case where
// only timeout changes occur.
func (r *B2bManagementPolicyAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan B2bManagementPolicyAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles the Delete operation for B2B Management Policy Assignment resources.
//
// Operation: Removes a B2B management policy from an application or service principal
// API Calls:
//   - DELETE /{applications|servicePrincipals}/{directoryObjectId}/policies/{b2bManagementPolicyId}/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-delete-appliesto?view=graph-rest-beta
func (r *B2bManagementPolicyAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object B2bManagementPolicyAssignmentResourceModel

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

	policyID := object.B2bManagementPolicyID.ValueString()
	directoryObjectID := object.DirectoryObjectID.ValueString()

	directoryObjectType := object.DirectoryObjectType.ValueString()
	if directoryObjectType == "" {
		resolvedType, err := r.resolveDirectoryObjectType(ctx, directoryObjectID)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}
		directoryObjectType = resolvedType
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing B2B management policy %s from %s %s", policyID, directoryObjectType, directoryObjectID))

	if err := r.removePolicyReference(ctx, directoryObjectType, directoryObjectID, policyID); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
