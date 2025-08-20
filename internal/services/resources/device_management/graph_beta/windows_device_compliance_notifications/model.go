// REF: https://learn.microsoft.com/en-au/graph/api/resources/intune-notification-notificationmessagetemplate?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-au/graph/api/intune-notification-notificationmessagetemplate-create?view=graph-rest-beta
package graphBetaWindowsDeviceComplianceNotifications

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsDeviceComplianceNotificationsResourceModel struct {
	ID                            types.String   `tfsdk:"id"`
	DisplayName                   types.String   `tfsdk:"display_name"`
	DefaultLocale                 types.String   `tfsdk:"default_locale"`
	BrandingOptions               types.Set      `tfsdk:"branding_options"`
	RoleScopeTagIds               types.Set      `tfsdk:"role_scope_tag_ids"`
	LastModifiedDateTime          types.String   `tfsdk:"last_modified_date_time"`
	LocalizedNotificationMessages types.Set      `tfsdk:"localized_notification_messages"`
	Timeouts                      timeouts.Value `tfsdk:"timeouts"`
}

type LocalizedNotificationMessageModel struct {
	ID              types.String `tfsdk:"id"`
	Locale          types.String `tfsdk:"locale"`
	Subject         types.String `tfsdk:"subject"`
	MessageTemplate types.String `tfsdk:"message_template"`
	IsDefault       types.Bool   `tfsdk:"is_default"`
}
