package graphBetaDeviceCategory

import (
	"context"
	"fmt"

	resource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/device_category"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for Device Category data sources.
//
// The function supports two methods of looking up a device category:
// 1. By ID - Uses a direct API call to fetch the specific category
// 2. By DisplayName - Lists all categories and finds the matching one
//
// The function ensures that:
// - Either ID or DisplayName is provided (but not both)
// - The lookup method is optimized based on the provided identifier
// - The remote state is properly mapped to the Terraform state
func (d *DeviceCategoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object resource.DeviceCategoryResourceModel

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
		category, err := d.client.
			DeviceManagement().
			DeviceCategories().
			ByDeviceCategoryId(object.ID.ValueString()).
			Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Device Category",
				fmt.Sprintf("Could not read device category ID %s: %s", object.ID.ValueString(), err),
			)
			return
		}

		resource.MapRemoteStateToTerraform(ctx, &object, category)
	} else {
		categories := d.client.
			DeviceManagement().
			DeviceCategories()

		result, err := categories.Get(ctx, nil)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error Reading Device Categories",
				fmt.Sprintf("Could not read device categories: %s", err),
			)
			return
		}

		var foundCategory graphmodels.DeviceCategoryable
		for _, category := range result.GetValue() {
			if category.GetDisplayName() != nil && *category.GetDisplayName() == object.DisplayName.ValueString() {
				foundCategory = category
				break
			}
		}

		if foundCategory == nil {
			resp.Diagnostics.AddError(
				"Error Reading Device Category",
				fmt.Sprintf("No device category found with display name: %s", object.DisplayName.ValueString()),
			)
			return
		}

		resource.MapRemoteStateToTerraform(ctx, &object, foundCategory)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
