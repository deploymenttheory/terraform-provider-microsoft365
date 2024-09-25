package graphbetamacospkgapp

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for the MacOSPkgApp resource.
func (r *MacOSPkgAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MacOSPkgAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	app, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	createdApp, err := r.client.DeviceAppManagement().MobileApps().Post(ctx, app, nil)
	if err != nil {
		if crud.PermissionError(err, "Create", r.WritePermissions, resp) {
			return
		} else {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Client error when creating %s_%s", r.ProviderTypeName, r.TypeName),
				err.Error(),
			)
		}
		return
	}
	plan.ID = types.StringValue(*createdApp.GetId())

	// Call Read to fetch the full state and set it in the response
	readResp := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &readResp)

	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation for the MacOSPkgApp resource.
func (r *MacOSPkgAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MacOSPkgAppResourceModel
	tflog.Debug(ctx, "Starting Read method for macOS PKG app")

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading macOS PKG app with ID: %s", state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	app, err := r.client.DeviceAppManagement().MobileApps().ByMobileAppId(state.ID.ValueString()).Get(ctx, nil)

	if err != nil {
		if crud.IsNotFoundError(err) {
			tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found on server, removing from state", r.TypeName, state.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}

		if crud.PermissionError(err, "Read", r.ReadPermissions, resp) {
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName),
			err.Error(),
		)
		return
	}

	// Type assertion to convert MobileAppable to MacOSPkgAppable
	macOSPkgApp, ok := app.(models.MacOSPkgAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Failed",
			fmt.Sprintf("Expected MacOSPkgAppable, got: %T. Please report this issue to the provider developers.", app),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, macOSPkgApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation for the MacOSPkgApp resource.
func (r *MacOSPkgAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state MacOSPkgAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	_, err = r.client.DeviceAppManagement().MobileApps().
		ByMobileAppId(plan.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		if crud.IsNotFoundError(err) {
			tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found on server, removing from state", r.TypeName, plan.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}

		if crud.PermissionError(err, "Update", r.WritePermissions, resp) {
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Client error when updating %s_%s", r.ProviderTypeName, r.TypeName),
			err.Error(),
		)
		return
	}

	// Call Read to fetch the full state and set it in the response
	readResp := resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{State: resp.State}, &readResp)

	resp.Diagnostics.Append(readResp.Diagnostics...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation for the MacOSPkgApp resource.
func (r *MacOSPkgAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MacOSPkgAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.DeviceAppManagement().MobileApps().
		ByMobileAppId(data.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		if crud.IsNotFoundError(err) {
			tflog.Warn(ctx, fmt.Sprintf("%s with ID %s not found on server, removing from state", r.TypeName, data.ID.ValueString()))
			resp.State.RemoveResource(ctx)
			return
		}

		if crud.PermissionError(err, "Delete", r.WritePermissions, resp) {
			return
		}

		resp.Diagnostics.AddError(
			fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName),
			err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
