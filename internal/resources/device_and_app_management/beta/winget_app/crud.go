package graphBetaWinGetApp

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *WinGetAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WinGetAppResourceModel

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

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	resource, err := r.client.DeviceAppManagement().MobileApps().Post(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Could not create %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(*resource.GetId())

	resourceAsWinGetApp, ok := resource.(models.WinGetAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Created resource is not of type WinGetApp: %s_%s", r.ProviderTypeName, r.TypeName),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, resourceAsWinGetApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

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

	tflog.Debug(ctx, fmt.Sprintf("Reading win get app with ID: %s", state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resource, err := r.client.DeviceAppManagement().MobileApps().ByMobileAppId(state.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		crud.HandleReadErrorIfNotFound(ctx, resp, r, &state, err)
		return
	}

	resourceAsWinGetApp, ok := resource.(models.WinGetAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Created resource is not of type WinGetApp: %s_%s", r.ProviderTypeName, r.TypeName),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, resourceAsWinGetApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *WinGetAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WinGetAppResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
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

	_, err = r.client.DeviceAppManagement().MobileApps().ByMobileAppId(plan.ID.ValueString()).Patch(ctx, requestBody, nil)
	if err != nil {
		crud.HandleUpdateErrorIfNotFound(ctx, resp, r, &plan, err)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *WinGetAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WinGetAppResourceModel

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

	err := r.client.DeviceAppManagement().MobileApps().ByMobileAppId(data.ID.ValueString()).Delete(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}
