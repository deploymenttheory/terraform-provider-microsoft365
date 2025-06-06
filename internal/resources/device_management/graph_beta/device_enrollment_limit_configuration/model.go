// https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentlimitconfiguration?view=graph-rest-beta

package graphBetaDeviceEnrollmentLimitConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceEnrollmentLimitConfigurationResourceModel struct {
	ID                                types.String   `tfsdk:"id"`
	DisplayName                       types.String   `tfsdk:"display_name"`
	Description                       types.String   `tfsdk:"description"`
	Priority                          types.Int32    `tfsdk:"priority"`
	CreatedDateTime                   types.String   `tfsdk:"created_date_time"`
	LastModifiedDateTime              types.String   `tfsdk:"last_modified_date_time"`
	Version                           types.Int32    `tfsdk:"version"`
	RoleScopeTagIds                   types.Set      `tfsdk:"role_scope_tag_ids"`
	DeviceEnrollmentConfigurationType types.String   `tfsdk:"device_enrollment_configuration_type"`
	Limit                             types.Int32    `tfsdk:"limit"`
	Timeouts                          timeouts.Value `tfsdk:"timeouts"`
}
