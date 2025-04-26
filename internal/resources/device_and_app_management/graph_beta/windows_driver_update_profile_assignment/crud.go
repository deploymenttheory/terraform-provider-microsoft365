package graphBetaWindowsDriverUpdateProfileAssignment

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
func (r *WindowsDriverUpdateProfileAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsDriverUpdateProfileAssignmentResourceModel

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

	if object.WindowsDriverUpdateProfileID.IsNull() || object.WindowsDriverUpdateProfileID.ValueString() == "" {
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

	profileID := object.WindowsDriverUpdateProfileID.ValueString()

	err = r.client.
		DeviceManagement().
		WindowsDriverUpdateProfiles().
		ByWindowsDriverUpdateProfileId(profileID).
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
func (r *WindowsDriverUpdateProfileAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsDriverUpdateProfileAssignmentResourceModel

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

	// Validate that Windows Driver Update Profile ID is provided
	if object.WindowsDriverUpdateProfileID.IsNull() || object.WindowsDriverUpdateProfileID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to read assignments.",
		)
		return
	}

	profileID := object.WindowsDriverUpdateProfileID.ValueString()

	// For this resource, we need to read all assignments and filter them
	// since we're managing multiple assignments in a single resource
	assignmentsResponse, err := r.client.
		DeviceManagement().
		WindowsDriverUpdateProfiles().
		ByWindowsDriverUpdateProfileId(profileID).
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
func (r *WindowsDriverUpdateProfileAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object WindowsDriverUpdateProfileAssignmentResourceModel

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

	if object.WindowsDriverUpdateProfileID.IsNull() || object.WindowsDriverUpdateProfileID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to update an assignment.",
		)
		return
	}

	// Validate that at least one assignment block is provided
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

	profileID := object.WindowsDriverUpdateProfileID.ValueString()

	// Update the assignments
	err = r.client.
		DeviceManagement().
		WindowsDriverUpdateProfiles().
		ByWindowsDriverUpdateProfileId(profileID).
		Assign().
		Post(ctx, assignRequest, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Set the state with what we know
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Now call Read to refresh the state with the latest data
	readResp := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}, &readResp)

	if readResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResp.Diagnostics...)
		return
	}

	// Update the state with the data returned from Read
	resp.State = readResp.State

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for Windows Driver Update Profile Assignment resources.
func (r *WindowsDriverUpdateProfileAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsDriverUpdateProfileAssignmentResourceModel

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

	if object.WindowsDriverUpdateProfileID.IsNull() || object.WindowsDriverUpdateProfileID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Required Parameter",
			"The windows_driver_update_profile_id field is required to delete assignments.",
		)
		return
	}

	profileID := object.WindowsDriverUpdateProfileID.ValueString()

	// For a delete operation, we'll submit an empty assignments list to effectively
	// remove all assignments managed by this resource
	assignRequest := devicemanagement.NewWindowsDriverUpdateProfilesItemAssignPostRequestBody()
	assignRequest.SetAssignments([]graphmodels.WindowsDriverUpdateProfileAssignmentable{})

	err := r.client.
		DeviceManagement().
		WindowsDriverUpdateProfiles().
		ByWindowsDriverUpdateProfileId(profileID).
		Assign().
		Post(ctx, assignRequest, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
