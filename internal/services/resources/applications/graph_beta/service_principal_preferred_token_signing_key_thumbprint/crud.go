package graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for the preferred token signing key thumbprint resource.
//
// Operation: Sets preferredTokenSigningKeyThumbprint on a service principal
// API Calls:
//   - PATCH /servicePrincipals/{servicePrincipalId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-update?view=graph-rest-beta
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel

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

	if err := r.patchThumbprint(ctx, &object, constants.TfOperationCreate, resp, &resp.Diagnostics); err != nil {
		return
	}

	object.Id = object.ServicePrincipalID

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = preferredTokenSigningKeyThumbprintConsistencyPredicate(&object)

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for the preferred token signing key thumbprint resource.
//
// Operation: Retrieves the service principal and maps preferredTokenSigningKeyThumbprint
// API Calls:
//   - GET /servicePrincipals/{servicePrincipalId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-get?view=graph-rest-beta
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel
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

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	identity.ID = object.Id.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	servicePrincipalID := object.ServicePrincipalID.ValueString()

	servicePrincipal, err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(servicePrincipalID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, servicePrincipal)

	// A cleared thumbprint on a normal refresh means the property no longer exists remotely.
	// During create/update retries the null value is kept so the consistency predicate
	// fails and the read is retried until propagation completes.
	if operation == constants.TfOperationRead && (object.Thumbprint.IsNull() || object.Thumbprint.ValueString() == "") {
		tflog.Debug(ctx, fmt.Sprintf("preferredTokenSigningKeyThumbprint is not set on service principal %s, removing %s from state", servicePrincipalID, ResourceName))
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for the preferred token signing key thumbprint resource.
//
// Operation: Sets the new thumbprint in place (certificate rotation)
// API Calls:
//   - PATCH /servicePrincipals/{servicePrincipalId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-update?view=graph-rest-beta
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if err := r.patchThumbprint(ctx, &object, constants.TfOperationUpdate, resp, &resp.Diagnostics); err != nil {
		return
	}

	object.Id = object.ServicePrincipalID

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = preferredTokenSigningKeyThumbprintConsistencyPredicate(&object)

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

// Delete handles the Delete operation for the preferred token signing key thumbprint resource.
//
// Operation: Clears preferredTokenSigningKeyThumbprint by sending an explicit JSON null,
// returning the service principal to Microsoft Entra's automatic signing key selection.
// The token signing certificate itself is not removed.
// API Calls:
//   - PATCH /servicePrincipals/{servicePrincipalId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-update?view=graph-rest-beta
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel

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

	servicePrincipalID := object.ServicePrincipalID.ValueString()

	requestBody, err := constructDeleteResource(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing delete request",
			fmt.Sprintf("Could not construct request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		ServicePrincipals().
		ByServicePrincipalId(servicePrincipalID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully cleared preferredTokenSigningKeyThumbprint on service principal: %s", servicePrincipalID))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// patchThumbprint constructs and sends the PATCH request that sets the thumbprint,
// reporting failures to the caller's Create or Update response.
func (r *ServicePrincipalPreferredTokenSigningKeyThumbprintResource) patchThumbprint(ctx context.Context, object *ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel, operation string, resp any, diagnostics *diag.Diagnostics) error {
	requestBody, err := constructResource(ctx, object)
	if err != nil {
		diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s", err.Error()),
		)
		return err
	}

	servicePrincipalID := object.ServicePrincipalID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Setting preferredTokenSigningKeyThumbprint on service principal: %s", servicePrincipalID))

	_, err = r.client.
		ServicePrincipals().
		ByServicePrincipalId(servicePrincipalID).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.WritePermissions)
		return err
	}

	return nil
}
