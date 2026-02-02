package graphBetaApplication

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/directory"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for application resources.
//
// Operation: Creates a new Microsoft Entra application
// API Calls:
//   - POST /applications
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-post-applications?view=graph-rest-beta
func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ApplicationResourceModel

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

	requestBody, err := constructResource(ctx, r.client, &object, "", true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing application",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdApplication, err := r.client.
		Applications().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.Id = types.StringValue(*createdApplication.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Waiting 25 seconds for eventual consistency after create")
	time.Sleep(25 * time.Second)

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

// Read handles the Read operation for application resources.
//
// Operation: Retrieves application details including owners
// API Calls:
//   - GET /applications/{id}?$expand=*
//   - GET /applications/{id}/owners
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ApplicationResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with Id: %s", ResourceName, object.Id.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	application, err := r.client.
		Applications().
		ByApplicationId(object.Id.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, application)

	owners, err := r.client.
		Applications().
		ByApplicationId(object.Id.ValueString()).
		Owners().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteOwnersToTerraform(ctx, &object, owners)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for application resources.
//
// Operation: Updates application properties and manages owners
// API Calls:
//   - PATCH /applications/{id} (for general properties)
//   - PATCH /applications/{id} (separate call for app_roles if changed)
//   - DELETE /applications/{id}/owners/{ownerId}/$ref (for removed owners)
//   - POST /applications/{id}/owners/$ref (for new owners)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ApplicationResourceModel
	var state ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting update of resource: %s with Id: %s", ResourceName, state.Id.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, r.client, &plan, state.Id.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing application",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Applications().
		ByApplicationId(state.Id.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// App roles require special handling for update operations: roles being removed must first be disabled
	if !plan.AppRoles.Equal(state.AppRoles) {
		tflog.Debug(ctx, fmt.Sprintf("App roles have changed for %s with ID: %s, performing app role update", ResourceName, state.Id.ValueString()))

		var stateRoles, planRoles []ApplicationAppRole

		state.AppRoles.ElementsAs(ctx, &stateRoles, false)
		plan.AppRoles.ElementsAs(ctx, &planRoles, false)

		planRoleIds := make(map[string]bool)
		for _, role := range planRoles {
			planRoleIds[role.Id.ValueString()] = true
		}

		var removedRoleIds []string
		for _, stateRole := range stateRoles {
			if !planRoleIds[stateRole.Id.ValueString()] {
				removedRoleIds = append(removedRoleIds, stateRole.Id.ValueString())
			}
		}

		// If roles are being removed, first disable them
		if len(removedRoleIds) > 0 {
			disableRoles, err := ConstructAppRoleIsEnabledToFalse(ctx, state.AppRoles, removedRoleIds)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing app roles for disable",
					fmt.Sprintf("Could not construct app roles with disabled state: %s", err.Error()),
				)
				return
			}

			disablePatch := graphmodels.NewApplication()
			disablePatch.SetAppRoles(disableRoles)

			_, err = r.client.
				Applications().
				ByApplicationId(state.Id.ValueString()).
				Patch(ctx, disablePatch, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
				return
			}

			tflog.Debug(ctx, fmt.Sprintf("Disabled %d app roles for %s with ID: %s, waiting before removal", len(removedRoleIds), ResourceName, state.Id.ValueString()))
			time.Sleep(5 * time.Second)
		}

		// Step 2: Apply updated app roles configuration
		appRoles, err := ConstructAppRolesForUpdate(ctx, plan.AppRoles)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing app roles",
				fmt.Sprintf("Could not construct app roles: %s", err.Error()),
			)
			return
		}

		appRolePatch := graphmodels.NewApplication()
		appRolePatch.SetAppRoles(appRoles)

		_, err = r.client.
			Applications().
			ByApplicationId(state.Id.ValueString()).
			Patch(ctx, appRolePatch, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully updated app roles for %s with ID: %s", ResourceName, state.Id.ValueString()))
	}

	// Key credentials are managed by the separate application_certificate_credential resource

	ownersToAdd, ownersToRemove := ResolveOwnerChanges(ctx, &state, &plan)

	adapter := r.client.GetAdapter()

	for _, ownerId := range ownersToRemove {
		if err := RemoveOwner(ctx, adapter, state.Id.ValueString(), ownerId); err != nil {
			resp.Diagnostics.AddError(
				"Error removing owner",
				fmt.Sprintf("Could not remove owner %s from %s: %s", ownerId, ResourceName, err.Error()),
			)
			return
		}
	}

	for _, ownerId := range ownersToAdd {
		if err := AddOwner(ctx, adapter, state.Id.ValueString(), ownerId); err != nil {
			resp.Diagnostics.AddError(
				"Error adding owner",
				fmt.Sprintf("Could not add owner %s to %s: %s", ownerId, ResourceName, err.Error()),
			)
			return
		}
	}

	tflog.Debug(ctx, "Waiting 15 seconds for eventual consistency")
	time.Sleep(15 * time.Second)

	plan.Id = state.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with Id: %s", ResourceName, state.Id.ValueString()))
}

// Delete handles the Delete operation for application resources.
//
// Operation: Deletes application with optional hard delete
// API Calls:
//   - DELETE /applications/{id}
//   - GET /directory/deletedItems/microsoft.graph.application (if hard_delete is true)
//   - DELETE /directory/deletedItems/{id} (if hard_delete is true)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-delete?view=graph-rest-beta
// Note: When hard_delete is true, performs soft delete then permanent deletion from directory deleted items after eventual consistency delay
func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	applicationId := data.Id.ValueString()

	// Define the soft delete function for applications
	softDeleteFunc := func(ctx context.Context) error {
		return r.client.
			Applications().
			ByApplicationId(applicationId).
			Delete(ctx, nil)
	}

	deleteOpts := directory.DeleteOptions{
		MaxRetries:    10,
		RetryInterval: 5 * time.Second,
		ResourceType:  directory.ResourceTypeApplication,
		ResourceID:    applicationId,
		ResourceName:  data.DisplayName.ValueString(),
	}

	err := directory.ExecuteDeleteWithVerification(
		ctx,
		r.client,
		softDeleteFunc,
		data.HardDelete.ValueBool(),
		deleteOpts,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Failed",
			fmt.Sprintf("Failed to delete %s: %s", ResourceName, err.Error()),
		)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
