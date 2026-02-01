package graphBetaServicePrincipalOwner

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Service Principal Owner resources.
//
// Operation: Adds an owner to a service principal
// API Calls:
//   - POST /servicePrincipals/{servicePrincipalId}/owners/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-post-owners?view=graph-rest-beta
func (r *ServicePrincipalOwnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ServicePrincipalOwnerResourceModel

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

	servicePrincipalId := object.ServicePrincipalID.ValueString()
	ownerId := object.OwnerID.ValueString()
	ownerObjectType := object.OwnerObjectType.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Adding owner of type %s to service principal %s", ownerObjectType, servicePrincipalId))

	requestBody, err := constructResource(ctx, &object, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		ServicePrincipals().
		ByServicePrincipalId(servicePrincipalId).
		Owners().
		Ref().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// Create composite ID since Microsoft Graph API doesn't return a unique assignment ID
	// Service principal owner assignments are just relationships, not objects with their own IDs
	// We construct a composite ID from service_principal_id/owner_id to uniquely identify this relationship
	compositeID := fmt.Sprintf("%s/%s", servicePrincipalId, ownerId)
	object.ID = types.StringValue(compositeID)

	object.OwnerType = types.StringValue("Unknown") // Will be updated in the read operation
	object.OwnerDisplayName = types.StringValue("")

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

// Read handles the Read operation for Service Principal Owner resources.
//
// Operation: Retrieves service principal owners to verify ownership exists
// API Calls:
//   - GET /servicePrincipals/{servicePrincipalId}/owners
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-list-owners?view=graph-rest-beta
func (r *ServicePrincipalOwnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ServicePrincipalOwnerResourceModel

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

	servicePrincipalId := object.ServicePrincipalID.ValueString()
	ownerId := object.OwnerID.ValueString()

	owners, err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(servicePrincipalId).
		Owners().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	var ownerObject graphmodels.DirectoryObjectable

	if owners != nil && owners.GetValue() != nil {
		for _, owner := range owners.GetValue() {
			if owner.GetId() != nil && *owner.GetId() == ownerId {
				ownerObject = owner
				break
			}
		}
	}

	MapRemoteStateToTerraform(ctx, &object, ownerObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Service Principal Owner resources.
//
// Operation: Updates service principal ownership (removes old and creates new if service principal or owner changes)
// API Calls:
//   - DELETE /servicePrincipals/{servicePrincipalId}/owners/{directoryObjectId}/$ref (if service_principal_id or owner_id changes)
//   - POST /servicePrincipals/{servicePrincipalId}/owners/$ref (if service_principal_id or owner_id changes)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-delete-owners?view=graph-rest-beta
// Note: Service principal owner assignments are relationships without their own IDs; changes require deletion and recreation
func (r *ServicePrincipalOwnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ServicePrincipalOwnerResourceModel

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

	// For service principal owner assignments, if either service_principal_id or owner_id changes,
	// we need to remove the old assignment and create a new one
	if plan.ServicePrincipalID.ValueString() != state.ServicePrincipalID.ValueString() ||
		plan.OwnerID.ValueString() != state.OwnerID.ValueString() {

		err := r.client.
			ServicePrincipals().
			ByServicePrincipalId(state.ServicePrincipalID.ValueString()).
			Owners().
			ByDirectoryObjectId(state.OwnerID.ValueString()).
			Ref().
			Delete(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		requestBody, err := constructResource(ctx, &plan, r.client)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing resource for Update method",
				fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			ServicePrincipals().
			ByServicePrincipalId(plan.ServicePrincipalID.ValueString()).
			Owners().
			Ref().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		compositeID := fmt.Sprintf("%s/%s", plan.ServicePrincipalID.ValueString(), plan.OwnerID.ValueString())
		plan.ID = types.StringValue(compositeID)
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Service Principal Owner resources.
//
// Operation: Removes an owner from a service principal
// API Calls:
//   - DELETE /servicePrincipals/{servicePrincipalId}/owners/{directoryObjectId}/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-delete-owners?view=graph-rest-beta
func (r *ServicePrincipalOwnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ServicePrincipalOwnerResourceModel

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

	servicePrincipalId := object.ServicePrincipalID.ValueString()
	ownerId := object.OwnerID.ValueString()

	err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(servicePrincipalId).
		Owners().
		ByDirectoryObjectId(ownerId).
		Ref().
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
