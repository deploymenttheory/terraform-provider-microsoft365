package graphBetaRoleScopeTag

import (
	"context"
	"fmt"

	resource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/beta/role_scope_tag"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Role Scope Tag data sources.
//
// The function supports two methods of looking up an Role Scope Tag:
// 1. By ID - Uses a direct API call to fetch the specific filter
// 2. By DisplayName - Lists all filters and finds the matching one
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
//   - Retrieves all filters and searches for matching display name
//   - Returns error if no matching filter is found
func (d *RoleScopeTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object resource.RoleScopeTagResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
		filter, err := d.client.
			DeviceManagement().
			RoleScopeTags().
			ByRoleScopeTagId(object.ID.ValueString()).
			Get(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Role Scope Tag",
				fmt.Sprintf("Could not read Role Scope Tag ID %s: %s", object.ID.ValueString(), err),
			)
			return
		}
		resource.MapRemoteResourceStateToTerraform(ctx, &object, filter)
	} else {
		filters := d.client.
			DeviceManagement().
			RoleScopeTags()
		result, err := filters.Get(ctx, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Role Scope Tags",
				fmt.Sprintf("Could not read Role Scope Tags: %s", err),
			)
			return
		}

		var foundFilter graphmodels.RoleScopeTagable
		for _, filter := range result.GetValue() {
			if *filter.GetDisplayName() == object.DisplayName.ValueString() {
				foundFilter = filter
				break
			}
		}

		if foundFilter == nil {
			resp.Diagnostics.AddError(
				"Error Reading Role Scope Tag",
				fmt.Sprintf("No Role Scope Tag found with display name: %s", object.DisplayName.ValueString()),
			)
			return
		}

		resource.MapRemoteResourceStateToTerraform(ctx, &object, foundFilter)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
