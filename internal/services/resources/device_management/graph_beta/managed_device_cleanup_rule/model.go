// Base resource REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-devices-manageddevicecleanuprule?view=graph-rest-beta
package graphBetaManagedDeviceCleanupRule

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ManagedDeviceCleanupRuleResourceModel struct {
	ID                                     types.String   `tfsdk:"id"`
	DisplayName                            types.String   `tfsdk:"display_name"`
	Description                            types.String   `tfsdk:"description"`
	DeviceCleanupRulePlatformType          types.String   `tfsdk:"device_cleanup_rule_platform_type"`
	DeviceInactivityBeforeRetirementInDays types.Int32    `tfsdk:"device_inactivity_before_retirement_in_days"`
	LastModifiedDateTime                   types.String   `tfsdk:"last_modified_date_time"`
	Timeouts                               timeouts.Value `tfsdk:"timeouts"`
}
