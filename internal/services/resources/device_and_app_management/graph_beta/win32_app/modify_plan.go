package graphBetaWin32App

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ModifyPlan invalidates the computed content version when its source changes.
// The shared content schema otherwise preserves the previous committed version.
func (r *Win32LobAppResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan, state Win32LobAppResourceModel
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("app_installer"), &plan.AppInstaller)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("app_installer_zip"), &plan.AppInstallerZip)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if req.State.Raw.IsNull() {
		if !hasInstallerSource(&plan) {
			resp.Diagnostics.AddError("Missing installer source", "Creating a Win32 application requires either app_installer or app_installer_zip.")
		}
		return
	}
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("app_installer"), &state.AppInstaller)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("app_installer_zip"), &state.AppInstallerZip)...)
	if resp.Diagnostics.HasError() || !hasInstallerSource(&plan) || !installerSourceChanged(&plan, &state) {
		return
	}
	var configured, planned types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("content_version"), &configured)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("content_version"), &planned)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !configured.IsNull() {
		resp.Diagnostics.AddAttributeError(path.Root("content_version"), "Content version is managed by the provider", "Omit content_version when changing the installer source; the new version is assigned by Intune.")
		return
	}
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("content_version"), types.ListUnknown(planned.ElementType(ctx)))...)
}
