package graphBetaIOSStoreApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	sharedConstructors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_and_app_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *IOSStoreAppResourceModel) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	baseApp := graphmodels.NewIosStoreApp()

	convert.FrameworkToGraphString(data.Description, baseApp.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, baseApp.SetPublisher)
	convert.FrameworkToGraphString(data.DisplayName, baseApp.SetDisplayName)
	convert.FrameworkToGraphString(data.InformationUrl, baseApp.SetInformationUrl)
	convert.FrameworkToGraphBool(data.IsFeatured, baseApp.SetIsFeatured)
	convert.FrameworkToGraphString(data.Owner, baseApp.SetOwner)
	convert.FrameworkToGraphString(data.Developer, baseApp.SetDeveloper)
	convert.FrameworkToGraphString(data.Notes, baseApp.SetNotes)
	convert.FrameworkToGraphString(data.PrivacyInformationUrl, baseApp.SetPrivacyInformationUrl)
	convert.FrameworkToGraphString(data.AppStoreUrl, baseApp.SetAppStoreUrl)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, baseApp.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Handle app icon (either from file path or web source)
	if data.AppIcon != nil {
		largeIcon, tempFiles, err := sharedConstructors.ConstructMobileAppIcon(ctx, data.AppIcon)
		if err != nil {
			return nil, err
		}

		defer func() {
			for _, tempFile := range tempFiles {
				helpers.CleanupTempFile(ctx, tempFile)
			}
		}()

		baseApp.SetLargeIcon(largeIcon)
	}

	// Set applicable device type
	if data.ApplicableDeviceType != nil {
		deviceType := graphmodels.NewIosDeviceType()
		convert.FrameworkToGraphBool(data.ApplicableDeviceType.IPad, deviceType.SetIPad)
		convert.FrameworkToGraphBool(data.ApplicableDeviceType.IPhoneAndIPod, deviceType.SetIPhoneAndIPod)
		baseApp.SetApplicableDeviceType(deviceType)
	}

	// Set minimum supported operating system
	if data.MinimumSupportedOperatingSystem != nil {
		minOS := graphmodels.NewIosMinimumOperatingSystem()
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V8_0, minOS.SetV80)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V9_0, minOS.SetV90)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V10_0, minOS.SetV100)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V11_0, minOS.SetV110)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V12_0, minOS.SetV120)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V13_0, minOS.SetV130)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V14_0, minOS.SetV140)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V15_0, minOS.SetV150)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V16_0, minOS.SetV160)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V17_0, minOS.SetV170)
		convert.FrameworkToGraphBool(data.MinimumSupportedOperatingSystem.V18_0, minOS.SetV180)
		baseApp.SetMinimumSupportedOperatingSystem(minOS)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), baseApp); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return baseApp, nil
}
