package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
func (r *ConditionalAccessPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ConditionalAccessPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing conditional access policy",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	conditionalAccessPolicy, err := r.client.Identity().ConditionalAccess().Policies().Post(ctx, requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating conditional access policy",
			fmt.Sprintf("Could not create conditional access policy: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(*conditionalAccessPolicy.GetId())

	MapRemoteStateToTerraform(ctx, &plan, conditionalAccessPolicy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *ConditionalAccessPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ConditionalAccessPolicyResourceModel
	tflog.Debug(ctx, "Starting Read method for conditional access policy")

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading conditional access policy with ID: %s", state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	conditionalAccessPolicy, err := r.client.Identity().ConditionalAccess().Policies().ByConditionalAccessPolicyId(state.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		crud.HandleReadErrorIfNotFound(ctx, resp, r, &state, err)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, conditionalAccessPolicy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *ConditionalAccessPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ConditionalAccessPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing conditional access policy",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.Identity().ConditionalAccess().Policies().ByConditionalAccessPolicyId(plan.ID.ValueString()).Patch(ctx, requestBody, nil)
	if err != nil {
		crud.HandleUpdateErrorIfNotFound(ctx, resp, r, &plan, err)
		return
	}

	updatedPolicy, err := r.client.Identity().ConditionalAccess().Policies().ByConditionalAccessPolicyId(plan.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Updated Conditional Access Policy",
			fmt.Sprintf("Could not read updated conditional access policy: %s", err.Error()),
		)
		return
	}

	// Map the updated policy back to the Terraform state
	MapRemoteStateToTerraform(ctx, &plan, updatedPolicy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *ConditionalAccessPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ConditionalAccessPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.Identity().ConditionalAccess().Policies().ByConditionalAccessPolicyId(data.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting %s_%s", r.ProviderTypeName, r.TypeName),
			fmt.Sprintf("Failed to delete conditional access policy: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
