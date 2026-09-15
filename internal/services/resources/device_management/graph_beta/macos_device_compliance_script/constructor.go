package graphBetaMacosDeviceComplianceScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource creates a new DeviceComplianceScript based on the resource model.
//
// This is the same graphmodels.DeviceComplianceScript object the Windows
// device compliance script resource uses. The one thing that resource never
// does -- because it didn't need to when it was the only platform -- is set
// the Platform property. We set it explicitly to macOS here so the API
// associates this script with macOS custom compliance policies rather than
// defaulting to windows10.
func constructResource(ctx context.Context, data *DeviceComplianceScriptResourceModel) (graphmodels.DeviceComplianceScriptable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceComplianceScript()

	platform := graphmodels.MACOS_DEVICECOMPLIANCESCRIPTPLATFORMTYPE
	requestBody.SetPlatform(&platform)

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.Publisher, requestBody.SetPublisher)
	convert.FrameworkToGraphBool(data.EnforceSignatureCheck, requestBody.SetEnforceSignatureCheck)

	if err := convert.FrameworkToGraphEnum(data.RunAsAccount, graphmodels.ParseRunAsAccountType, requestBody.SetRunAsAccount); err != nil {
		return nil, fmt.Errorf("invalid run as account type: %s", err)
	}

	// Bash script content -- same field Windows uses for its PowerShell
	// content, just bash bytes instead of PowerShell bytes.
	convert.FrameworkToGraphBytes(data.DetectionScriptContent, requestBody.SetDetectionScriptContent)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}