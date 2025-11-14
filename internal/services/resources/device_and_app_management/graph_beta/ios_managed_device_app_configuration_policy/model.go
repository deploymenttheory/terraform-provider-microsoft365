// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-apps-iosmobileappconfiguration?view=graph-rest-beta
package graphBetaIOSManagedDeviceAppConfigurationPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IOSManagedDeviceAppConfigurationPolicyResourceModel struct {
	ID                 types.String   `tfsdk:"id"`
	DisplayName        types.String   `tfsdk:"display_name"`
	Description        types.String   `tfsdk:"description"`
	TargetedMobileApps types.Set      `tfsdk:"targeted_mobile_apps"`
	RoleScopeTagIds    types.Set      `tfsdk:"role_scope_tag_ids"`
	Version            types.Int32    `tfsdk:"version"`
	EncodedSettingXml  types.String   `tfsdk:"encoded_setting_xml"`
	Settings           types.Set      `tfsdk:"settings"`
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
}

type AppConfigurationSettingItemModel struct {
	AppConfigKey      types.String `tfsdk:"app_config_key"`
	AppConfigKeyType  types.String `tfsdk:"app_config_key_type"`
	AppConfigKeyValue types.String `tfsdk:"app_config_key_value"`
}
