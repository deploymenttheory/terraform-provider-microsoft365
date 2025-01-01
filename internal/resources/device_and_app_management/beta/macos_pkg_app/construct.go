package graphBetaMacosPkgApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *MacOSPkgAppResourceModel) (graphmodels.MacOSPkgAppable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewMacOSPkgApp()

	// Set string properties using helper
	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetStringProperty(data.Publisher, requestBody.SetPublisher)
	constructors.SetStringProperty(data.PrivacyInformationUrl, requestBody.SetPrivacyInformationUrl)
	constructors.SetStringProperty(data.InformationUrl, requestBody.SetInformationUrl)
	constructors.SetStringProperty(data.Owner, requestBody.SetOwner)
	constructors.SetStringProperty(data.Developer, requestBody.SetDeveloper)
	constructors.SetStringProperty(data.Notes, requestBody.SetNotes)
	constructors.SetStringProperty(data.FileName, requestBody.SetFileName)
	constructors.SetStringProperty(data.PrimaryBundleId, requestBody.SetPrimaryBundleId)
	constructors.SetStringProperty(data.PrimaryBundleVersion, requestBody.SetPrimaryBundleVersion)

	// Set boolean properties using helper
	constructors.SetBoolProperty(data.IgnoreVersionDetection, requestBody.SetIgnoreVersionDetection)
	constructors.SetBoolProperty(data.IsFeatured, requestBody.SetIsFeatured)

	// Handle role scope tags
	if err := constructors.SetStringList(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %v", err)
	}

	// Handle large icon
	if !data.LargeIcon.Type.IsNull() && !data.LargeIcon.Value.IsNull() {
		largeIcon := graphmodels.NewMimeContent()
		largeIcon.SetValue([]byte(data.LargeIcon.Value.ValueString()))
		requestBody.SetLargeIcon(largeIcon)
	}

	// Handle included apps
	if len(data.IncludedApps) > 0 {
		includedApps := make([]graphmodels.MacOSIncludedAppable, 0, len(data.IncludedApps))
		for _, v := range data.IncludedApps {
			includedApp := graphmodels.NewMacOSIncludedApp()
			constructors.SetStringProperty(v.BundleId, includedApp.SetBundleId)
			constructors.SetStringProperty(v.BundleVersion, includedApp.SetBundleVersion)
			includedApps = append(includedApps, includedApp)
		}
		requestBody.SetIncludedApps(includedApps)
	}

	// Handle minimum OS version
	minOS := graphmodels.NewMacOSMinimumOperatingSystem()
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_7, minOS.SetV107)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_8, minOS.SetV108)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_9, minOS.SetV109)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_10, minOS.SetV1010)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_11, minOS.SetV1011)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_12, minOS.SetV1012)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_13, minOS.SetV1013)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_14, minOS.SetV1014)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V10_15, minOS.SetV1015)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V11_0, minOS.SetV110)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V12_0, minOS.SetV120)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V13_0, minOS.SetV130)
	constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V14_0, minOS.SetV140)
	requestBody.SetMinimumSupportedOperatingSystem(minOS)

	// Handle pre/post install scripts
	if !data.PreInstallScript.ScriptContent.IsNull() {
		preInstallScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PreInstallScript.ScriptContent, preInstallScript.SetScriptContent)
		requestBody.SetPreInstallScript(preInstallScript)
	}

	if !data.PostInstallScript.ScriptContent.IsNull() {
		postInstallScript := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PostInstallScript.ScriptContent, postInstallScript.SetScriptContent)
		requestBody.SetPostInstallScript(postInstallScript)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
