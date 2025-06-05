// https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deviceenrollmentnotificationconfiguration?view=graph-rest-beta

package graphBetaDeviceEnrollmentNotificationConfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceEnrollmentNotificationConfigurationResourceModel struct {
	ID                                types.String                       `tfsdk:"id"`
	DisplayName                       types.String                       `tfsdk:"display_name"`
	Description                       types.String                       `tfsdk:"description"`
	Priority                          types.Int32                        `tfsdk:"priority"`
	CreatedDateTime                   types.String                       `tfsdk:"created_date_time"`
	LastModifiedDateTime              types.String                       `tfsdk:"last_modified_date_time"`
	Version                           types.Int32                        `tfsdk:"version"`
	RoleScopeTagIds                   types.Set                          `tfsdk:"role_scope_tag_ids"`
	DeviceEnrollmentConfigurationType types.String                       `tfsdk:"device_enrollment_configuration_type"`
	PlatformType                      types.String                       `tfsdk:"platform_type"`
	TemplateTypes                     types.Set                          `tfsdk:"template_types"`
	NotificationMessageTemplateId     types.String                       `tfsdk:"notification_message_template_id"`
	NotificationTemplates             types.Set                          `tfsdk:"notification_templates"`
	BrandingOptions                   types.String                       `tfsdk:"branding_options"`
	PushLocalizedMessage              *LocalizedNotificationMessageModel `tfsdk:"push_localized_message"`
	EmailLocalizedMessage             *LocalizedNotificationMessageModel `tfsdk:"email_localized_message"`
	Timeouts                          timeouts.Value                     `tfsdk:"timeouts"`
}

type LocalizedNotificationMessageModel struct {
	Locale          types.String `tfsdk:"locale"`
	Subject         types.String `tfsdk:"subject"`
	MessageTemplate types.String `tfsdk:"message_template"`
	IsDefault       types.Bool   `tfsdk:"is_default"`
}
