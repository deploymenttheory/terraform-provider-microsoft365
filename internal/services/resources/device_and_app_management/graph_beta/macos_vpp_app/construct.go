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

	requestBody := graphmodels.NewMacOsVppApp()

	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, requestBody.SetPublisher)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.InformationUrl, requestBody.SetInformationUrl)
	convert.FrameworkToGraphBool(data.IsFeatured, requestBody.SetIsFeatured)
	convert.FrameworkToGraphString(data.Owner, requestBody.SetOwner)
	convert.FrameworkToGraphString(data.Developer, requestBody.SetDeveloper)
	convert.FrameworkToGraphString(data.Notes, requestBody.SetNotes)
	convert.FrameworkToGraphString(data.PrivacyInformationUrl, requestBody.SetPrivacyInformationUrl)
	convert.FrameworkToGraphString(data.BundleId, requestBody.SetBundleId)
	convert.FrameworkToGraphString(data.VppTokenId, requestBody.SetVppTokenId)
	convert.FrameworkToGraphString(data.VppTokenAppleId, requestBody.SetVppTokenAppleId)
	convert.FrameworkToGraphString(data.VppTokenOrganizationName, requestBody.SetVppTokenOrganizationName)

	if err := convert.FrameworkToGraphEnum(data.VppTokenAccountType, graphmodels.ParseVppTokenAccountType, requestBody.SetVppTokenAccountType); err != nil {
		return nil, fmt.Errorf("failed to set VPP token account type: %s", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
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

		requestBody.SetLargeIcon(largeIcon)
	}

	if data.LicensingType != nil {
		licensingType := graphmodels.NewVppLicensingType()
		convert.FrameworkToGraphBool(data.LicensingType.SupportUserLicensing, licensingType.SetSupportUserLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportDeviceLicensing, licensingType.SetSupportDeviceLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportsUserLicensing, licensingType.SetSupportsUserLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportsDeviceLicensing, licensingType.SetSupportsDeviceLicensing)
		requestBody.SetLicensingType(licensingType)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
