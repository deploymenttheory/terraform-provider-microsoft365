package graphBetaLinuxPlatformScript

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteSettingsStateToTerraform restores the script and execution settings
// from Graph, including during import when no configured values are available.
func MapRemoteSettingsStateToTerraform(ctx context.Context, data *LinuxPlatformScriptResourceModel, settingsResponse []byte) error {
	parseNode, err := jsonserialization.NewJsonParseNode(settingsResponse)
	if err != nil {
		return fmt.Errorf("parse settings response: %w", err)
	}
	parsed, err := parseNode.GetObjectValue(graphmodels.CreateDeviceManagementConfigurationSettingCollectionResponseFromDiscriminatorValue)
	if err != nil {
		return fmt.Errorf("deserialize settings response: %w", err)
	}
	response, ok := parsed.(graphmodels.DeviceManagementConfigurationSettingCollectionResponseable)
	if !ok || response == nil {
		return fmt.Errorf("the API returned an invalid settings collection")
	}

	values := make(map[string]string, 4)
	for _, setting := range response.GetValue() {
		if setting == nil || setting.GetSettingInstance() == nil {
			return fmt.Errorf("the API returned a setting without an instance")
		}
		instance := setting.GetSettingInstance()
		if instance.GetSettingDefinitionId() == nil {
			return fmt.Errorf("the API returned a setting without a definition ID")
		}
		id := *instance.GetSettingDefinitionId()
		switch id {
		case "linux_customconfig_executioncontext", "linux_customconfig_executionfrequency", "linux_customconfig_executionretries":
			choice, ok := instance.(graphmodels.DeviceManagementConfigurationChoiceSettingInstanceable)
			if !ok || choice.GetChoiceSettingValue() == nil || choice.GetChoiceSettingValue().GetValue() == nil {
				return fmt.Errorf("the API returned an invalid choice for %s", id)
			}
			value := *choice.GetChoiceSettingValue().GetValue()
			prefix := id + "_"
			if !strings.HasPrefix(value, prefix) || len(value) == len(prefix) {
				return fmt.Errorf("the API returned an invalid choice value for %s", id)
			}
			values[id] = strings.TrimPrefix(value, prefix)
		case "linux_customconfig_script":
			simple, ok := instance.(graphmodels.DeviceManagementConfigurationSimpleSettingInstanceable)
			if !ok || simple.GetSimpleSettingValue() == nil {
				return fmt.Errorf("the API returned an invalid script setting")
			}
			value, ok := simple.GetSimpleSettingValue().(graphmodels.DeviceManagementConfigurationStringSettingValueable)
			if !ok || value.GetValue() == nil {
				return fmt.Errorf("the API returned an invalid script value")
			}
			decoded, err := base64.StdEncoding.DecodeString(*value.GetValue())
			if err != nil {
				return fmt.Errorf("decode script content: %w", err)
			}
			values[id] = string(decoded)
		}
	}
	for _, id := range []string{"linux_customconfig_executioncontext", "linux_customconfig_executionfrequency", "linux_customconfig_executionretries", "linux_customconfig_script"} {
		if _, ok := values[id]; !ok {
			return fmt.Errorf("the API response is missing required setting %s", id)
		}
	}

	data.ExecutionContext = types.StringValue(values["linux_customconfig_executioncontext"])
	data.ExecutionFrequency = types.StringValue(values["linux_customconfig_executionfrequency"])
	data.ExecutionRetries = types.StringValue(values["linux_customconfig_executionretries"])
	data.ScriptContent = types.StringValue(values["linux_customconfig_script"])
	tflog.Debug(ctx, "Finished mapping settings state to Terraform state")
	return nil
}
