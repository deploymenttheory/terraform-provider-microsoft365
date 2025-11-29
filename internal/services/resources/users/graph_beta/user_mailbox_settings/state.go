package graphBetaUsersUserMailboxSettings

import (
	"context"

	commonattr "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/attr"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote mailbox settings to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *UserMailboxSettingsResourceModel, remoteResource graphmodels.MailboxSettingsable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote mailbox settings to Terraform state")

	data.DateFormat = convert.GraphToFrameworkString(remoteResource.GetDateFormat())
	data.DelegateMeetingMessageDeliveryOptions = convert.GraphToFrameworkEnum(remoteResource.GetDelegateMeetingMessageDeliveryOptions())
	data.TimeFormat = convert.GraphToFrameworkString(remoteResource.GetTimeFormat())
	data.TimeZone = convert.GraphToFrameworkString(remoteResource.GetTimeZone())
	data.UserPurpose = convert.GraphToFrameworkEnum(remoteResource.GetUserPurpose())

	if automaticReplies := remoteResource.GetAutomaticRepliesSetting(); automaticReplies != nil {
		tflog.Debug(ctx, "Mapping automaticRepliesSetting")
		data.AutomaticRepliesSetting = mapAutomaticRepliesSetting(ctx, automaticReplies)
	} else {
		tflog.Debug(ctx, "automaticRepliesSetting not found")
		data.AutomaticRepliesSetting = nil
	}

	if language := remoteResource.GetLanguage(); language != nil {
		tflog.Debug(ctx, "Mapping language")
		data.Language = mapLocaleInfoToObject(ctx, language)
	} else {
		tflog.Debug(ctx, "language not found")
		data.Language = types.ObjectNull(map[string]attr.Type{
			"locale":       types.StringType,
			"display_name": types.StringType,
		})
	}

	if workingHours := remoteResource.GetWorkingHours(); workingHours != nil {
		tflog.Debug(ctx, "Mapping workingHours")
		data.WorkingHours = mapWorkingHoursToObject(ctx, workingHours)
	} else {
		tflog.Debug(ctx, "workingHours not found")
		data.WorkingHours = types.ObjectNull(map[string]attr.Type{
			"days_of_week": types.SetType{ElemType: types.StringType},
			"start_time":   types.StringType,
			"end_time":     types.StringType,
			"time_zone":    types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType}},
		})
	}

	tflog.Debug(ctx, "Completed mapping remote mailbox settings to Terraform state")
}

// mapAutomaticRepliesSetting maps the automatic replies setting from API response
func mapAutomaticRepliesSetting(ctx context.Context, automaticReplies graphmodels.AutomaticRepliesSettingable) *AutomaticRepliesSetting {
	if automaticReplies == nil {
		return nil
	}

	result := &AutomaticRepliesSetting{}
	result.Status = convert.GraphToFrameworkEnum(automaticReplies.GetStatus())
	result.ExternalAudience = convert.GraphToFrameworkEnum(automaticReplies.GetExternalAudience())
	result.InternalReplyMessage = convert.GraphToFrameworkString(automaticReplies.GetInternalReplyMessage())
	result.ExternalReplyMessage = convert.GraphToFrameworkString(automaticReplies.GetExternalReplyMessage())

	if scheduledStart := automaticReplies.GetScheduledStartDateTime(); scheduledStart != nil {
		result.ScheduledStartDateTime = mapDateTimeTimeZone(ctx, scheduledStart)
	}

	if scheduledEnd := automaticReplies.GetScheduledEndDateTime(); scheduledEnd != nil {
		result.ScheduledEndDateTime = mapDateTimeTimeZone(ctx, scheduledEnd)
	}

	return result
}

// mapDateTimeTimeZone maps a date time with time zone from API response
func mapDateTimeTimeZone(ctx context.Context, dateTime graphmodels.DateTimeTimeZoneable) *DateTimeTimeZone {
	if dateTime == nil {
		return nil
	}

	result := &DateTimeTimeZone{}
	result.DateTime = convert.GraphToFrameworkString(dateTime.GetDateTime())
	result.TimeZone = convert.GraphToFrameworkString(dateTime.GetTimeZone())
	return result
}

// mapLocaleInfo maps the locale information from API response (for nested use)
func mapLocaleInfo(ctx context.Context, locale graphmodels.LocaleInfoable) *LocaleInfo {
	if locale == nil {
		return nil
	}

	result := &LocaleInfo{}
	result.Locale = convert.GraphToFrameworkString(locale.GetLocale())
	result.DisplayName = convert.GraphToFrameworkString(locale.GetDisplayName())
	return result
}

