package graphBetaApplicationOwner

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

// Create handles the Create operation for Application Owner resources.
//
// Operation: Adds an owner to an application
// API Calls:
//   - POST /applications/{applicationId}/owners/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-post-owners?view=graph-rest-beta
func (r *ApplicationOwnerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ApplicationOwnerResourceModel

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

	applicationId := object.ApplicationID.ValueString()
	ownerId := object.OwnerID.ValueString()
	ownerObjectType := object.OwnerObjectType.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Adding owner of type %s to application %s", ownerObjectType, applicationId))

	requestBody, err := constructResource(ctx, &object, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		Applications().
		ByApplicationId(applicationId).
		Owners().
		Ref().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// Create composite ID since Microsoft Graph API doesn't return a unique assignment ID
	// Application owner assignments are just relationships, not objects with their own IDs
	// We construct a composite ID from application_id/owner_id to uniquely identify this relationship
	compositeID := fmt.Sprintf("%s/%s", applicationId, ownerId)
	object.ID = types.StringValue(compositeID)

	object.OwnerType = types.StringValue("Unknown") // Will be updated in the read operation
	object.OwnerDisplayName = types.StringValue("")

	tflog.Debug(ctx, "Waiting 5 seconds for eventual consistency after create")
	time.Sleep(5 * time.Second)

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

// Read handles the Read operation for Application Owner resources.
//
// Operation: Retrieves application owners to verify ownership exists
// API Calls:
//   - GET /applications/{applicationId}/owners
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-list-owners?view=graph-rest-beta
func (r *ApplicationOwnerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ApplicationOwnerResourceModel

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

	applicationId := object.ApplicationID.ValueString()
	ownerId := object.OwnerID.ValueString()

	owners, err := r.client.
		Applications().
		ByApplicationId(applicationId).
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

	if ownerObject == nil {
		if operation == constants.TfOperationRead {
			tflog.Warn(ctx, fmt.Sprintf("Owner %s not found on application, removing from state", ownerId))
			resp.State.RemoveResource(ctx)
			return
		}
		// During Create/Update retry, return error to allow retry
		resp.Diagnostics.AddError(
			"Owner not found",
			fmt.Sprintf("Owner %s not yet available on application, retry may be needed", ownerId),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, ownerObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Application Owner resources.
//
// Operation: Updates application ownership (removes old and creates new if application or owner changes)
// API Calls:
//   - DELETE /applications/{applicationId}/owners/{directoryObjectId}/$ref (if application_id or owner_id changes)
//   - POST /applications/{applicationId}/owners/$ref (if application_id or owner_id changes)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-delete-owners?view=graph-rest-beta
// Note: Application owner assignments are relationships without their own IDs; changes require deletion and recreation
func (r *ApplicationOwnerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state ApplicationOwnerResourceModel

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

	// For application owner assignments, if either application_id or owner_id changes,
	// we need to remove the old assignment and create a new one
	if plan.ApplicationID.ValueString() != state.ApplicationID.ValueString() ||
		plan.OwnerID.ValueString() != state.OwnerID.ValueString() {

		err := r.client.
			Applications().
			ByApplicationId(state.ApplicationID.ValueString()).
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
			Applications().
			ByApplicationId(plan.ApplicationID.ValueString()).
			Owners().
			Ref().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		compositeID := fmt.Sprintf("%s/%s", plan.ApplicationID.ValueString(), plan.OwnerID.ValueString())
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

// Delete handles the Delete operation for Application Owner resources.
//
// Operation: Removes an owner from an application
// API Calls:
//   - DELETE /applications/{applicationId}/owners/{directoryObjectId}/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-delete-owners?view=graph-rest-beta
func (r *ApplicationOwnerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ApplicationOwnerResourceModel

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

	applicationId := object.ApplicationID.ValueString()
	ownerId := object.OwnerID.ValueString()

	err := r.client.
		Applications().
		ByApplicationId(applicationId).
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
