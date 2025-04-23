package graphBetaRoleDefinition

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object RoleDefinitionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	defer cancel()

	// Intune roles require unique display_names
	isBuiltIn := object.IsBuiltInRoleDefinition.ValueBool() || object.IsBuiltIn.ValueBool()
	if !isBuiltIn && !object.DisplayName.IsNull() && !object.DisplayName.IsUnknown() {
		if err := checkRoleNameUniqueness(ctx, r.client, object.DisplayName.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Role Name Not Unique",
				err.Error(),
			)
			return
		}
	}

	requestBody, err := constructResource(ctx, r.client, &object, resp, r.ReadPermissions, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		DeviceManagement().
		RoleDefinitions().
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

	readResp := &resource.ReadResponse{
		State: resp.State,
	}
	r.Read(ctx, resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}, readResp)

	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State = readResp.State

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object RoleDefinitionResourceModel

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

	// 1️⃣ Fetch base resource
	resource, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(object.ID.ValueString()).
		Get(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}
	MapRemoteResourceStateToTerraform(ctx, &object, resource)

	// 2️⃣ List the assignments
	assignmentsList, err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(object.ID.ValueString()).
		RoleAssignments().
		Get(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read Assignments List", r.ReadPermissions)
		return
	}

	// 3️⃣ Pull each assignment’s full details
	detailedResponse := graphmodels.NewRoleAssignmentCollectionResponse()
	var detailedAssignments []graphmodels.RoleAssignmentable

	if assignmentsList != nil && assignmentsList.GetValue() != nil {
		for _, listAssignment := range assignmentsList.GetValue() {
			if listAssignment == nil || listAssignment.GetId() == nil {
				continue
			}
			assignmentID := *listAssignment.GetId()
			full, err := r.client.
				DeviceManagement().
				RoleAssignments().
				ByDeviceAndAppManagementRoleAssignmentId(assignmentID).
				Get(ctx, nil)
			if err != nil {
				tflog.Warn(ctx, fmt.Sprintf("Failed to fetch details for assignment ID %s: %s", assignmentID, err))
				continue
			}
			// append to the RoleAssignmentable slice
			detailedAssignments = append(detailedAssignments, full)
		}
	}
	detailedResponse.SetValue(detailedAssignments)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read method for: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for role definitions and assignments,
// tracking assignments strictly by ID
func (r *RoleDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planObj, stateObj RoleDefinitionResourceModel

	tflog.Info(ctx, "Starting Update for RoleDefinition", map[string]interface{}{"resource_type": r.TypeName})

	// Load plan & state
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planObj)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &stateObj)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, planObj.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	defer cancel()

	// 1️⃣ PATCH the base RoleDefinition
	builder := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(planObj.ID.ValueString())
	requestBody, err := constructResource(ctx, r.client, &planObj, resp, r.ReadPermissions, true)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing resource", err.Error())
		return
	}
	if _, err := builder.Patch(ctx, requestBody, nil); err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update RoleDefinition", r.WritePermissions)
		return
	}
	tflog.Debug(ctx, "Patched base RoleDefinition successfully")

	// 9️⃣ Let Read function handle state
	tflog.Debug(ctx, "Using Read to refresh final state")
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}, readResp)
	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State = readResp.State

	tflog.Info(ctx, "Update completed successfully")
}

// Delete handles the Delete operation for the RoleDefinition resource.
func (r *RoleDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RoleDefinitionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(data.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Custom role definition deleted successfully")

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
