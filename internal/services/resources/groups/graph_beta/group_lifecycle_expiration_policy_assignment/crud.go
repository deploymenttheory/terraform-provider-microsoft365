package graphBetaGroupLifecycleExpirationPolicyAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for assigning a group to the lifecycle expiration policy.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupLifecycleExpirationPolicyAssignmentResourceModel

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

	groupID := object.GroupID.ValueString()

	// Construct the request for adding a group to the lifecycle expiration policy
	policyID, requestBody, err := r.constructAddGroupRequest(ctx, groupID, &resp.Diagnostics)
	if err != nil || policyID == "" {
		return
	}

	// Call the addGroup endpoint
	tflog.Debug(ctx, fmt.Sprintf("Adding group %s to lifecycle expiration policy %s", groupID, policyID))

	result, err := r.client.
		GroupLifecyclePolicies().
		ByGroupLifecyclePolicyId(policyID).
		AddGroup().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	// Check the result - API returns true if successful, false otherwise
	if result != nil && result.GetValue() != nil {
		success := *result.GetValue()
		if !success {
			resp.Diagnostics.AddError(
				"Failed to add group to lifecycle expiration policy",
				fmt.Sprintf("The API returned false when attempting to add group %s to the lifecycle expiration policy. "+
					"This may occur if the group is not a Microsoft 365 group or if the policy's managedGroupTypes is not set to 'Selected'.", groupID),
			)
			return
		}
	}

	object.ID = object.GroupID

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
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

// Read handles the Read operation.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupLifecycleExpirationPolicyAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := "Read"
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

	groupID := object.GroupID.ValueString()

	policies, err := r.client.
		GroupLifecyclePolicies().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if policies == nil || policies.GetValue() == nil || len(policies.GetValue()) == 0 || policies.GetValue()[0].GetId() == nil {
		tflog.Debug(ctx, "No valid lifecycle expiration policy found in tenant - removing assignment from state", map[string]any{
			"groupId": groupID,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	policy := policies.GetValue()[0]
	policyID := policy.GetId()
	managedGroupTypes := policy.GetManagedGroupTypes()

	// Check if policy is configured to manage selected groups
	// If managedGroupTypes is not "Selected", individual assignments don't apply
	if managedGroupTypes == nil || *managedGroupTypes != "Selected" {
		tflog.Debug(ctx, "lifecycle expiration policy managedGroupTypes is not set to 'Selected' - removing assignment from state", map[string]any{
			"groupId":           groupID,
			"policyId":          *policyID,
			"managedGroupTypes": managedGroupTypes,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, groupID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation by removing the old group and adding the new group.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state GroupLifecycleExpirationPolicyAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

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

	oldGroupID := state.GroupID.ValueString()
	newGroupID := plan.GroupID.ValueString()

	// If group_id changed, remove old group and add new group
	if oldGroupID != newGroupID {
		tflog.Debug(ctx, fmt.Sprintf("Group ID changed from %s to %s, performing remove and add operations", oldGroupID, newGroupID))

		oldPolicyID, removeRequestBody, err := r.constructRemoveGroupRequest(ctx, oldGroupID, &resp.Diagnostics)
		if err != nil || oldPolicyID == "" {
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Removing old group %s from lifecycle expiration policy %s", oldGroupID, oldPolicyID))

		removeResult, err := r.client.
			GroupLifecyclePolicies().
			ByGroupLifecyclePolicyId(oldPolicyID).
			RemoveGroup().
			Post(ctx, removeRequestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}

		if removeResult != nil && removeResult.GetValue() != nil {
			success := *removeResult.GetValue()
			if !success {
				tflog.Warn(ctx, fmt.Sprintf("API returned false when removing old group %s - continuing with add operation", oldGroupID))
			}
		}

		newPolicyID, addRequestBody, err := r.constructAddGroupRequest(ctx, newGroupID, &resp.Diagnostics)
		if err != nil || newPolicyID == "" {
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Adding new group %s to lifecycle expiration policy %s", newGroupID, newPolicyID))

		addResult, err := r.client.
			GroupLifecyclePolicies().
			ByGroupLifecyclePolicyId(newPolicyID).
			AddGroup().
			Post(ctx, addRequestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}

		if addResult != nil && addResult.GetValue() != nil {
			success := *addResult.GetValue()
			if !success {
				resp.Diagnostics.AddError(
					"Failed to add new group to lifecycle expiration policy",
					fmt.Sprintf("The API returned false when attempting to add group %s to the lifecycle expiration policy.", newGroupID),
				)
				return
			}
		}

		plan.ID = plan.GroupID
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = ResourceName

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupLifecycleExpirationPolicyAssignmentResourceModel

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

	groupID := object.GroupID.ValueString()

	policyID, requestBody, err := r.constructRemoveGroupRequest(ctx, groupID, &resp.Diagnostics)
	if err != nil || policyID == "" {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing group %s from lifecycle expiration policy %s", groupID, policyID))

	result, err := r.client.
		GroupLifecyclePolicies().
		ByGroupLifecyclePolicyId(policyID).
		RemoveGroup().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	if result != nil && result.GetValue() != nil {
		success := *result.GetValue()
		if !success {
			tflog.Warn(ctx, fmt.Sprintf("API returned false when removing group %s from lifecycle expiration policy - group may already be removed", groupID))
		}
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
