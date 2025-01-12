package sharedmodels

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ReuseablePolicySettingsResourceModel holds the configuration for a Settings Catalog profile.
type ReuseablePolicySettingsResourceModel struct {
	ID                                  types.String   `tfsdk:"id"`
	DisplayName                         types.String   `tfsdk:"display_name"`
	Description                         types.String   `tfsdk:"description"`
	CreatedDateTime                     types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime                types.String   `tfsdk:"last_modified_date_time"`
	ReferencingConfigurationPolicies    types.List     `tfsdk:"referencing_configuration_policies"`
	ReferencingConfigurationPolicyCount types.Int32    `tfsdk:"referencing_configuration_policy_count"`
	Settings                            types.String   `tfsdk:"settings"`
	Version                             types.Int32    `tfsdk:"version"`
	Timeouts                            timeouts.Value `tfsdk:"timeouts"`
}
