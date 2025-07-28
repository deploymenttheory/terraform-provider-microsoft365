package graphBetaLinuxDeviceComplianceScript

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

const (
	LinuxComplianceScriptSettingDefinitionId = "linux_customcompliance_discoveryscript_reusablesetting"
)

// constructResource creates a new DeviceManagementReusablePolicySetting based on the resource model
func constructResource(ctx context.Context, data *LinuxDeviceComplianceScriptResourceModel) (graphmodels.DeviceManagementReusablePolicySettingable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewDeviceManagementReusablePolicySetting()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	settingDefinitionId := LinuxComplianceScriptSettingDefinitionId
	requestBody.SetSettingDefinitionId(&settingDefinitionId)

	version := int32(0)
	requestBody.SetVersion(&version)

	settingInstance := graphmodels.NewDeviceManagementConfigurationSimpleSettingInstance()
	settingInstance.SetSettingDefinitionId(&settingDefinitionId)
	odataTypeSettingInstance := "#microsoft.graph.deviceManagementConfigurationSimpleSettingInstance"
	settingInstance.SetOdataType(&odataTypeSettingInstance)

	simpleSettingValue := graphmodels.NewDeviceManagementConfigurationStringSettingValue()
	odataTypeStringValue := "#microsoft.graph.deviceManagementConfigurationStringSettingValue"
	simpleSettingValue.SetOdataType(&odataTypeStringValue)

	if !data.DetectionScriptContent.IsNull() && !data.DetectionScriptContent.IsUnknown() {
		scriptContent := data.DetectionScriptContent.ValueString()
		encodedContent := base64.StdEncoding.EncodeToString([]byte(scriptContent))
		simpleSettingValue.SetValue(&encodedContent)
	}

	settingInstance.SetSimpleSettingValue(simpleSettingValue)
	requestBody.SetSettingInstance(settingInstance)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))
	return requestBody, nil
}
