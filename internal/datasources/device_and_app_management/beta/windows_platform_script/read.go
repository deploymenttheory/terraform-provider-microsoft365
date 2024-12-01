package graphbetadevicemanagementscript

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	resource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/windows_platform_script"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	// object is the resource model for the device management script resource
	object resource.WindowsPlatformScriptResourceModel
)

// Read handles the Read operation for the WindowsPlatformScriptDataSource.
func (d *WindowsPlatformScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", d.ProviderTypeName, d.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, resource.ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Read base resource
	respResource, err := d.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	resource.MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	// Read assignments
	respAssignments, err := d.client.
		DeviceManagement().
		DeviceManagementScripts().
		ByDeviceManagementScriptId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
		return
	}

	resource.MapRemoteAssignmentStateToTerraform(ctx, &object, respAssignments)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
