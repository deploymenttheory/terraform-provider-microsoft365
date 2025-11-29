package graphBetaUsersUserMailboxSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource converts the Terraform resource model to the SDK model
func constructResource(ctx context.Context, data *UserMailboxSettingsResourceModel) (graphmodels.MailboxSettingsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	mailboxSettings := graphmodels.NewMailboxSettings()

	if data.AutomaticRepliesSetting != nil {
		automaticReplies := graphmodels.NewAutomaticRepliesSetting()

		if err := convert.FrameworkToGraphEnum(data.AutomaticRepliesSetting.Status,
			graphmodels.ParseAutomaticRepliesStatus,
			automaticReplies.SetStatus); err != nil {
			return nil, fmt.Errorf("failed to set automatic replies status: %w", err)
		}

		if err := convert.FrameworkToGraphEnum(data.AutomaticRepliesSetting.ExternalAudience,
			graphmodels.ParseExternalAudienceScope,
			automaticReplies.SetExternalAudience); err != nil {
			return nil, fmt.Errorf("failed to set external audience: %w", err)
		}

		convert.FrameworkToGraphString(data.AutomaticRepliesSetting.InternalReplyMessage, automaticReplies.SetInternalReplyMessage)
		convert.FrameworkToGraphString(data.AutomaticRepliesSetting.ExternalReplyMessage, automaticReplies.SetExternalReplyMessage)

		if data.AutomaticRepliesSetting.ScheduledStartDateTime != nil {
			scheduledStart := graphmodels.NewDateTimeTimeZone()
			convert.FrameworkToGraphString(data.AutomaticRepliesSetting.ScheduledStartDateTime.DateTime, scheduledStart.SetDateTime)
			convert.FrameworkToGraphString(data.AutomaticRepliesSetting.ScheduledStartDateTime.TimeZone, scheduledStart.SetTimeZone)
			automaticReplies.SetScheduledStartDateTime(scheduledStart)
		}

		if data.AutomaticRepliesSetting.ScheduledEndDateTime != nil {
			scheduledEnd := graphmodels.NewDateTimeTimeZone()
			convert.FrameworkToGraphString(data.AutomaticRepliesSetting.ScheduledEndDateTime.DateTime, scheduledEnd.SetDateTime)
			convert.FrameworkToGraphString(data.AutomaticRepliesSetting.ScheduledEndDateTime.TimeZone, scheduledEnd.SetTimeZone)
			automaticReplies.SetScheduledEndDateTime(scheduledEnd)
		}

		mailboxSettings.SetAutomaticRepliesSetting(automaticReplies)
	}

	convert.FrameworkToGraphString(data.DateFormat, mailboxSettings.SetDateFormat)

	if err := convert.FrameworkToGraphEnum(data.DelegateMeetingMessageDeliveryOptions,
		graphmodels.ParseDelegateMeetingMessageDeliveryOptions,
		mailboxSettings.SetDelegateMeetingMessageDeliveryOptions); err != nil {
		return nil, fmt.Errorf("failed to set delegate meeting message delivery options: %w", err)
	}

	convert.FrameworkToGraphString(data.TimeFormat, mailboxSettings.SetTimeFormat)
	convert.FrameworkToGraphString(data.TimeZone, mailboxSettings.SetTimeZone)

	// Handle language (types.Object)
	if !data.Language.IsNull() && !data.Language.IsUnknown() {
		var languageData LocaleInfo
		diags := data.Language.As(ctx, &languageData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract language data: %s", diags.Errors()[0].Detail())
		}

		language := graphmodels.NewLocaleInfo()
		convert.FrameworkToGraphString(languageData.Locale, language.SetLocale)
		mailboxSettings.SetLanguage(language)
	}

	if err := convert.FrameworkToGraphEnum(data.UserPurpose,
		graphmodels.ParseUserPurpose,
		mailboxSettings.SetUserPurpose); err != nil {
		return nil, fmt.Errorf("failed to set user purpose: %w", err)
	}

	// Handle working_hours (types.Object)
	if !data.WorkingHours.IsNull() && !data.WorkingHours.IsUnknown() {
		var workingHoursData WorkingHours
		diags := data.WorkingHours.As(ctx, &workingHoursData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract working hours data: %s", diags.Errors()[0].Detail())
		}

		workingHours := graphmodels.NewWorkingHours()

		// Convert days of week from string set to enum slice using helper
		if !workingHoursData.DaysOfWeek.IsNull() && !workingHoursData.DaysOfWeek.IsUnknown() {
			if err := convert.FrameworkToGraphObjectsFromStringSet(
				ctx,
				workingHoursData.DaysOfWeek,
				func(_ context.Context, values []string) []graphmodels.DayOfWeek {
					result := make([]graphmodels.DayOfWeek, 0, len(values))
					for _, val := range values {
						if dayEnum, err := graphmodels.ParseDayOfWeek(val); err == nil && dayEnum != nil {
							if enumValue, ok := dayEnum.(*graphmodels.DayOfWeek); ok && enumValue != nil {
								result = append(result, *enumValue)
							}
						}
					}
					return result
				},
				workingHours.SetDaysOfWeek,
			); err != nil {
				return nil, fmt.Errorf("failed to set working hours days of week: %w", err)
			}
		}

		// Use precision 0 for time values (HH:MM:SS format)
		if err := convert.FrameworkToGraphTimeOnlyWithPrecision(workingHoursData.StartTime, 0, workingHours.SetStartTime); err != nil {
			return nil, fmt.Errorf("failed to set working hours start time: %w", err)
		}

		if err := convert.FrameworkToGraphTimeOnlyWithPrecision(workingHoursData.EndTime, 0, workingHours.SetEndTime); err != nil {
			return nil, fmt.Errorf("failed to set working hours end time: %w", err)
		}

		// Handle time_zone (types.Object) within working_hours
		if !workingHoursData.TimeZone.IsNull() && !workingHoursData.TimeZone.IsUnknown() {
			var timeZoneData TimeZoneBase
			diags := workingHoursData.TimeZone.As(ctx, &timeZoneData, basetypes.ObjectAsOptions{})
			if !diags.HasError() {
				timeZone := graphmodels.NewTimeZoneBase()
				convert.FrameworkToGraphString(timeZoneData.Name, timeZone.SetName)
				workingHours.SetTimeZone(timeZone)
			}
		}

		mailboxSettings.SetWorkingHours(workingHours)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), mailboxSettings); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return mailboxSettings, nil
}
