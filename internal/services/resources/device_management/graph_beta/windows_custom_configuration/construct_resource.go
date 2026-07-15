package graphBetaWindowsCustomConfiguration

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the Windows custom configuration resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsCustomConfigurationResourceModel) (graphmodels.DeviceConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindows10CustomConfiguration()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	var settingModels []OmaSettingResourceModel
	diags := data.OmaSettings.ElementsAs(ctx, &settingModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to extract oma settings: %v", diags.Errors())
	}

	omaSettings := make([]graphmodels.OmaSettingable, 0, len(settingModels))
	for idx, settingModel := range settingModels {
		omaSetting, err := constructOmaSetting(settingModel)
		if err != nil {
			return nil, fmt.Errorf("failed to construct oma setting at index %d (%s): %s", idx, settingModel.OmaUri.ValueString(), err)
		}
		omaSettings = append(omaSettings, omaSetting)
	}
	requestBody.SetOmaSettings(omaSettings)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructOmaSetting constructs the correct microsoft.graph.omaSetting subtype from the
// Terraform model, converting the string value to the type expected by the Graph API.
func constructOmaSetting(data OmaSettingResourceModel) (graphmodels.OmaSettingable, error) {
	value := data.Value.ValueString()

	var omaSetting graphmodels.OmaSettingable

	switch data.OdataType.ValueString() {
	case "#microsoft.graph.omaSettingString":
		setting := graphmodels.NewOmaSettingString()
		setting.SetValue(&value)
		omaSetting = setting

	case "#microsoft.graph.omaSettingInteger":
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("value %q is not a valid integer: %s", value, err)
		}
		setting := graphmodels.NewOmaSettingInteger()
		int32Value := int32(intValue)
		setting.SetValue(&int32Value)
		omaSetting = setting

	case "#microsoft.graph.omaSettingBoolean":
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("value %q is not a valid boolean: %s", value, err)
		}
		setting := graphmodels.NewOmaSettingBoolean()
		setting.SetValue(&boolValue)
		omaSetting = setting

	case "#microsoft.graph.omaSettingBase64":
		setting := graphmodels.NewOmaSettingBase64()
		setting.SetValue(&value)
		convert.FrameworkToGraphString(data.FileName, setting.SetFileName)
		omaSetting = setting

	case "#microsoft.graph.omaSettingDateTime":
		timeValue, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return nil, fmt.Errorf("value %q is not a valid RFC3339 timestamp: %s", value, err)
		}
		setting := graphmodels.NewOmaSettingDateTime()
		setting.SetValue(&timeValue)
		omaSetting = setting

	case "#microsoft.graph.omaSettingFloatingPoint":
		floatValue, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, fmt.Errorf("value %q is not a valid floating point number: %s", value, err)
		}
		setting := graphmodels.NewOmaSettingFloatingPoint()
		float32Value := float32(floatValue)
		setting.SetValue(&float32Value)
		omaSetting = setting

	case "#microsoft.graph.omaSettingStringXml":
		setting := graphmodels.NewOmaSettingStringXml()
		setting.SetValue([]byte(value))
		convert.FrameworkToGraphString(data.FileName, setting.SetFileName)
		omaSetting = setting

	default:
		return nil, fmt.Errorf("unsupported oma setting odata type: %s", data.OdataType.ValueString())
	}

	convert.FrameworkToGraphString(data.DisplayName, omaSetting.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, omaSetting.SetDescription)
	convert.FrameworkToGraphString(data.OmaUri, omaSetting.SetOmaUri)

	return omaSetting, nil
}
