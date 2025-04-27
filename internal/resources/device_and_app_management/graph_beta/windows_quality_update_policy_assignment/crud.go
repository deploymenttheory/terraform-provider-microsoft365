package graphBetaWindowsQualityUpdatePolicyAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Windows Driver Update Profile Assignments.
func (r *WindowsQualityUpdateProfileAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsQualityUpdateProfileAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if object.WindowsQualityUpdatePolicyID.IsNull() || object.WindowsQualityUpdatePolicyID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to create assignments.",
		)
		return
	}
	if len(object.Assignments) == 0 {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"At least one assignment block is required.",
		)
		return
	}

	assignRequest, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assign request",
			fmt.Sprintf("Could not construct assign request: %s", err.Error()),
		)
		return
	}

	profileID := object.WindowsQualityUpdatePolicyID.ValueString()

	err = r.client.
		DeviceManagement().
		WindowsQualityUpdateProfiles().
		ByWindowsQualityUpdateProfileId(profileID).
		Assign().
		Post(ctx, assignRequest, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResp := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}, &readResp)

	if readResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResp.Diagnostics...)
		return
	}

	resp.State = readResp.State

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation for Windows Driver Update Profile Assignment resources.
func (r *WindowsQualityUpdateProfileAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsQualityUpdateProfileAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if object.WindowsQualityUpdatePolicyID.IsNull() || object.WindowsQualityUpdatePolicyID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to read assignments.",
		)
		return
	}

	profileID := object.WindowsQualityUpdatePolicyID.ValueString()

	// For this resource, we need to read all assignments and filter them
	// since we're managing multiple assignments in a single resource
	assignmentsResponse, err := r.client.
		DeviceManagement().
		WindowsQualityUpdateProfiles().
		ByWindowsQualityUpdateProfileId(profileID).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	if assignmentsResponse == nil || assignmentsResponse.GetValue() == nil || len(assignmentsResponse.GetValue()) == 0 {
		tflog.Debug(ctx, "No assignments found for profile", map[string]interface{}{
			"profileID": profileID,
		})
		resp.State.RemoveResource(ctx)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, assignmentsResponse.GetValue())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for Windows Driver Update Profile Assignment resources.
func (r *WindowsQualityUpdateProfileAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object WindowsQualityUpdateProfileAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if object.WindowsQualityUpdatePolicyID.IsNull() || object.WindowsQualityUpdatePolicyID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to update an assignment.",
		)
		return
	}

	if len(object.Assignments) == 0 {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"At least one assignment block is required.",
		)
		return
	}

	assignRequest, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assign request",
			fmt.Sprintf("Could not construct assign request: %s", err.Error()),
		)
		return
	}

	profileID := object.WindowsQualityUpdatePolicyID.ValueString()

	err = r.client.
		DeviceManagement().
		WindowsQualityUpdateProfiles().
		ByWindowsQualityUpdateProfileId(profileID).
		Assign().
		Post(ctx, assignRequest, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResp := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}, &readResp)

	if readResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResp.Diagnostics...)
		return
	}

	resp.State = readResp.State

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for Windows Driver Update Profile Assignment resources.
func (r *WindowsQualityUpdateProfileAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsQualityUpdateProfileAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if object.WindowsQualityUpdatePolicyID.IsNull() || object.WindowsQualityUpdatePolicyID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to delete assignments.",
		)
		return
	}

	profileID := object.WindowsQualityUpdatePolicyID.ValueString()

	// For a delete operation, we submit an empty assignments list to effectively
	// remove all assignments managed by this resource
	assignRequest := devicemanagement.NewWindowsQualityUpdateProfilesItemAssignPostRequestBody()
	assignRequest.SetAssignments([]graphmodels.WindowsQualityUpdateProfileAssignmentable{})

	err := r.client.
		DeviceManagement().
		WindowsQualityUpdateProfiles().
		ByWindowsQualityUpdateProfileId(profileID).
		Assign().
		Post(ctx, assignRequest, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