// mapLocaleInfoToObject maps the locale information to types.Object
func mapLocaleInfoToObject(ctx context.Context, locale graphmodels.LocaleInfoable) types.Object {
	if locale == nil {
		return types.ObjectNull(map[string]attr.Type{
			"locale":       types.StringType,
			"display_name": types.StringType,
		})
	}

	attrTypes := map[string]attr.Type{
		"locale":       types.StringType,
		"display_name": types.StringType,
	}

	attrValues := map[string]attr.Value{
		"locale":       convert.GraphToFrameworkString(locale.GetLocale()),
		"display_name": convert.GraphToFrameworkString(locale.GetDisplayName()),
	}

	return commonattr.ObjectValue(attrTypes, attrValues)
}

// mapWorkingHours maps the working hours from API response (for nested use)
func mapWorkingHours(ctx context.Context, workingHours graphmodels.WorkingHoursable) *WorkingHours {
	if workingHours == nil {
		return nil
	}

	result := &WorkingHours{}
	result.StartTime = convert.GraphToFrameworkTimeOnly(workingHours.GetStartTime())
	result.EndTime = convert.GraphToFrameworkTimeOnly(workingHours.GetEndTime())

	daysOfWeek := workingHours.GetDaysOfWeek()
	if len(daysOfWeek) > 0 {
		dayStrings := convert.GraphToFrameworkEnumSlice(daysOfWeek)
		setValue, diags := types.SetValueFrom(ctx, types.StringType, dayStrings)
		if diags.HasError() {
			tflog.Error(ctx, "Error creating set for daysOfWeek", map[string]any{"diags": diags})
			result.DaysOfWeek = types.SetNull(types.StringType)
		} else {
			result.DaysOfWeek = setValue
		}
	} else {
		emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		result.DaysOfWeek = emptySet
	}

	// Handle time_zone as types.Object
	if timeZone := workingHours.GetTimeZone(); timeZone != nil {
		result.TimeZone = commonattr.ObjectValue(
			map[string]attr.Type{"name": types.StringType},
			map[string]attr.Value{"name": convert.GraphToFrameworkString(timeZone.GetName())},
		)
	} else {
		result.TimeZone = types.ObjectNull(map[string]attr.Type{"name": types.StringType})
	}

	return result
}

// mapWorkingHoursToObject maps the working hours to types.Object
func mapWorkingHoursToObject(ctx context.Context, workingHours graphmodels.WorkingHoursable) types.Object {
	if workingHours == nil {
		return types.ObjectNull(map[string]attr.Type{
			"days_of_week": types.SetType{ElemType: types.StringType},
			"start_time":   types.StringType,
			"end_time":     types.StringType,
			"time_zone":    types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType}},
		})
	}

	attrTypes := map[string]attr.Type{
		"days_of_week": types.SetType{ElemType: types.StringType},
		"start_time":   types.StringType,
		"end_time":     types.StringType,
		"time_zone":    types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType}},
	}

	// Build days_of_week set
	daysOfWeek := workingHours.GetDaysOfWeek()
	var daysSet types.Set
	if len(daysOfWeek) > 0 {
		dayStrings := convert.GraphToFrameworkEnumSlice(daysOfWeek)
		setValue, diags := types.SetValueFrom(ctx, types.StringType, dayStrings)
		if diags.HasError() {
			tflog.Error(ctx, "Error creating set for daysOfWeek", map[string]any{"diags": diags})
			daysSet = types.SetNull(types.StringType)
		} else {
			daysSet = setValue
		}
	} else {
		emptySet, _ := types.SetValueFrom(ctx, types.StringType, []string{})
		daysSet = emptySet
	}

	// Build time_zone object
	var timeZoneObj types.Object
	if timeZone := workingHours.GetTimeZone(); timeZone != nil {
		timeZoneObj = commonattr.ObjectValue(
			map[string]attr.Type{"name": types.StringType},
			map[string]attr.Value{"name": convert.GraphToFrameworkString(timeZone.GetName())},
		)
	} else {
		timeZoneObj = types.ObjectNull(map[string]attr.Type{"name": types.StringType})
	}

	attrValues := map[string]attr.Value{
		"days_of_week": daysSet,
		"start_time":   convert.GraphToFrameworkTimeOnly(workingHours.GetStartTime()),
		"end_time":     convert.GraphToFrameworkTimeOnly(workingHours.GetEndTime()),
		"time_zone":    timeZoneObj,
	}

	return commonattr.ObjectValue(attrTypes, attrValues)
}

// mapTimeZoneBase maps the time zone base from API response
func mapTimeZoneBase(ctx context.Context, timeZone graphmodels.TimeZoneBaseable) *TimeZoneBase {
	if timeZone == nil {
		return nil
	}

	result := &TimeZoneBase{}
	result.Name = convert.GraphToFrameworkString(timeZone.GetName())
	return result
}
