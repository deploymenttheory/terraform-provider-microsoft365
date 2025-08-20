package graphBetaWindowsDeviceComplianceNotifications

import (
	"context"
	"fmt"
)

// validateRequest validates that exactly one localized message has is_default = true
func validateRequest(ctx context.Context, data *WindowsDeviceComplianceNotificationsResourceModel) error {
	// Validate that exactly one localized notification message has is_default = true
	if !data.LocalizedNotificationMessages.IsNull() && !data.LocalizedNotificationMessages.IsUnknown() {
		var localizedMessages []LocalizedNotificationMessageModel
		data.LocalizedNotificationMessages.ElementsAs(ctx, &localizedMessages, false)

		defaultCount := 0
		var defaultLocales []string

		for _, msg := range localizedMessages {
			if !msg.IsDefault.IsNull() && !msg.IsDefault.IsUnknown() && msg.IsDefault.ValueBool() {
				defaultCount++
				if !msg.Locale.IsNull() && !msg.Locale.IsUnknown() {
					defaultLocales = append(defaultLocales, msg.Locale.ValueString())
				}
			}
		}

		if defaultCount == 0 {
			return fmt.Errorf("at least one localized notification message must have 'is_default = true'. This message will be used as the default locale for the notification template")
		} else if defaultCount > 1 {
			return fmt.Errorf("only one localized notification message can have 'is_default = true'. Found %d messages with is_default = true for locales: %v. Please set is_default = true for only one message", defaultCount, defaultLocales)
		}
	}

	return nil
}
