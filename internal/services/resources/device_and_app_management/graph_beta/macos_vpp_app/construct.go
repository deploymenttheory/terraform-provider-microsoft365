package graphBetaMacOSVppApp

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
func constructResource(ctx context.Context, data *MacOSVppAppResourceModel) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	baseApp := graphmodels.NewMacOsVppApp()

	convert.FrameworkToGraphString(data.Description, baseApp.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, baseApp.SetPublisher)
	convert.FrameworkToGraphString(data.DisplayName, baseApp.SetDisplayName)
	convert.FrameworkToGraphString(data.InformationUrl, baseApp.SetInformationUrl)
	convert.FrameworkToGraphBool(data.IsFeatured, baseApp.SetIsFeatured)
	convert.FrameworkToGraphString(data.Owner, baseApp.SetOwner)
	convert.FrameworkToGraphString(data.Developer, baseApp.SetDeveloper)
	convert.FrameworkToGraphString(data.Notes, baseApp.SetNotes)
	convert.FrameworkToGraphString(data.PrivacyInformationUrl, baseApp.SetPrivacyInformationUrl)
	convert.FrameworkToGraphString(data.BundleId, baseApp.SetBundleId)
	convert.FrameworkToGraphString(data.VppTokenId, baseApp.SetVppTokenId)
	convert.FrameworkToGraphString(data.VppTokenAppleId, baseApp.SetVppTokenAppleId)
	convert.FrameworkToGraphString(data.VppTokenOrganizationName, baseApp.SetVppTokenOrganizationName)

	if err := convert.FrameworkToGraphEnum(data.VppTokenAccountType, graphmodels.ParseVppTokenAccountType, baseApp.SetVppTokenAccountType); err != nil {
		return nil, fmt.Errorf("failed to set VPP token account type: %s", err)
	}

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

	if data.LicensingType != nil {
		licensingType := graphmodels.NewVppLicensingType()
		convert.FrameworkToGraphBool(data.LicensingType.SupportUserLicensing, licensingType.SetSupportUserLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportDeviceLicensing, licensingType.SetSupportDeviceLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportsUserLicensing, licensingType.SetSupportsUserLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportsDeviceLicensing, licensingType.SetSupportsDeviceLicensing)
		baseApp.SetLicensingType(licensingType)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), baseApp); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return baseApp, nil
}
