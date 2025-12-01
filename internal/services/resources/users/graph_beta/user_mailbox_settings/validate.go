package graphBetaUsersUserMailboxSettings

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates the mailbox settings request against supported values from the API
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, userID string, data *UserMailboxSettingsResourceModel) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating mailbox settings request for user: %s", userID))

	// Validate language/locale if provided
	if !data.Language.IsNull() && !data.Language.IsUnknown() {
		var languageData LocaleInfo
		diags := data.Language.As(ctx, &languageData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fmt.Errorf("failed to extract language data for validation: %s", diags.Errors()[0].Detail())
		}

		if !languageData.Locale.IsNull() && !languageData.Locale.IsUnknown() {
			locale := languageData.Locale.ValueString()
			if err := validateSupportedLanguage(ctx, client, userID, locale); err != nil {
				return err
			}
		}
	}

	// Validate timezone if provided
	if !data.TimeZone.IsNull() && !data.TimeZone.IsUnknown() {
		timezone := data.TimeZone.ValueString()
		if err := validateSupportedTimeZones(ctx, client, userID, timezone); err != nil {
			return err
		}
	}

	// Validate working hours timezone if provided
	if !data.WorkingHours.IsNull() && !data.WorkingHours.IsUnknown() {
		var workingHoursData WorkingHours
		diags := data.WorkingHours.As(ctx, &workingHoursData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return fmt.Errorf("failed to extract working hours data for validation: %s", diags.Errors()[0].Detail())
		}

		if !workingHoursData.TimeZone.IsNull() && !workingHoursData.TimeZone.IsUnknown() {
			var timeZoneData TimeZoneBase
			diags := workingHoursData.TimeZone.As(ctx, &timeZoneData, basetypes.ObjectAsOptions{})
			if !diags.HasError() && !timeZoneData.Name.IsNull() && !timeZoneData.Name.IsUnknown() {
				timezone := timeZoneData.Name.ValueString()
				if err := validateSupportedTimeZones(ctx, client, userID, timezone); err != nil {
					return fmt.Errorf("working hours timezone validation failed: %w", err)
				}
			}
		}
	}

	tflog.Debug(ctx, "Mailbox settings validation completed successfully")
	return nil
}

// validateSupportedLanguage validates that the provided locale is supported by the user's mailbox
func validateSupportedLanguage(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, userID string, locale string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating locale '%s' against supported languages for user %s", locale, userID))

	// Call the Microsoft Graph API to get supported languages
	supportedLanguages, err := client.
		Users().
		ByUserId(userID).
		Outlook().
		SupportedLanguages().
		Get(ctx, nil)

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to retrieve supported languages for validation: %v", err))
		// Don't fail the operation if we can't retrieve the list - the API will validate it
		return nil
	}

	if supportedLanguages == nil {
		tflog.Debug(ctx, "No supported languages list returned, skipping validation")
		return nil
	}

	// Debug log the full API response as JSON
	if err := constructors.DebugLogGraphObject(ctx, "Full supportedLanguages API response", supportedLanguages); err != nil {
		tflog.Error(ctx, "Failed to debug log supportedLanguages response", map[string]any{"error": err.Error()})
	}

	// Log the full response for debugging
	languages := supportedLanguages.GetValue()
	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d supported languages from API", len(languages)))

	// Build and log the list of supported locales
	var supportedLocales []string
	for _, lang := range languages {
		if lang.GetLocale() != nil {
			locale := *lang.GetLocale()
			displayName := ""
			if lang.GetDisplayName() != nil {
				displayName = *lang.GetDisplayName()
			}
			supportedLocales = append(supportedLocales, fmt.Sprintf("%s (%s)", locale, displayName))
			tflog.Trace(ctx, fmt.Sprintf("Supported locale: %s - %s", locale, displayName))
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("Full list of supported locales: %v", supportedLocales))

	// Check if the provided locale is in the supported list
	for _, lang := range languages {
		if lang.GetLocale() != nil && *lang.GetLocale() == locale {
			tflog.Debug(ctx, fmt.Sprintf("✓ Locale '%s' is supported", locale))
			return nil
		}
	}

	// Locale not found in supported list
	tflog.Warn(ctx, fmt.Sprintf("✗ Locale '%s' not found in supported languages list", locale))
	return fmt.Errorf("locale '%s' is not supported by the user's mailbox. To get the list of supported locales, query GET /users/%s/outlook/supportedLanguages", locale, userID)
}

// validateSupportedTimeZones validates that the provided timezone is supported by the user's mailbox
func validateSupportedTimeZones(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, userID string, timezone string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating timezone '%s' against supported timezones for user %s", timezone, userID))

	// Call the Microsoft Graph API to get supported time zones
	// Try both Windows and IANA formats
	supportedTimeZones, err := client.
		Users().
		ByUserId(userID).
		Outlook().
		SupportedTimeZones().
		Get(ctx, nil)

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to retrieve supported time zones for validation: %v", err))
		// Don't fail the operation if we can't retrieve the list - the API will validate it
		return nil
	}

	if supportedTimeZones == nil {
		tflog.Debug(ctx, "No supported time zones list returned, skipping validation")
		return nil
	}

	// Debug log the full API response as JSON
	if err := constructors.DebugLogGraphObject(ctx, "Full supportedTimeZones API response", supportedTimeZones); err != nil {
		tflog.Error(ctx, "Failed to debug log supportedTimeZones response", map[string]any{"error": err.Error()})
	}

	// Log the full response for debugging
	timeZones := supportedTimeZones.GetValue()
	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d supported time zones from API", len(timeZones)))

	// Build and log the list of supported time zones
	var supportedTZList []string
	for _, tz := range timeZones {
		alias := ""
		displayName := ""
		if tz.GetAlias() != nil {
			alias = *tz.GetAlias()
		}
		if tz.GetDisplayName() != nil {
			displayName = *tz.GetDisplayName()
		}

		if alias != "" || displayName != "" {
			supportedTZList = append(supportedTZList, fmt.Sprintf("%s (%s)", alias, displayName))
			tflog.Trace(ctx, fmt.Sprintf("Supported timezone: alias=%s, displayName=%s", alias, displayName))
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("Full list of supported time zones (first 20): %v", limitSlice(supportedTZList, 20)))

	// Check if the provided timezone is in the supported list
	for _, tz := range timeZones {
		if tz.GetAlias() != nil && *tz.GetAlias() == timezone {
			tflog.Debug(ctx, fmt.Sprintf("✓ Timezone '%s' is supported (matched alias)", timezone))
			return nil
		}
		if tz.GetDisplayName() != nil && *tz.GetDisplayName() == timezone {
			tflog.Debug(ctx, fmt.Sprintf("✓ Timezone '%s' is supported (matched displayName)", timezone))
			return nil
		}
	}

	// Timezone not found in supported list
	tflog.Warn(ctx, fmt.Sprintf("✗ Timezone '%s' not found in supported time zones list", timezone))
	return fmt.Errorf("timezone '%s' is not supported by the user's mailbox. To get the list of supported time zones, query GET /users/%s/outlook/supportedTimeZones", timezone, userID)
}

// limitSlice returns the first n elements of a slice, or the entire slice if it's smaller than n
func limitSlice(slice []string, n int) []string {
	if len(slice) <= n {
		return slice
	}
	return slice[:n]
}
