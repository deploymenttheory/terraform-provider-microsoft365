package graphBetaDeviceEnrollmentNotification

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceEnrollmentNotificationResourceModel struct {
	ID                                types.String   `tfsdk:"id"`
	DisplayName                       types.String   `tfsdk:"display_name"`
	Description                       types.String   `tfsdk:"description"`
	PlatformType                      types.String   `tfsdk:"platform_type"`
	DefaultLocale                     types.String   `tfsdk:"default_locale"`
	RoleScopeTagIds                   types.Set      `tfsdk:"role_scope_tag_ids"`
	BrandingOptions                   types.Set      `tfsdk:"branding_options"`
	NotificationTemplates             types.Set      `tfsdk:"notification_templates"`
	Priority                          types.Int32    `tfsdk:"priority"`
	DeviceEnrollmentConfigurationType types.String   `tfsdk:"device_enrollment_configuration_type"`
	LocalizedNotificationMessages     types.Set      `tfsdk:"localized_notification_messages"`
	Assignments                       types.Set      `tfsdk:"assignments"`
	Timeouts                          timeouts.Value `tfsdk:"timeouts"`
}

type LocalizedNotificationMessageModel struct {
	Locale          types.String `tfsdk:"locale"`
	Subject         types.String `tfsdk:"subject"`
	MessageTemplate types.String `tfsdk:"message_template"`
	IsDefault       types.Bool   `tfsdk:"is_default"`
	TemplateType    types.String `tfsdk:"template_type"`
}
