package graphBetaServicePrincipalSynchronizationSchemaJson

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for service principal synchronization schema resources.
//
// Operation: The schema of a synchronization job always exists, so Create replaces it with the configured schema.
// API Calls:
//   - PUT /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
//
// Reference: https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationschema-update?view=graph-rest-beta
func (r *ServicePrincipalSynchronizationSchemaJsonResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ServicePrincipalSynchronizationSchemaJsonResourceModel

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

	err = r.putSchema(ctx, object.ServicePrincipalId.ValueString(), object.JobId.ValueString(), requestBody)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(resourceId(object.ServicePrincipalId.ValueString(), object.JobId.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
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

// Read handles the Read operation for service principal synchronization schema resources.
//
// Operation: Retrieves the schema of a synchronization job
// API Calls:
//   - GET /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
//
// Reference: https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationschema-get?view=graph-rest-beta
func (r *ServicePrincipalSynchronizationSchemaJsonResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ServicePrincipalSynchronizationSchemaJsonResourceModel
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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

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

	remoteResource, err := r.getSchema(ctx, object.ServicePrincipalId.ValueString(), object.JobId.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	object, err = MapRemoteResourceStateToTerraform(ctx, object, remoteResource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error mapping resource state",
			fmt.Sprintf("Could not map resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for service principal synchronization schema resources.
//
// Operation: Replaces the schema of a synchronization job with the configured schema
// API Calls:
//   - PUT /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
//
// Reference: https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationschema-update?view=graph-rest-beta
func (r *ServicePrincipalSynchronizationSchemaJsonResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ServicePrincipalSynchronizationSchemaJsonResourceModel
	var state ServicePrincipalSynchronizationSchemaJsonResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s", ResourceName))

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
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.putSchema(ctx, state.ServicePrincipalId.ValueString(), state.JobId.ValueString(), requestBody)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// resp.State still holds the prior state at this point, so a read against it would drop configuration-only values
	// such as a timeouts block added in this update. Reading against the plan keeps them; the response state is still
	// written only by Read.
	readReq := resource.ReadRequest{State: tfsdk.State{Schema: req.Plan.Schema, Raw: req.Plan.Raw}, ProviderMeta: req.ProviderMeta}
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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for service principal synchronization schema resources.
//
// Operation: The schema of a synchronization job cannot be removed. Delete resets it to the default schema of the
// job's synchronization template and removes the resource from Terraform state. A job or service principal that no
// longer exists (404) is treated as already deleted.
// API Calls:
//   - DELETE /servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
//
// Reference: https://learn.microsoft.com/en-us/graph/api/synchronization-synchronizationschema-delete?view=graph-rest-beta
func (r *ServicePrincipalSynchronizationSchemaJsonResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ServicePrincipalSynchronizationSchemaJsonResourceModel

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

	tflog.Info(ctx, fmt.Sprintf("Resetting %s with ID: %s to the template default", ResourceName, object.ID.ValueString()))

	err := r.schemaRequestBuilder(object.ServicePrincipalId.ValueString(), object.JobId.ValueString()).Delete(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
