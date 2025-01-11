package sharedmodels

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReuseablePolicySettingsResourceModel holds the configuration for a Settings Catalog profile.
type ReuseablePolicySettingsResourceModel struct {
	ID                                  types.String   `tfsdk:"id"`
	Name                                types.String   `tfsdk:"name"` // Maps to DisplayName in SDK
	Description                         types.String   `tfsdk:"description"`
	CreatedDateTime                     types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime                types.String   `tfsdk:"last_modified_date_time"`
	ReferencingConfigurationPolicies    types.List     `tfsdk:"referencing_configuration_policies"`
	ReferencingConfigurationPolicyCount types.Int64    `tfsdk:"referencing_configuration_policy_count"`
	SettingDefinitionId                 types.String   `tfsdk:"setting_definition_id"`
	SettingInstance                     types.String   `tfsdk:"settings"` // Maps to Settings in previous model
	Version                             types.Int64    `tfsdk:"version"`
	Timeouts                            timeouts.Value `tfsdk:"timeouts"`
}
