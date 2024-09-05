package graphCloudPcDeviceImage

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
func (r *CloudPcDeviceImageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudPcDeviceImageResourceModel

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
			"Error constructing Cloud PC Device Image",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	provisioningPolicy, err := r.client.DeviceManagement().VirtualEndpoint().DeviceImages().Post(ctx, requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Cloud PC Device Image",
			fmt.Sprintf("Could not create Cloud PC Device Image: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(*provisioningPolicy.GetId())

	MapRemoteStateToTerraform(ctx, &plan, provisioningPolicy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *CloudPcDeviceImageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CloudPcDeviceImageResourceModel
	tflog.Debug(ctx, "Starting Read method for Cloud PC Device Image")

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading Cloud PC Device Image with ID: %s", state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	provisioningPolicy, err := r.client.DeviceManagement().VirtualEndpoint().DeviceImages().ByCloudPcDeviceImageId(state.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		crud.HandleReadErrorIfNotFound(ctx, resp, r, &state, err)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, provisioningPolicy)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *CloudPcDeviceImageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudPcDeviceImageResourceModel

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
			"Error constructing Cloud PC Device Image",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.DeviceManagement().VirtualEndpoint().DeviceImages().ByCloudPcDeviceImageId(plan.ID.ValueString()).Patch(ctx, requestBody, nil)
	if err != nil {
		crud.HandleUpdateErrorIfNotFound(ctx, resp, r, &plan, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *CloudPcDeviceImageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudPcDeviceImageResourceModel

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

	err := r.client.DeviceManagement().VirtualEndpoint().DeviceImages().ByCloudPcDeviceImageId(data.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
