package graphBetaRoleScopeTags

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/retry"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Role Scope Tags in Microsoft Intune.
//
// The function:
// - Gets the planned Role Scope Tag configuration from Terraform
// - Constructs the role scope tag body with display name and description
// - Sends a POST request to create the role scope tag in Intune
// - Maps the response to the Terraform state
// - Performs a Read operation to ensure state consistency
func (r *RoleScopeTagsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object RoleScopeTagsProfileResourceModel

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

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	err = retry.RetryableIntuneOperation(ctx, "create resource", retry.IntuneWrite, func() error {
		var reqErr error
		requestBody, reqErr = r.client.
			DeviceManagement().
			RoleScopeTags().
			Post(ctx, requestBody, nil)
		return reqErr
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

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

// Read handles the Read operation for Role Scope Tags in Microsoft Intune.
//
// The function:
// - Gets the role scope tag ID from the current Terraform state
// - Retrieves the role scope tag details from Intune using the ID
// - Maps the role scope tag properties (display name, description, isBuiltIn)
// - Maps any auto-assignments if they exist
// - Updates the Terraform state with the current Intune configuration
func (r *RoleScopeTagsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object RoleScopeTagsProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var baseResource graphmodels.RoleScopeTagable
	err := retry.RetryableIntuneOperation(ctx, "read base resource", retry.IntuneRead, func() error {
		var err error
		baseResource, err = r.client.
			DeviceManagement().
			RoleScopeTags().
			ByRoleScopeTagId(object.ID.ValueString()).
			Get(ctx, nil)
		return err
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, baseResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for Role Scope Tags in Microsoft Intune.
//
// The function:
// - Gets the planned changes for the role scope tag from Terraform
// - Constructs the update request with modified display name or description
// - Sends a PATCH request to update the role scope tag in Intune
// - Performs a Read operation to refresh the Terraform state
// Note: Built-in role scope tags cannot be modified
func (r *RoleScopeTagsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object RoleScopeTagsProfileResourceModel

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

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	err = retry.RetryableAssignmentOperation(ctx, "update assignment", func() error {
		_, err := r.client.
			DeviceManagement().
			RoleScopeTags().
			ByRoleScopeTagId(object.ID.ValueString()).
			Patch(ctx, requestBody, nil)
		return err
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for Role Scope Tags in Microsoft Intune.
//
// The function:
// - Gets the role scope tag ID from the Terraform state
// - Sends a DELETE request to remove the role scope tag from Intune
// - Removes the resource from the Terraform state
// Note: Built-in role scope tags cannot be deleted
func (r *RoleScopeTagsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object RoleScopeTagsProfileResourceModel

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

	err := retry.RetryableIntuneOperation(ctx, "delete resource", retry.IntuneWrite, func() error {
		return r.client.
			DeviceManagement().
			RoleScopeTags().
			ByRoleScopeTagId(object.ID.ValueString()).
			Delete(ctx, nil)
	})

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.ReadPermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
