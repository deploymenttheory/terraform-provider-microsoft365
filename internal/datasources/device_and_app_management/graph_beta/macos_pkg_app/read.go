package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	resource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/device_and_app_management/graph_beta/macos_pkg_app"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Read handles the Read operation for MacOS PKG Application data sources.
//
// The function supports two methods of looking up a macOS PKG app:
// 1. By ID - Uses a direct API call to fetch the specific app
// 2. By DisplayName - Lists all mobile apps and finds the matching macOS PKG app
//
// The function ensures that:
// - Either ID or DisplayName is provided (but not both)
// - The lookup method is optimized based on the provided identifier
// - The remote state is properly mapped to the Terraform state
func (d *MacOSPKGAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object resource.MacOSPKGAppResourceModel

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
		// 1. Get base resource with expanded query to return categories
		requestParameters := &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetQueryParameters{
				Expand: []string{"categories"},
			},
		}

		respBaseResource, err := d.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(object.ID.ValueString()).
			Get(ctx, requestParameters)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		// This ensures type safety as the Graph API returns a base interface that needs
		// to be converted to the specific app type
		macOSPkgApp, ok := respBaseResource.(graphmodels.MacOSPkgAppable)
		if !ok {
			resp.Diagnostics.AddError(
				"Resource type mismatch",
				fmt.Sprintf("Expected resource of type MacOSPkgAppable but got %T", respBaseResource),
			)
			return
		}

		resource.MapRemoteResourceStateToTerraform(ctx, &object, macOSPkgApp)
	} else {
		// When looking up by display name, we need to list all mobile apps and filter
		mobileApps := d.client.
			DeviceAppManagement().
			MobileApps()

		// Use expanded query to get categories for all apps
		requestParameters := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
				Expand: []string{"categories"},
			},
		}

		result, err := mobileApps.Get(ctx, requestParameters)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		var foundApp graphmodels.MacOSPkgAppable
		for _, app := range result.GetValue() {
			// Check if app is a MacOS PKG app
			macOSApp, ok := app.(graphmodels.MacOSPkgAppable)
			if !ok {
				continue
			}

			if macOSApp.GetDisplayName() != nil && *macOSApp.GetDisplayName() == object.DisplayName.ValueString() {
				foundApp = macOSApp
				break
			}
		}

		if foundApp == nil {
			resp.Diagnostics.AddError(
				"Error Reading MacOS PKG App",
				fmt.Sprintf("No MacOS PKG app found with display name: %s", object.DisplayName.ValueString()),
			)
			return
		}

		resource.MapRemoteResourceStateToTerraform(ctx, &object, foundApp)
	}

	// Set the data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished MacOS PKG App Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
