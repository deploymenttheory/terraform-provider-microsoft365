package graphBetaWinGetApp

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *WinGetAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WinGetAppResourceModel

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

	resource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		Post(context.Background(), requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*resource.GetId())

	resourceAsWinGetApp, ok := resource.(models.WinGetAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Error",
			fmt.Sprintf("Expected resource of type WinGetApp for %s_%s, but got %T",
				r.ProviderTypeName, r.TypeName, resource),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, resourceAsWinGetApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *WinGetAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WinGetAppResourceModel
	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resource, err := r.client.DeviceAppManagement().MobileApps().
		ByMobileAppId(state.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Resource type: %T", resource))

	mobileApp, ok := resource.(models.MobileAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Error",
			fmt.Sprintf("Expected resource of type MobileApp for %s_%s, but got %T",
				r.ProviderTypeName, r.TypeName, resource),
		)
		return
	}

	resourceAsWinGetApp, ok := mobileApp.(models.WinGetAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Error",
			fmt.Sprintf("Expected resource of type WinGetApp for %s_%s, but got %T",
				r.ProviderTypeName, r.TypeName, mobileApp),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, resourceAsWinGetApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *WinGetAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object WinGetAppResourceModel

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
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	resource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Resource type after update: %T", resource))

	mobileApp, ok := resource.(models.MobileAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Error",
			fmt.Sprintf("Expected resource of type MobileApp for %s_%s, but got %T",
				r.ProviderTypeName, r.TypeName, resource),
		)
		return
	}

	resourceAsWinGetApp, ok := mobileApp.(models.WinGetAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Error",
			fmt.Sprintf("Expected resource of type WinGetApp for %s_%s, but got %T",
				r.ProviderTypeName, r.TypeName, mobileApp),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, resourceAsWinGetApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *WinGetAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WinGetAppResourceModel

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

	err := r.client.DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
