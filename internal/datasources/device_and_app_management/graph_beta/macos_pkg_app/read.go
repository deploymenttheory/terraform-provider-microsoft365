package graphBetaMacOSPKGApp

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
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
	var object MacOSPKGAppResourceModel

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

		macOSPkgApp, ok := respBaseResource.(graphmodels.MacOSPkgAppable)
		if !ok {
			resp.Diagnostics.AddError(
				"Resource type mismatch",
				fmt.Sprintf("Expected resource of type MacOSPkgAppable but got %T", respBaseResource),
			)
			return
		}

		MapRemoteResourceStateToTerraform(ctx, &object, macOSPkgApp)
	} else {
		mobileApps := d.client.
			DeviceAppManagement().
			MobileApps()

		result, err := mobileApps.Get(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		var appId string
		for _, app := range result.GetValue() {
			if app.GetDisplayName() != nil && *app.GetDisplayName() == object.DisplayName.ValueString() {
				appId = *app.GetId()
				break
			}
		}

		if appId == "" {
			resp.Diagnostics.AddError(
				"Error Reading MacOS PKG App",
				fmt.Sprintf("No mobile app found with display name: %s", object.DisplayName.ValueString()),
			)
			return
		}

		requestParameters := &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.MobileAppsMobileAppItemRequestBuilderGetQueryParameters{
				Expand: []string{"categories"},
			},
		}

		respBaseResource, err := d.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(appId).
			Get(ctx, requestParameters)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Read", d.ReadPermissions)
			return
		}

		macOSPkgApp, ok := respBaseResource.(graphmodels.MacOSPkgAppable)
		if !ok {
			resp.Diagnostics.AddError(
				"Resource type mismatch",
				fmt.Sprintf("Expected resource of type MacOSPkgAppable but got %T", respBaseResource),
			)
			return
		}

		MapRemoteResourceStateToTerraform(ctx, &object, macOSPkgApp)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished MacOS PKG App Datasource Read Method: %s_%s", d.ProviderTypeName, d.TypeName))
}
