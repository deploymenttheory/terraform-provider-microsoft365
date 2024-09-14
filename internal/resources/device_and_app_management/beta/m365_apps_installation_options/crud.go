package graphBetaM365AppsInstallationOptions

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
func (r *M365AppsInstallationOptionsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan M365AppsInstallationOptionsResourceModel

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

	// Since there's no Create API, we'll use the Update API to set the initial state
	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing M365AppsInstallationOptions",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	options, err := r.client.Admin().Microsoft365Apps().InstallationOptions().Patch(ctx, requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating M365AppsInstallationOptions",
			fmt.Sprintf("Could not create M365AppsInstallationOptions: %s", err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(*options.GetId())

	MapRemoteStateToTerraform(ctx, &plan, options)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *M365AppsInstallationOptionsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state M365AppsInstallationOptionsResourceModel
	tflog.Debug(ctx, "Starting Read method for M365AppsInstallationOptions")

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading M365 Apps Installation Options with ID: %s", state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, 30*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	options, err := r.client.Admin().Microsoft365Apps().InstallationOptions().Get(ctx, nil)
	if err != nil {
		crud.HandleReadErrorIfNotFound(ctx, resp, r, &state, err)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, options)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *M365AppsInstallationOptionsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan M365AppsInstallationOptionsResourceModel

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
			"Error constructing M365AppsInstallationOptions",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	options, err := r.client.Admin().Microsoft365Apps().InstallationOptions().Patch(ctx, requestBody, nil)
	if err != nil {
		crud.HandleUpdateErrorIfNotFound(ctx, resp, r, &plan, err)
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, options)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *M365AppsInstallationOptionsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data M365AppsInstallationOptionsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Since there's no Delete API, we'll just remove the resource from Terraform state
	tflog.Debug(ctx, fmt.Sprintf("Removing %s_%s from Terraform state", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
