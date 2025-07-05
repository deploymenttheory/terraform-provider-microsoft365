package graphBetaMacOSVppApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
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

	// Set VPP token account type
	if !data.VppTokenAccountType.IsNull() && data.VppTokenAccountType.ValueString() != "" {
		vppTokenAccountType := data.VppTokenAccountType.ValueString()
		baseApp.SetVppTokenAccountType(&vppTokenAccountType)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, baseApp.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Handle large icon if provided
	if data.LargeIcon != nil {
		largeIcon := graphmodels.NewMimeContent()
		if !data.LargeIcon.Type.IsNull() {
			convert.FrameworkToGraphString(data.LargeIcon.Type, largeIcon.SetTypeEscaped)
		}
		if !data.LargeIcon.Value.IsNull() {
			value := []byte(data.LargeIcon.Value.ValueString())
			largeIcon.SetValue(value)
		}
		baseApp.SetLargeIcon(largeIcon)
	}

	// Set licensing type if provided
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

// constructUpdateResource maps the Terraform schema to the SDK model for updates
func constructUpdateResource(ctx context.Context, data *MacOSVppAppResourceModel) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing update for %s resource from model", ResourceName))

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

	// Set VPP token account type
	if !data.VppTokenAccountType.IsNull() && data.VppTokenAccountType.ValueString() != "" {
		vppTokenAccountType := data.VppTokenAccountType.ValueString()
		baseApp.SetVppTokenAccountType(&vppTokenAccountType)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, baseApp.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Handle large icon if provided
	if data.LargeIcon != nil {
		largeIcon := graphmodels.NewMimeContent()
		if !data.LargeIcon.Type.IsNull() {
			convert.FrameworkToGraphString(data.LargeIcon.Type, largeIcon.SetTypeEscaped)
		}
		if !data.LargeIcon.Value.IsNull() {
			value := []byte(data.LargeIcon.Value.ValueString())
			largeIcon.SetValue(value)
		}
		baseApp.SetLargeIcon(largeIcon)
	}

	// Set licensing type if provided
	if data.LicensingType != nil {
		licensingType := graphmodels.NewVppLicensingType()
		convert.FrameworkToGraphBool(data.LicensingType.SupportUserLicensing, licensingType.SetSupportUserLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportDeviceLicensing, licensingType.SetSupportDeviceLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportsUserLicensing, licensingType.SetSupportsUserLicensing)
		convert.FrameworkToGraphBool(data.LicensingType.SupportsDeviceLicensing, licensingType.SetSupportsDeviceLicensing)
		baseApp.SetLicensingType(licensingType)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource update %s", ResourceName), baseApp); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing update for %s resource", ResourceName))

	return baseApp, nil
}
