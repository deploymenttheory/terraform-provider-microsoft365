package graphBetaMobileAppSupersedence

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	apprelationships "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management/app_relationships"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Mobile App Supersedence resources.
//
// Operation: Creates a mobile app supersedence while preserving other app relationships
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/updateRelationships
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapp-updaterelationships?view=graph-rest-beta
// Note: updateRelationships returns no content; the relationship ID is discovered during read-back.
func (r *MobileAppSupersedenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object MobileAppSupersedenceResourceModel

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

	sourceAppID := object.SourceID.ValueString()

	// Serialize the collection update and read-back with other supersedence relationships for this app.
	unlock, err := lockAppRelationships(ctx, sourceAppID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to lock app relationships", err.Error())
		return
	}
	defer unlock()

	_, err = r.changeSupersedence(ctx, &object, constants.TfOperationCreate)
	if err != nil {
		if stderrors.Is(err, errExistingRelationship) {
			resp.Diagnostics.AddError("App relationship already exists", err.Error())
			return
		}
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
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
		var refreshedState MobileAppSupersedenceResourceModel
		if state.Get(ctx, &refreshedState).HasError() {
			return false
		}

		return !refreshedState.ID.IsNull() && !refreshedState.ID.IsUnknown() &&
			refreshedState.ID.ValueString() != "" &&
			refreshedState.SourceID.Equal(object.SourceID) &&
			refreshedState.TargetID.Equal(object.TargetID) &&
			refreshedState.SupersedenceType.Equal(object.SupersedenceType)
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

// Read handles the Read operation for Mobile App Supersedence resources.
//
// Operation: Retrieves a mobile app supersedence by ID or discovers its ID after creation
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileappsupersedence-get?view=graph-rest-beta
// Note: Composite relationship IDs resolve the source app for ID-only imports.
func (r *MobileAppSupersedenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object MobileAppSupersedenceResourceModel
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

	sourceAppID := object.SourceID.ValueString()
	relationshipID := object.ID.ValueString()

	var respBaseResource graphmodels.MobileAppRelationshipable
	var err error
	if sourceAppID == "" {
		sourceAppID = apprelationships.SourceID(relationshipID)
	}
	if sourceAppID != "" {
		object.SourceID = types.StringValue(sourceAppID)
		var relationships []graphmodels.MobileAppRelationshipable
		relationships, err = r.listRelationships(ctx, sourceAppID)
		if err == nil {
			respBaseResource, err = apprelationships.Find(relationships, relationshipID, object.TargetID.ValueString())
		}
	} else {
		// Retain support for legacy opaque IDs whose source cannot be derived.
		respBaseResource, err = r.client.
			DeviceAppManagement().
			MobileAppRelationships().
			ByMobileAppRelationshipId(relationshipID).
			Get(ctx, nil)
	}

	if stderrors.Is(err, apprelationships.ErrNotFound) {
		if operation == constants.TfOperationRead {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				"App supersedence not found (404)",
				"The supersedence is not yet visible after the relationship update.",
			)
		}
		return
	}
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions, errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}

	// Convert the base Graph relationship to the supersedence type before mapping state.
	supersedence, ok := respBaseResource.(graphmodels.MobileAppSupersedenceable)
	if !ok {
		resp.Diagnostics.AddError(
			"Resource type mismatch",
			fmt.Sprintf("Expected resource of type MobileAppSupersedenceable but got %T", respBaseResource),
		)
		return
	}

	diags := mapResourceToState(ctx, &object, supersedence)
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

// Update handles the Update operation for Mobile App Supersedence resources.
//
// Operation: Updates the supersedence behavior while preserving other app relationships
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/updateRelationships
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapp-updaterelationships?view=graph-rest-beta
func (r *MobileAppSupersedenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MobileAppSupersedenceResourceModel
	var state MobileAppSupersedenceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	plan.ID = state.ID

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
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

	// Only supersedence_type can be updated; source_id and target_id require replacement.
	if !plan.SupersedenceType.Equal(state.SupersedenceType) {
		_, err = r.changeSupersedence(ctx, &plan, constants.TfOperationUpdate)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
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
		var object MobileAppSupersedenceResourceModel
		if state.Get(ctx, &object).HasError() {
			return false
		}

		return !object.ID.IsNull() && !object.ID.IsUnknown() && object.ID.ValueString() != "" &&
			object.SourceID.Equal(plan.SourceID) &&
			object.TargetID.Equal(plan.TargetID) &&
			object.SupersedenceType.Equal(plan.SupersedenceType)
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

// Delete handles the Delete operation for Mobile App Supersedence resources.
//
// Operation: Removes the managed supersedence while preserving other app relationships
// API Calls:
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//   - POST /deviceAppManagement/mobileApps/{mobileAppId}/updateRelationships
//   - GET /deviceAppManagement/mobileApps/{mobileAppId}/relationships
//
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapp-updaterelationships?view=graph-rest-beta
// Note: An empty outgoing relationship collection is sent when deleting the last supersedence.
func (r *MobileAppSupersedenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object MobileAppSupersedenceResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting to delete resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	sourceAppID := object.SourceID.ValueString()
	relationshipID := object.ID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Deleting Mobile App Supersedence with ID: %s for source app ID: %s", relationshipID, sourceAppID))

	unlock, err := lockAppRelationships(ctx, sourceAppID)
	if err != nil {
		resp.Diagnostics.AddError("Unable to lock app relationships", err.Error())
		return
	}
	defer unlock()

	changed, err := r.changeSupersedence(ctx, &object, constants.TfOperationDelete)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	if changed {
		// Wait until Graph confirms deletion before removing the supersedence from state.
		err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
			relationships, err := r.listRelationships(ctx, sourceAppID)
			if err == nil {
				_, findErr := apprelationships.Find(relationships, relationshipID, object.TargetID.ValueString())
				if stderrors.Is(findErr, apprelationships.ErrNotFound) {
					return true, nil
				}
				return false, errMissingSupersedenceDeletion
			}

			errorInfo := errors.GraphError(ctx, err)
			if errorInfo.StatusCode == 404 {
				return true, nil
			}
			if errors.IsRetryableReadError(&errorInfo) {
				return false, fmt.Errorf("verify supersedence deletion: %w", err)
			}

			return false, &crud.FatalPollError{
				Err: fmt.Errorf("verify supersedence deletion: %w", err),
			}
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error verifying resource deletion",
				fmt.Sprintf("Could not verify resource deletion: %s: %s", ResourceName, err.Error()),
			)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
