package graphBetaDeviceAndAppManagementAppAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Mobile App Assignment resources.
//
//   - Retrieves the planned configuration from the create request
//   - Constructs the resource request body from the plan
//   - Sends POST request to create the base resource and settings
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API with retry
//   - Updates the final state with the fresh data from the API
func (r *MobileAppAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object MobileAppAssignmentResourceModel

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

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	requestBody, err := ConstructMobileAppAssignment(ctx, object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Create the assignment for the mobile app
	createdResource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.MobileAppId.ValueString()).
		Assignments().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Create Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource creation",
			fmt.Sprintf("Failed to verify resource creation: %s", err),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for Mobile App Assignment resources.
//
//   - Retrieves the current state from the read request
//   - Gets the base resource details from the API
//   - Maps the base resource details to Terraform state
func (r *MobileAppAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object MobileAppAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s for Mobile App: %s",
		ResourceName, object.ID.ValueString(), object.MobileAppId.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// If we don't have an ID yet (during initial creation), try to find the assignment by target and intent
	if object.ID.IsNull() {
		assignments, err := r.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.MobileAppId.ValueString()).
			Assignments().
			Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
			return
		}

		// Find the assignment that matches our configuration
		for _, assign := range assignments.GetValue() {
			// Compare relevant fields to identify our assignment
			// This is a simplistic approach - you may need more sophisticated matching
			if matchesAssignment(ctx, object, assign) {
				object.ID = types.StringValue(*assign.GetId())
				break
			}
		}

		if object.ID.IsNull() {
			// Couldn't find a matching assignment
			resp.Diagnostics.AddWarning(
				"Resource not found",
				fmt.Sprintf("Could not find assignment matching specified criteria for Mobile App: %s",
					object.MobileAppId.ValueString()),
			)
			resp.State.RemoveResource(ctx)
			return
		}
	}

	// Get the specific assignment by ID
	resource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.MobileAppId.ValueString()).
		Assignments().
		ByMobileAppAssignmentId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, object, resource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Mobile App Assignment resources.
func (r *MobileAppAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object MobileAppAssignmentResourceModel

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

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	requestBody, err := ConstructMobileAppAssignment(ctx, object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Perform the PATCH operation to update the assignment
	_, err = r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.MobileAppId.ValueString()).
		Assignments().
		ByMobileAppAssignmentId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Update Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource update",
			fmt.Sprintf("Failed to verify resource update: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Mobile App Assignment resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
func (r *MobileAppAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object MobileAppAssignmentResourceModel

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
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.MobileAppId.ValueString()).
		Assignments().
		ByMobileAppAssignmentId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// Helper function to determine if an assignment matches our resource
func matchesAssignment(ctx context.Context, object MobileAppAssignmentResourceModel, assign interface{}) bool {
	assignment, ok := assign.(graphmodels.MobileAppAssignmentable)
	if !ok {
		return false
	}

	// Check if intent matches
	if assignment.GetIntent() != nil {
		intentStr := assignment.GetIntent().String()
		if !object.Intent.IsNull() && object.Intent.ValueString() != intentStr {
			return false
		}
	}

	// Check if target matches (simplified - you may need more complex logic)
	target := assignment.GetTarget()
	if target != nil {
		// Check target type
		odataType := target.GetOdataType()
		if odataType != nil {
			targetType := getTargetTypeFromOdataType(*odataType)
			if !object.Target.TargetType.IsNull() && object.Target.TargetType.ValueString() != targetType {
				return false
			}
		}

		// Check group ID if applicable
		if !object.Target.GroupId.IsNull() {
			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				if groupTarget.GetGroupId() == nil || object.Target.GroupId.ValueString() != *groupTarget.GetGroupId() {
					return false
				}
			}
		}

		// Check filter ID if applicable
		if !object.Target.DeviceAndAppManagementAssignmentFilterId.IsNull() {
			filterId := target.GetDeviceAndAppManagementAssignmentFilterId()
			if filterId == nil || object.Target.DeviceAndAppManagementAssignmentFilterId.ValueString() != *filterId {
				return false
			}
		}
	}

	return true
}

// Helper to get target type from odata.type
func getTargetTypeFromOdataType(odataType string) string {
	switch odataType {
	case "#microsoft.graph.allDevicesAssignmentTarget":
		return "allDevices"
	case "#microsoft.graph.allLicensedUsersAssignmentTarget":
		return "allLicensedUsers"
	case "#microsoft.graph.groupAssignmentTarget":
		return "groupAssignment"
	case "#microsoft.graph.exclusionGroupAssignmentTarget":
		return "exclusionGroupAssignment"
	case "#microsoft.graph.androidFotaDeploymentAssignmentTarget":
		return "androidFotaDeployment"
	case "#microsoft.graph.configurationManagerCollectionAssignmentTarget":
		return "configurationManagerCollection"
	default:
		return ""
	}
}
