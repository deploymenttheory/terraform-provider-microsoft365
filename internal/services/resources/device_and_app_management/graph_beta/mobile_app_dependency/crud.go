package graphBetaMobileAppDependency

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for Mobile App Dependency resources.
//
// Operation: Creates a mobile app dependency while preserving other app relationships
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/updateRelationships
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapp-updaterelationships?view=graph-rest-beta
// Note: updateRelationships returns no content; the relationship ID is discovered during read-back.
func (r *MobileAppDependencyResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var object MobileAppDependencyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	sourceAppID := object.SourceID.ValueString()

	// Serialize the collection update and read-back with other dependencies for this app.
	unlock, err := lockAppRelationships(ctx, sourceAppID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to lock app relationships", err.Error())
		return
	}
	defer unlock()

	_, err = r.changeDependency(ctx, &object, constants.TfOperationCreate)
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = func(ctx context.Context, state tfsdk.State) bool {
		var refreshedState MobileAppDependencyResourceModel
		if state.Get(ctx, &refreshedState).HasError() {
			return false
		}

		return !refreshedState.ID.IsNull() && !refreshedState.ID.IsUnknown() &&
			refreshedState.ID.ValueString() != "" &&
			refreshedState.SourceID.Equal(object.SourceID) &&
			refreshedState.TargetID.Equal(object.TargetID) &&
			refreshedState.DependencyType.Equal(object.DependencyType)
	}

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

// Read handles the Read operation for Mobile App Dependency resources.
//
// Operation: Retrieves a mobile app dependency by ID or discovers its ID after creation
// API Calls:
//   - GET /deviceAppManagement/mobileAppRelationships/{relationshipId}
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileappdependency-get?view=graph-rest-beta
// Note: The global relationship endpoint supports imports using only the relationship ID.
func (r *MobileAppDependencyResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var object MobileAppDependencyResourceModel
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

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	sourceAppID := object.SourceID.ValueString()
	relationshipID := object.ID.ValueString()

	var respBaseResource graphmodels.MobileAppRelationshipable

	if !object.ID.IsNull() && !object.ID.IsUnknown() && relationshipID != "" {
		// Read by relationship ID so imports do not require source_id in existing state.
		var err error
		respBaseResource, err = r.client.
			DeviceAppManagement().
			MobileAppRelationships().
			ByMobileAppRelationshipId(relationshipID).
			Get(ctx, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
			return
		}
	} else {
		// Discover the relationship ID because updateRelationships returns no content.
		relationships, err := r.listRelationships(ctx, sourceAppID)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
			return
		}

		for _, relationship := range relationships {
			if !isIncomingRelationship(relationship) && relationship.GetTargetId() != nil &&
				strings.EqualFold(*relationship.GetTargetId(), object.TargetID.ValueString()) {
				respBaseResource = relationship
				break
			}
		}

		if respBaseResource == nil {
			resp.Diagnostics.AddError(
				"App dependency not found (404)",
				"The dependency is not yet visible after the relationship update.",
			)
			return
		}
	}

	// Convert the base Graph relationship to the dependency type before mapping state.
	dependency, ok := respBaseResource.(graphmodels.MobileAppDependencyable)
	if !ok {
		resp.Diagnostics.AddError(
			"Resource type mismatch",
			fmt.Sprintf(
				"Expected resource of type MobileAppDependencyable but got %T",
				respBaseResource,
			),
		)
		return
	}

	diags := mapResourceToState(ctx, &object, dependency)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set identity after read-back has resolved the relationship ID.
	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Mobile App Dependency resources.
//
// Operation: Updates the dependency behavior while preserving other app relationships
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/updateRelationships
//   - GET /deviceAppManagement/mobileAppRelationships/{relationshipId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapp-updaterelationships?view=graph-rest-beta
func (r *MobileAppDependencyResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan MobileAppDependencyResourceModel
	var state MobileAppDependencyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	plan.ID = state.ID

	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	sourceAppID := state.SourceID.ValueString()

	unlock, err := lockAppRelationships(ctx, sourceAppID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to lock app relationships", err.Error())
		return
	}
	defer unlock()

	// Only dependency_type can be updated; source_id and target_id require replacement.
	if !plan.DependencyType.Equal(state.DependencyType) {
		_, err = r.changeDependency(ctx, &plan, constants.TfOperationUpdate)
		if err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationUpdate,
				r.WritePermissions,
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = func(ctx context.Context, state tfsdk.State) bool {
		var object MobileAppDependencyResourceModel
		if state.Get(ctx, &object).HasError() {
			return false
		}

		return !object.ID.IsNull() && !object.ID.IsUnknown() && object.ID.ValueString() != "" &&
			object.SourceID.Equal(plan.SourceID) &&
			object.TargetID.Equal(plan.TargetID) &&
			object.DependencyType.Equal(plan.DependencyType)
	}

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

// Delete handles the Delete operation for Mobile App Dependency resources.
//
// Operation: Removes the managed dependency while preserving other app relationships
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/updateRelationships
//   - GET /deviceAppManagement/mobileAppRelationships/{relationshipId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapp-updaterelationships?view=graph-rest-beta
// Note: An empty outgoing relationship collection is sent when deleting the last dependency.
func (r *MobileAppDependencyResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var object MobileAppDependencyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting to delete resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Delete,
		DeleteTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	sourceAppID := object.SourceID.ValueString()
	relationshipID := object.ID.ValueString()

	tflog.Debug(
		ctx,
		fmt.Sprintf(
			"Deleting Mobile App Dependency with ID: %s for source app ID: %s",
			relationshipID,
			sourceAppID,
		),
	)

	unlock, err := lockAppRelationships(ctx, sourceAppID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to lock app relationships", err.Error())
		return
	}
	defer unlock()

	changed, err := r.changeDependency(ctx, &object, constants.TfOperationDelete)
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationDelete,
			r.WritePermissions,
		)
		return
	}

	if changed {
		// Wait until Graph confirms deletion before removing the dependency from state.
		err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
			_, err := r.client.
				DeviceAppManagement().
				MobileAppRelationships().
				ByMobileAppRelationshipId(relationshipID).
				Get(ctx, nil)

			if err == nil {
				return false, errMissingDependencyDeletion
			}

			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 {
				return true, nil
			}
			if errors.IsRetryableReadError(&errorInfo) {
				return false, fmt.Errorf("verify dependency deletion: %w", err)
			}

			return false, &crud.FatalPollError{
				Err: fmt.Errorf("verify dependency deletion: %w", err),
			}
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error verifying resource deletion",
				fmt.Sprintf(
					"Could not verify resource deletion: %s: %s",
					ResourceName,
					err.Error(),
				),
			)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
