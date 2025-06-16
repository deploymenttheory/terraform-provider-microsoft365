package graphBetaDeviceManagementScript

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Windows Platform Script data sources.
//
// The function supports two methods of looking up a Windows Platform Script:
// 1. By ID - Uses a direct API call to fetch the specific script
// 2. By DisplayName - Lists all scripts and finds the matching one
//
// The function ensures that:
// - Either ID or DisplayName is provided (but not both)
// - The lookup method is optimized based on the provided identifier
// - The remote state is properly mapped to the Terraform state
//
// The function will:
//  1. Extract and validate the configuration
//  2. Verify that exactly one identifier (ID or DisplayName) is provided
//  3. Perform the appropriate API call based on the provided identifier
//  4. Map the remote state to the Terraform state
//  5. Handle any errors and return appropriate diagnostics
//
// If using ID:
//   - Makes a direct GET request to the specific resource endpoint
//   - Returns error if the ID is not found
//
// If using DisplayName:
//   - Retrieves all scripts and searches for matching display name
//   - Returns error if no matching script is found
func (d *WindowsPlatformScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object WindowsPlatformScriptDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", d.ProviderTypeName, d.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Validate that either ID or display_name is provided, but not both
	if object.ID.IsNull() && object.DisplayName.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Either id or display_name must be provided",
		)
		return
	}
	if !object.ID.IsNull() && !object.DisplayName.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"Only one of id or display_name should be provided, not both",
		)
		return
	}

	if !object.ID.IsNull() {
		// Direct lookup by ID
		respResource, err := d.client.
			DeviceManagement().
			DeviceManagementScripts().
			ByDeviceManagementScriptId(object.ID.ValueString()).
			Get(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Platform Script",
				fmt.Sprintf("Could not read Windows Platform Script ID %s: %s", object.ID.ValueString(), err),
			)
			return
		}

		MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	} else {
		// Lookup by display name
		result, err := d.client.
			DeviceManagement().
			DeviceManagementScripts().
			Get(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Platform Scripts",
				fmt.Sprintf("Could not read Windows Platform Scripts: %s", err),
			)
			return
		}

		var foundScript graphmodels.DeviceManagementScriptable
		for _, script := range result.GetValue() {
			if *script.GetDisplayName() == object.DisplayName.ValueString() {
				foundScript = script
				break
			}
		}

		if foundScript == nil {
			resp.Diagnostics.AddError(
				"Error Reading Windows Platform Script",
				fmt.Sprintf("No Windows Platform Script found with display name: %s", object.DisplayName.ValueString()),
			)
			return
		}

		MapRemoteResourceStateToTerraform(ctx, &object, foundScript)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
