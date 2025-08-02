package graphIOSMobileAppConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IOSMobileAppConfigurationDataSourceModel represents the data source model
type IOSMobileAppConfigurationDataSourceModel struct {
	Id                   types.String                          `tfsdk:"id"`
	DisplayName          types.String                          `tfsdk:"display_name"`
	Description          types.String                          `tfsdk:"description"`
	TargetedMobileApps   types.List                            `tfsdk:"targeted_mobile_apps"`
	EncodedSettingXml    types.String                          `tfsdk:"encoded_setting_xml"`
	Settings             []IOSMobileAppConfigurationSetting    `tfsdk:"settings"`
	Assignments          []IOSMobileAppConfigurationAssignment `tfsdk:"assignments"`
	CreatedDateTime      types.String                          `tfsdk:"created_date_time"`
	LastModifiedDateTime types.String                          `tfsdk:"last_modified_date_time"`
	Version              types.Int32                           `tfsdk:"version"`
	Timeouts             timeouts.Value                        `tfsdk:"timeouts"`
}

// IOSMobileAppConfigurationSetting represents a configuration setting
type IOSMobileAppConfigurationSetting struct {
	AppConfigKey      types.String `tfsdk:"app_config_key"`
	AppConfigKeyType  types.String `tfsdk:"app_config_key_type"`
	AppConfigKeyValue types.String `tfsdk:"app_config_key_value"`
}

// IOSMobileAppConfigurationAssignment represents an assignment
type IOSMobileAppConfigurationAssignment struct {
	Id     types.String                               `tfsdk:"id"`
	Target *IOSMobileAppConfigurationAssignmentTarget `tfsdk:"target"`
}

// IOSMobileAppConfigurationAssignmentTarget represents an assignment target
type IOSMobileAppConfigurationAssignmentTarget struct {
	ODataType types.String `tfsdk:"odata_type"`
	GroupId   types.String `tfsdk:"group_id"`
}
