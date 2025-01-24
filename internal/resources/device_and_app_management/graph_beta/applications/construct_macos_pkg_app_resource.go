package graphBetaApplications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructMacOSPkgAppResource(ctx context.Context, data *MacOSPkgAppResourceModel, baseApp graphmodels.MacOSPkgAppable) (graphmodels.MacOSPkgAppable, error) {
	constructors.SetBoolProperty(data.IgnoreVersionDetection, baseApp.SetIgnoreVersionDetection)
	constructors.SetStringProperty(data.PrimaryBundleId, baseApp.SetPrimaryBundleId)
	constructors.SetStringProperty(data.PrimaryBundleVersion, baseApp.SetPrimaryBundleVersion)

	if len(data.IncludedApps) > 500 {
		return nil, fmt.Errorf("included_apps exceeds maximum limit of 500")
	}

	if len(data.IncludedApps) > 0 {
		includedApps := make([]graphmodels.MacOSIncludedAppable, 0, len(data.IncludedApps))
		for _, app := range data.IncludedApps {
			includedApp := graphmodels.NewMacOSIncludedApp()
			constructors.SetStringProperty(app.BundleId, includedApp.SetBundleId)
			constructors.SetStringProperty(app.BundleVersion, includedApp.SetBundleVersion)
			includedApps = append(includedApps, includedApp)
		}
		baseApp.SetIncludedApps(includedApps)
	}

	if data.MinimumSupportedOperatingSystem != nil {
		minOS := graphmodels.NewMacOSMinimumOperatingSystem()
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V107, minOS.SetV107)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V108, minOS.SetV108)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V109, minOS.SetV109)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1010, minOS.SetV1010)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1011, minOS.SetV1011)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1012, minOS.SetV1012)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1013, minOS.SetV1013)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1014, minOS.SetV1014)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V1015, minOS.SetV1015)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V110, minOS.SetV110)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V120, minOS.SetV120)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V130, minOS.SetV130)
		constructors.SetBoolProperty(data.MinimumSupportedOperatingSystem.V140, minOS.SetV140)
		baseApp.SetMinimumSupportedOperatingSystem(minOS)
	}

	if data.PreInstallScript != nil {
		script := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PreInstallScript.ScriptContent, script.SetScriptContent)
		baseApp.SetPreInstallScript(script)
	}

	if data.PostInstallScript != nil {
		script := graphmodels.NewMacOSAppScript()
		constructors.SetStringProperty(data.PostInstallScript.ScriptContent, script.SetScriptContent)
		baseApp.SetPostInstallScript(script)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s MacOSPkgApp resource values", ResourceName))

	return baseApp, nil
}
