package graphBetaWindowsUpdatesAutopatchDeploymentState

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for Windows Updates autopatch deployment state resources.
//
// Operation: Sets the state of a deployment (e.g., scheduled, offering, paused)
// API Calls:
//   - PATCH /admin/windows/updates/deployments/{deploymentId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-update?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeploymentStateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsUpdatesAutopatchDeploymentStateResourceModel

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

	requestBody, err := constructStateResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Admin().
		Windows().
		Updates().
		Deployments().
		ByDeploymentId(object.DeploymentId.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = object.DeploymentId

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	time.Sleep(10 * time.Second)

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
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

// Read handles the Read operation for Windows Updates autopatch deployment state resources.
//
// Operation: Retrieves the current state of a deployment
// API Calls:
//   - GET /admin/windows/updates/deployments/{deploymentId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-get?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeploymentStateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsUpdatesAutopatchDeploymentStateResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read of resource: %s", ResourceName))

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with deployment ID: %s", ResourceName, object.DeploymentId.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	remoteResource, err := r.client.
		Admin().
		Windows().
		Updates().
		Deployments().
		ByDeploymentId(object.DeploymentId.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, remoteResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Updates autopatch deployment state resources.
//
// Operation: Updates the state of a deployment (e.g., transitions between scheduled, offering, paused)
// API Calls:
//   - PATCH /admin/windows/updates/deployments/{deploymentId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-update?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeploymentStateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsUpdatesAutopatchDeploymentStateResourceModel
	var state WindowsUpdatesAutopatchDeploymentStateResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s for deployment ID: %s", ResourceName, state.DeploymentId.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructStateResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing update request",
			fmt.Sprintf("Could not construct update request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Admin().
		Windows().
		Updates().
		Deployments().
		ByDeploymentId(state.DeploymentId.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	time.Sleep(10 * time.Second)

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Windows Updates autopatch deployment state resources.
//
// Operation: Resets the deployment state to default by patching the deployment
// API Calls:
//   - PATCH /admin/windows/updates/deployments/{deploymentId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-deployment-update?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeploymentStateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsUpdatesAutopatchDeploymentStateResourceModel

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

	requestBody, err := constructResetStateResource(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing delete request",
			fmt.Sprintf("Could not construct delete request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Admin().
		Windows().
		Updates().
		Deployments().
		ByDeploymentId(object.DeploymentId.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
