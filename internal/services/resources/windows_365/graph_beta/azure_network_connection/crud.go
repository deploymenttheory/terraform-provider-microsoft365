package graphBetaAzureNetworkConnection

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
func (r *CloudPcOnPremisesConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudPcOnPremisesConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing resource", fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()))
		return
	}

	created, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		OnPremisesConnections().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, created, plan.AdDomainPassword)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

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
func (r *CloudPcOnPremisesConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object CloudPcOnPremisesConnectionResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s (operation: %s)", ResourceName, object.ID.ValueString(), operation))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}

	defer cancel()
	item, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		OnPremisesConnections().
		ByCloudPcOnPremisesConnectionId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, item, object.AdDomainPassword)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *CloudPcOnPremisesConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudPcOnPremisesConnectionResourceModel
	var state CloudPcOnPremisesConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

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
		resp.Diagnostics.AddError("Error constructing resource for update method", fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()))
		return
	}

	updated, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		OnPremisesConnections().
		ByCloudPcOnPremisesConnectionId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, updated, plan.AdDomainPassword)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

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
func (r *CloudPcOnPremisesConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object CloudPcOnPremisesConnectionResourceModel

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

	err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		OnPremisesConnections().
		ByCloudPcOnPremisesConnectionId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
