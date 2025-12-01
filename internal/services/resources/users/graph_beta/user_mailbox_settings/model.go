// REF: https://learn.microsoft.com/en-us/graph/api/user-get-mailboxsettings?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/user-update-mailboxsettings?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/outlookuser-supportedlanguages?view=graph-rest-beta&tabs=http
// REF: https://learn.microsoft.com/en-us/graph/api/outlookuser-supportedtimezones?view=graph-rest-beta&tabs=http
package graphBetaUsersUserMailboxSettings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserMailboxSettingsResourceModel represents the schema for the User Mailbox Settings resource
type UserMailboxSettingsResourceModel struct {
	ID                                    types.String   `tfsdk:"id"`
	UserID                                types.String   `tfsdk:"user_id"`
	AutomaticRepliesSetting               types.Object   `tfsdk:"automatic_replies_setting"` // Optional+Computed must be types.Object
	DateFormat                            types.String   `tfsdk:"date_format"`
	DelegateMeetingMessageDeliveryOptions types.String   `tfsdk:"delegate_meeting_message_delivery_options"`
	Language                              types.Object   `tfsdk:"language"` // Optional+Computed must be types.Object
	TimeFormat                            types.String   `tfsdk:"time_format"`
	TimeZone                              types.String   `tfsdk:"time_zone"`
	WorkingHours                          types.Object   `tfsdk:"working_hours"` // Optional+Computed must be types.Object
	UserPurpose                           types.String   `tfsdk:"user_purpose"`
	Timeouts                              timeouts.Value `tfsdk:"timeouts"`
}

// AutomaticRepliesSetting represents the automatic replies settings for a user's mailbox
type AutomaticRepliesSetting struct {
	Status                 types.String      `tfsdk:"status"`
	ExternalAudience       types.String      `tfsdk:"external_audience"`
	ScheduledStartDateTime *DateTimeTimeZone `tfsdk:"scheduled_start_date_time"`
	ScheduledEndDateTime   *DateTimeTimeZone `tfsdk:"scheduled_end_date_time"`
	InternalReplyMessage   types.String      `tfsdk:"internal_reply_message"`
	ExternalReplyMessage   types.String      `tfsdk:"external_reply_message"`
}

// DateTimeTimeZone represents a date and time with time zone information
type DateTimeTimeZone struct {
	DateTime types.String `tfsdk:"date_time"`
	TimeZone types.String `tfsdk:"time_zone"`
}

// LocaleInfo represents the locale (language and country/region) information
type LocaleInfo struct {
	Locale      types.String `tfsdk:"locale"`
	DisplayName types.String `tfsdk:"display_name"`
}

// WorkingHours represents the working hours for a user
type WorkingHours struct {
	DaysOfWeek types.Set    `tfsdk:"days_of_week"`
	StartTime  types.String `tfsdk:"start_time"`
	EndTime    types.String `tfsdk:"end_time"`
	TimeZone   types.Object `tfsdk:"time_zone"`
}

// TimeZoneBase represents time zone information (can be standard or custom)
type TimeZoneBase struct {
	Name types.String `tfsdk:"name"`
}
