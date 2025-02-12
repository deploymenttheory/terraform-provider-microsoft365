package graphBetaMacOSPKGApp

import (
	"context"
	"encoding/base64"
	"fmt"
	"path/filepath"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	utility "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_pkg"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *MacOSPKGAppResourceModel) (graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var baseApp graphmodels.MobileAppable

	baseApp = graphmodels.NewMacOSPkgApp()
	macOSApp := baseApp.(graphmodels.MacOSPkgAppable)

	// Set base properties
	constructors.SetStringProperty(data.Description, baseApp.SetDescription)
	constructors.SetStringProperty(data.Publisher, baseApp.SetPublisher)
	constructors.SetStringProperty(data.DisplayName, baseApp.SetDisplayName)
	constructors.SetStringProperty(data.InformationUrl, baseApp.SetInformationUrl)
	constructors.SetBoolProperty(data.IsFeatured, baseApp.SetIsFeatured)
	constructors.SetStringProperty(data.Owner, baseApp.SetOwner)
	constructors.SetStringProperty(data.Developer, baseApp.SetDeveloper)
	constructors.SetStringProperty(data.Notes, baseApp.SetNotes)

	if err := constructors.SetStringList(ctx, data.RoleScopeTagIds, baseApp.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if len(data.Categories) > 0 {
		categories := make([]graphmodels.MobileAppCategoryable, 0, len(data.Categories))
		for _, category := range data.Categories {
			mobileAppCategory := graphmodels.NewMobileAppCategory()
			constructors.SetStringProperty(category.ID, mobileAppCategory.SetId)
			constructors.SetStringProperty(category.DisplayName, mobileAppCategory.SetDisplayName)
			categories = append(categories, mobileAppCategory)
		}
		baseApp.SetCategories(categories)
	}

	if !data.LargeIcon.IsNull() {
		largeIcon := graphmodels.NewMimeContent()
		var iconData map[string]attr.Value
		data.LargeIcon.As(ctx, &iconData, basetypes.ObjectAsOptions{})

		iconType := "image/png"
		largeIcon.SetTypeEscaped(&iconType)

		if valueVal, ok := iconData["value"].(types.String); ok {
			iconBytes, err := base64.StdEncoding.DecodeString(valueVal.ValueString())
			if err != nil {
				return nil, fmt.Errorf("failed to decode icon base64: %v", err)
			}
			largeIcon.SetValue(iconBytes)
		}
		baseApp.SetLargeIcon(largeIcon)
	}

	// Set MacOS PKG specific properties
	if data.MacOSPkgApp.PackageInstallerFileSource.IsNull() || data.MacOSPkgApp.PackageInstallerFileSource.ValueString() == "" {
		return nil, fmt.Errorf("package_installer_file_source is required but not provided")
	}

	filename := filepath.Base(data.MacOSPkgApp.PackageInstallerFileSource.ValueString())
	constructors.SetStringProperty(types.StringValue(filename), macOSApp.SetFileName)

	fields := []utility.Field{
		{Key: "CFBundleIdentifier", Required: true},
		{Key: "CFBundleShortVersionString", Required: true},
	}

	// Extract fields from all Info.plist files in the package
	extractedFields, err := utility.ExtractFieldsFromFiles(
		ctx,
		data.MacOSPkgApp.PackageInstallerFileSource.ValueString(),
		"Info.plist",
		fields,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to extract metadata from pkg file: %w", err)
	}

	// First Info.plist becomes primary bundle
	primaryBundleId := extractedFields[0].Values["CFBundleIdentifier"]
	primaryBundleVersion := extractedFields[0].Values["CFBundleShortVersionString"]

	constructors.SetStringProperty(types.StringValue(primaryBundleId), macOSApp.SetPrimaryBundleId)
	constructors.SetStringProperty(types.StringValue(primaryBundleVersion), macOSApp.SetPrimaryBundleVersion)

	// All entries are set as included apps (including primary)
	var includedApps []graphmodels.MacOSIncludedAppable
	for _, fields := range extractedFields {
		includedApp := graphmodels.NewMacOSIncludedApp()
		constructors.SetStringProperty(
			types.StringValue(fields.Values["CFBundleIdentifier"]),
			includedApp.SetBundleId,
		)
		constructors.SetStringProperty(
			types.StringValue(fields.Values["CFBundleShortVersionString"]),
			includedApp.SetBundleVersion,
		)
		includedApps = append(includedApps, includedApp)
	}

	macOSApp.SetIncludedApps(includedApps)

	constructors.SetBoolProperty(data.MacOSPkgApp.IgnoreVersionDetection, macOSApp.SetIgnoreVersionDetection)

	if data.MacOSPkgApp.MinimumSupportedOperatingSystem != nil {
		minOS := graphmodels.NewMacOSMinimumOperatingSystem()
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
		constructors.SetBoolProperty(data.MacOSPkgApp.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
		macOSApp.SetMinimumSupportedOperatingSystem(minOS)
	}

	if data.MacOSPkgApp.PreInstallScript != nil {
		preScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.MacOSPkgApp.PreInstallScript.ScriptContent, preScript.SetScriptContent)
		macOSApp.SetPreInstallScript(preScript)
	}

	if data.MacOSPkgApp.PostInstallScript != nil {
		postScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.MacOSPkgApp.PostInstallScript.ScriptContent, postScript.SetScriptContent)
		macOSApp.SetPostInstallScript(postScript)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), macOSApp); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully constructed MacOSPkgApp resource with %d included apps", len(includedApps)))
	return macOSApp, nil
}
