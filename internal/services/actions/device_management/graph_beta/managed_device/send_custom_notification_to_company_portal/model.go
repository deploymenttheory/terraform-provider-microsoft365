// REF: https://learn.microsoft.com/en-us/graph/api/intune-devices-manageddevice-sendcustomnotificationtocompanyportal?view=graph-rest-beta
package graphBetaSendCustomNotificationToCompanyPortal

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/action/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SendCustomNotificationToCompanyPortalActionModel struct {
	ManagedDevices        []ManagedDeviceNotification   `tfsdk:"managed_devices"`
	ComanagedDevices      []ComanagedDeviceNotification `tfsdk:"comanaged_devices"`
	IgnorePartialFailures types.Bool                    `tfsdk:"ignore_partial_failures"`
	ValidateDeviceExists  types.Bool                    `tfsdk:"validate_device_exists"`
	Timeouts              timeouts.Value                `tfsdk:"timeouts"`
}

type ManagedDeviceNotification struct {
	DeviceID          types.String `tfsdk:"device_id"`
	NotificationTitle types.String `tfsdk:"notification_title"`
	NotificationBody  types.String `tfsdk:"notification_body"`
}

type ComanagedDeviceNotification struct {
	DeviceID          types.String `tfsdk:"device_id"`
	NotificationTitle types.String `tfsdk:"notification_title"`
	NotificationBody  types.String `tfsdk:"notification_body"`
}
