package graphBetaGroupLifecycleExpirationPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
// Only 1 of this resource type can exist per tenant. The behavior depends on the overwrite_existing_policy flag:
// - If overwrite_existing_policy = true: Finds the existing tenant policy and overwrites it (PATCH)
// - If overwrite_existing_policy = false (default): Attempts to create a new policy (POST)
func (r *GroupLifecycleExpirationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupLifecycleExpirationPolicyResourceModel
	var policyObject graphmodels.GroupLifecyclePolicyable

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
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	if !object.OverwriteExistingPolicy.IsNull() && object.OverwriteExistingPolicy.ValueBool() {
		tflog.Info(ctx, "overwrite_existing_policy is true, attempting to find and update existing policy")

		existingPolicies, err := r.client.
			GroupLifecyclePolicies().
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Create - Get Existing Policy", r.ReadPermissions)
			return
		}

		var existingPolicyID *string
		if existingPolicies != nil && existingPolicies.GetValue() != nil {
			policies := existingPolicies.GetValue()
			if len(policies) > 0 {
				existingPolicyID = policies[0].GetId()
				tflog.Info(ctx, fmt.Sprintf("Found existing lifecycle policy with ID: %s. Overwriting with Terraform configuration.", *existingPolicyID))
			}
		}

		if existingPolicyID != nil {
			tflog.Debug(ctx, "Overwriting existing lifecycle policy")
			policyObject, err = r.client.
				GroupLifecyclePolicies().
				ByGroupLifecyclePolicyId(*existingPolicyID).
				Patch(ctx, requestBody, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, "Create - Overwrite Existing Policy", r.WritePermissions)
				return
			}
		} else {
			resp.Diagnostics.AddError(
				"No existing policy found",
				"overwrite_existing_policy is set to true, but no existing policy was found in the tenant to overwrite.",
			)
			return
		}
	} else {
		// Default behavior: attempt to create new policy
		tflog.Debug(ctx, "Creating new lifecycle policy")
		policyObject, err = r.client.
			GroupLifecyclePolicies().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

	object.ID = types.StringValue(*policyObject.GetId())

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
func (r *GroupLifecycleExpirationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupLifecycleExpirationPolicyResourceModel

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

	policyId := object.ID.ValueString()

	policyObject, err := r.client.
		GroupLifecyclePolicies().
		ByGroupLifecyclePolicyId(policyId).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, policyObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *GroupLifecycleExpirationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupLifecycleExpirationPolicyResourceModel
	var state GroupLifecycleExpirationPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

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

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	policyId := state.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("[DEBUG] Using resource ID for Update (PATCH): %s", policyId))

	_, err = r.client.
		GroupLifecyclePolicies().
		ByGroupLifecyclePolicyId(policyId).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = ResourceName

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

// Delete handles the Delete operation.
func (r *GroupLifecycleExpirationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupLifecycleExpirationPolicyResourceModel

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

	policyId := object.ID.ValueString()

	err := r.client.
		GroupLifecyclePolicies().
		ByGroupLifecyclePolicyId(policyId).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
