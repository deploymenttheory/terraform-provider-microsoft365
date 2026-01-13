package graphBetaAppleUserInitiatedEnrollmentProfileAssignment

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

// Create handles the Create operation for Apple User Initiated Enrollment Profile Assignment resources.
func (r *AppleUserInitiatedEnrollmentProfileAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AppleUserInitiatedEnrollmentProfileAssignmentResourceModel

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

	requestBody, err := ConstructAppleUserInitiatedEnrollmentProfileAssignment(ctx, r.client, object, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		DeviceManagement().
		AppleUserInitiatedEnrollmentProfiles().
		ByAppleUserInitiatedEnrollmentProfileId(object.AppleUserInitiatedEnrollmentProfileId.ValueString()).
		Assignments().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

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

// Read handles the Read operation for Apple User Initiated Enrollment Profile Assignment resources.
func (r *AppleUserInitiatedEnrollmentProfileAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AppleUserInitiatedEnrollmentProfileAssignmentResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s for Apple User Initiated Enrollment Profile: %s",
		ResourceName, object.ID.ValueString(), object.AppleUserInitiatedEnrollmentProfileId.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Get all assignments and find the one with matching ID or criteria
	assignments, err := r.client.
		DeviceManagement().
		AppleUserInitiatedEnrollmentProfiles().
		ByAppleUserInitiatedEnrollmentProfileId(object.AppleUserInitiatedEnrollmentProfileId.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	// Find the assignment with matching ID or criteria
	var foundAssignment graphmodels.AppleEnrollmentProfileAssignmentable
	if assignments != nil && assignments.GetValue() != nil {
		for _, assignment := range assignments.GetValue() {
			// If we have an ID, match by ID, otherwise match by criteria
			if !object.ID.IsNull() && !object.ID.IsUnknown() {
				if assignment.GetId() != nil && *assignment.GetId() == object.ID.ValueString() {
					foundAssignment = assignment
					break
				}
			} else if matchesAssignment(ctx, object, assignment) {
				foundAssignment = assignment
				object.ID = types.StringValue(*assignment.GetId())
				break
			}
		}
	}

	if foundAssignment == nil {
		tflog.Debug(ctx, fmt.Sprintf("Assignment with ID %s not found in collection", object.ID.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}

	object = MapRemoteStateToTerraform(ctx, object, foundAssignment)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Apple User Initiated Enrollment Profile Assignment resources.
func (r *AppleUserInitiatedEnrollmentProfileAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AppleUserInitiatedEnrollmentProfileAssignmentResourceModel
	var state AppleUserInitiatedEnrollmentProfileAssignmentResourceModel

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

	requestBody, err := ConstructAppleUserInitiatedEnrollmentProfileAssignment(ctx, r.client, plan, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Update the assignment
	updatedResource, err := r.client.
		DeviceManagement().
		AppleUserInitiatedEnrollmentProfiles().
		ByAppleUserInitiatedEnrollmentProfileId(plan.AppleUserInitiatedEnrollmentProfileId.ValueString()).
		Assignments().
		ByAppleEnrollmentProfileAssignmentId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(*updatedResource.GetId())

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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Apple User Initiated Enrollment Profile Assignment resources.
func (r *AppleUserInitiatedEnrollmentProfileAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AppleUserInitiatedEnrollmentProfileAssignmentResourceModel

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
		AppleUserInitiatedEnrollmentProfiles().
		ByAppleUserInitiatedEnrollmentProfileId(object.AppleUserInitiatedEnrollmentProfileId.ValueString()).
		Assignments().
		ByAppleEnrollmentProfileAssignmentId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfTfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// Helper function to determine if an assignment matches our resource
func matchesAssignment(ctx context.Context, object AppleUserInitiatedEnrollmentProfileAssignmentResourceModel, assign graphmodels.AppleEnrollmentProfileAssignmentable) bool {
	tflog.Debug(ctx, fmt.Sprintf("Matching assignment for resource: %s", ResourceName))

	if assign == nil {
		return false
	}

	// Compare target details based on the API response structure
	if target := assign.GetTarget(); target != nil {
		targetType := getTargetTypeFromTarget(target)
		if targetType != object.Target.TargetType.ValueString() {
			return false
		}

		// For group assignments, compare identifiers
		if targetType == "group" || targetType == "exclusionGroup" {
			// Extract group ID or Entra Object ID from target
			if !object.Target.GroupId.IsNull() {
				if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
					if groupTarget.GetGroupId() != nil && *groupTarget.GetGroupId() != object.Target.GroupId.ValueString() {
						return false
					}
				}
			}
		}

		// For user assignments, compare Entra Object ID
		if targetType == "user" {
			if !object.Target.EntraObjectId.IsNull() {
				// For user assignments, check if the group ID matches the user's Entra Object ID
				if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
					if groupTarget.GetGroupId() != nil && *groupTarget.GetGroupId() != object.Target.EntraObjectId.ValueString() {
						return false
					}
				}
			}
		}
	}

	return true
}

// Helper function to extract target type from the target object
func getTargetTypeFromTarget(target graphmodels.DeviceAndAppManagementAssignmentTargetable) string {
	if target == nil {
		return ""
	}

	odataType := target.GetOdataType()
	if odataType == nil {
		return ""
	}

	switch *odataType {
	case "#microsoft.graph.allLicensedUsersAssignmentTarget":
		return "allUsers"
	case "#microsoft.graph.groupAssignmentTarget":
		return "group" // Note: This could also be "user" - context needed to determine
	case "#microsoft.graph.exclusionGroupAssignmentTarget":
		return "exclusionGroup"
	default:
		return ""
	}
}
