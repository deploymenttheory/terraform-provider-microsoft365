package graphBetaWindowsQualityUpdatePolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the Create operation.
func (r *WindowsQualityUpdatePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsQualityUpdatePolicyResourceModel

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

	createdResource, err := r.client.
		DeviceManagement().
		WindowsQualityUpdatePolicies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		WindowsQualityUpdatePolicies().
		ByWindowsQualityUpdatePolicyId(object.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully posted assignments for policy ID: %s", object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
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

// Read handles the Read operation for macos platform scripts resources.
//
//   - Retrieves the current state from the read request
//   - Gets the resource details including assignments from the API using expand
//   - Maps both resource and assignment details to Terraform state
//
// The function ensures all components are properly read and mapped into the
// Terraform state in a single API call, providing a complete view of the
// resource's current configuration on the server.
// Read handles the Read operation for Windows Quality Update Policies.
func (r *WindowsQualityUpdatePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsQualityUpdatePolicyResourceModel

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

	if object.ID.IsNull() || object.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing ID",
			"Cannot read resource without an ID.",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resourceResp, err := r.client.
		DeviceManagement().
		WindowsQualityUpdatePolicies().
		ByWindowsQualityUpdatePolicyId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.WindowsQualityUpdatePoliciesWindowsQualityUpdatePolicyItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.WindowsQualityUpdatePoliciesWindowsQualityUpdatePolicyItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, resourceResp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Quality Update Policy resources.
func (r *WindowsQualityUpdatePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsQualityUpdatePolicyResourceModel
	var state WindowsQualityUpdatePolicyResourceModel

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
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		WindowsQualityUpdatePolicies().
		ByWindowsQualityUpdatePolicyId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	if !plan.Assignments.IsNull() && !plan.Assignments.IsUnknown() {
		tflog.Debug(ctx, fmt.Sprintf("Assignments detected, constructing assignment request for policy ID: %s", state.ID.ValueString()))

		assignRequestBody, err := constructAssignment(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Failed to construct assignments for policy: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			WindowsQualityUpdatePolicies().
			ByWindowsQualityUpdatePolicyId(state.ID.ValueString()).
			Assign().
			Post(ctx, assignRequestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully posted assignments for policy ID: %s", state.ID.ValueString()))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
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

// Delete handles the Delete operation for windows driver update profile resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *WindowsQualityUpdatePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsQualityUpdatePolicyResourceModel

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

	err := r.client.
		DeviceManagement().
		WindowsQualityUpdatePolicies().
		ByWindowsQualityUpdatePolicyId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
