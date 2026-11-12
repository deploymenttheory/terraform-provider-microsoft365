package graphBetaIOSManagedDeviceAppConfigurationPolicy

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	serialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

var configurationSettingsType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"app_config_key": types.StringType, "app_config_key_type": types.StringType,
	"app_config_key_value": types.StringType,
}}

func TestUnitIOSConfigurationXMLClearsCustomSettingsInRequest(t *testing.T) {
	model := IOSManagedDeviceAppConfigurationPolicyResourceModel{
		EncodedSettingXml:  types.StringValue("<dict><key>testKey</key><string>testValue</string></dict>"),
		Settings:           types.SetNull(configurationSettingsType),
		TargetedMobileApps: types.SetNull(types.StringType),
		RoleScopeTagIds:    types.SetNull(types.StringType),
	}
	request, err := constructResource(context.Background(), nil, &model)
	require.NoError(t, err)
	writer := serialization.NewJsonSerializationWriter()
	require.NoError(t, writer.WriteObjectValue("", request))
	serialized, err := writer.GetSerializedContent()
	require.NoError(t, err)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(serialized, &payload))
	settings, present := payload["settings"]
	require.True(t, present, "Graph requires an explicit empty collection to switch an existing policy to XML")
	require.Equal(t, []any{}, settings)
}

func TestUnitIOSConfigurationRemoteEmptySettingsClearPreviousValues(t *testing.T) {
	previous := types.SetValueMust(configurationSettingsType, []attr.Value{
		types.ObjectValueMust(configurationSettingsType.AttrTypes, map[string]attr.Value{
			"app_config_key": types.StringValue("testKey"), "app_config_key_type": types.StringValue("stringType"),
			"app_config_key_value": types.StringValue("testValue"),
		}),
	})
	for _, testCase := range []struct {
		name     string
		settings types.Set
		null     bool
	}{
		{name: "previous custom settings", settings: previous, null: true},
		{name: "unset", settings: types.SetNull(configurationSettingsType), null: true},
		{name: "unknown", settings: types.SetUnknown(configurationSettingsType), null: true},
		{name: "explicit empty collection", settings: types.SetValueMust(configurationSettingsType, []attr.Value{}), null: false},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			remote := graphmodels.NewIosMobileAppConfiguration()
			remote.SetEncodedSettingXml([]byte("<dict/>"))
			remote.SetSettings([]graphmodels.AppConfigurationSettingItemable{})
			model := IOSManagedDeviceAppConfigurationPolicyResourceModel{Settings: testCase.settings}
			MapRemoteStateToTerraform(context.Background(), &model, remote)
			require.Equal(t, testCase.null, model.Settings.IsNull())
			require.Empty(t, model.Settings.Elements())
			require.Equal(t, "<dict/>", model.EncodedSettingXml.ValueString())
		})
	}
}
