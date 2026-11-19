package graphBetaServicePrincipalSynchronizationJob

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for service principal synchronization job resources.
//
// Operation: Creates a disabled provisioning job and starts it only when enabled is true
// API Calls:
//   - POST /servicePrincipals/{servicePrincipalId}/synchronization/jobs
//   - POST /servicePrincipals/{servicePrincipalId}/synchronization/jobs/{jobId}/start (only when enabled)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/synchronization-synchronizationjob?view=graph-rest-beta
func (r *ServicePrincipalSynchronizationJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServicePrincipalSynchronizationJobResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	// Never blindly replay an ambiguous POST: Graph does not promise idempotency.
	remote, err := r.client.ServicePrincipals().ByServicePrincipalId(data.ServicePrincipalID.ValueString()).Synchronization().Jobs().Post(ctx, constructResource(&data), nil)
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" {
		resp.Diagnostics.AddError("Invalid synchronization job create response", errJobResponse.Error())
		return
	}
	data.JobID = types.StringValue(*remote.GetId())
	data.ID = types.StringValue(data.ServicePrincipalID.ValueString() + "/" + data.JobID.ValueString())
	// Record ownership before enabling or reading so failures retain a recoverable job.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if data.Enabled.ValueBool() {
		err = r.client.ServicePrincipals().
			ByServicePrincipalId(data.ServicePrincipalID.ValueString()).
			Synchronization().
			Jobs().
			BySynchronizationJobId(data.JobID.ValueString()).
			Start().
			Post(ctx, nil)
		if err != nil {
			grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = jobConsistencyPredicate(&data)
	if err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading synchronization job after create", err.Error())
	}
}

// Read handles the Read operation for service principal synchronization job resources.
//
// Operation: Retrieves the service principal synchronization job by ID
// API Calls:
//   - GET /servicePrincipals/{servicePrincipalId}/synchronization/jobs/{jobId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/synchronization-synchronizationjob?view=graph-rest-beta
func (r *ServicePrincipalSynchronizationJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServicePrincipalSynchronizationJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	parent, job, err := parseResourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid synchronization job state ID", err.Error())
		return
	}
	remote, err := r.client.ServicePrincipals().ByServicePrincipalId(parent).Synchronization().Jobs().BySynchronizationJobId(job).Get(ctx, nil)
	if err != nil {
		operation := constants.TfOperationRead
		if op, ok := ctx.Value("retry_operation").(string); ok {
			operation = op
		}
		grapherrors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions, grapherrors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}
	data.ServicePrincipalID = types.StringValue(parent)
	if err = MapRemoteStateToTerraform(&data, remote); err != nil {
		resp.Diagnostics.AddError("Invalid synchronization job response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: data.ID.ValueString()})...)
	}
}

// Update handles the Update operation for service principal synchronization job resources.
//
// Operation: Starts or pauses the job only when enabled changes; parent and template changes require replacement
// API Calls:
//   - POST /servicePrincipals/{servicePrincipalId}/synchronization/jobs/{jobId}/start
//   - POST /servicePrincipals/{servicePrincipalId}/synchronization/jobs/{jobId}/pause
//
// Reference: https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationjob-start?view=graph-rest-beta
func (r *ServicePrincipalSynchronizationJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ServicePrincipalSynchronizationJobResourceModel
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
	parent, job, err := parseResourceID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid synchronization job state ID", err.Error())
		return
	}
	if !plan.Enabled.Equal(state.Enabled) {
		if plan.Enabled.ValueBool() {
			err = r.client.ServicePrincipals().ByServicePrincipalId(parent).Synchronization().Jobs().BySynchronizationJobId(job).Start().Post(ctx, nil)
		} else {
			err = r.client.ServicePrincipals().ByServicePrincipalId(parent).Synchronization().Jobs().BySynchronizationJobId(job).Pause().Post(ctx, nil)
		}
		if err != nil {
			grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}
	plan.ID = state.ID
	plan.JobID = state.JobID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = jobConsistencyPredicate(&plan)
	if err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading synchronization job after update", err.Error())
	}
}

// Delete handles the Delete operation for service principal synchronization job resources.
//
// Operation: Deletes the service principal synchronization job
// API Calls:
//   - DELETE /servicePrincipals/{servicePrincipalId}/synchronization/jobs/{jobId}
//   - GET /servicePrincipals/{servicePrincipalId}/synchronization/jobs/{jobId} (verifies deletion)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/synchronization-synchronizationjob?view=graph-rest-beta
func (r *ServicePrincipalSynchronizationJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ServicePrincipalSynchronizationJobResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	parent, job, err := parseResourceID(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid synchronization job state ID", err.Error())
		return
	}
	missing := 0
	err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		deleteErr := r.client.ServicePrincipals().ByServicePrincipalId(parent).Synchronization().Jobs().BySynchronizationJobId(job).Delete(ctx, nil)
		if deleteErr != nil && grapherrors.GraphError(ctx, deleteErr).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: deleteErr}
		}
		remote, getErr := r.client.ServicePrincipals().ByServicePrincipalId(parent).Synchronization().Jobs().BySynchronizationJobId(job).Get(ctx, nil)
		if getErr != nil {
			if grapherrors.GraphError(ctx, getErr).StatusCode == 404 {
				missing++
				return missing >= 3, nil
			}
			return false, &crud.FatalPollError{Err: getErr}
		}
		missing = 0
		if remote == nil {
			return false, &crud.FatalPollError{Err: errJobResponse}
		}
		return false, nil
	})
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
