package graphBetaAllowNextEnrollment

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AllowNextEnrollmentActionModel represents the Terraform configuration for the allow next enrollment action
type AllowNextEnrollmentActionModel struct {
	WindowsAutopilotDeviceIdentityID types.String   `tfsdk:"windows_autopilot_device_identity_id"`
	Timeouts                         timeouts.Value `tfsdk:"timeouts"`
}
