package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// Create handles the Create operation for Conditional Access Policy resources.
//
//   - Retrieves the planned configuration from the create request
//   - Constructs the resource request body from the plan
//   - Sends POST request to create the base resource and settings
//   - Captures the new resource ID from the response
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API with retry
//   - Updates the final state with the fresh data from the API
//
// The function ensures that both the Conditional Access Policy profile and its assignments
// (if specified) are created properly. The settings must be defined during creation
// as they are required for a successful deployment, while assignments are optional.
func (r *ConditionalAccessPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ConditionalAccessPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	conditionalAccessPolicy, err := r.client.
		Identity().
		ConditionalAccess().
		Policies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*conditionalAccessPolicy.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Create Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource creation",
			fmt.Sprintf("Failed to verify resource creation: %s", err),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *ConditionalAccessPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ConditionalAccessPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	conditionalAccessPolicy, err := r.client.
		Identity().
		ConditionalAccess().
		Policies().
		ByConditionalAccessPolicyId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, conditionalAccessPolicy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for Conditional Access Policy resources.
//
//   - Retrieves the planned changes from the update request
//   - Constructs the resource request body from the plan
//   - Sends PATCH request to update the base resource and settings
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API with retry
//   - Updates the final state with the fresh data from the API
//
// The function ensures that both the settings and assignments are updated atomically,
// and the final state reflects the actual state of the resource on the server.
func (r *ConditionalAccessPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object ConditionalAccessPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Identity().
		ConditionalAccess().
		Policies().
		ByConditionalAccessPolicyId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Update Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource update",
			fmt.Sprintf("Failed to verify resource update: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *ConditionalAccessPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ConditionalAccessPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		Identity().
		ConditionalAccess().
		Policies().
		ByConditionalAccessPolicyId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
