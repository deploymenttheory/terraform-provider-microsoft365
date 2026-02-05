package graphBetaWindowsFeatureUpdatePolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateRequest performs all validation checks on the resource before sending to API
func validateRequest(ctx context.Context, data *WindowsFeatureUpdatePolicyResourceModel) error {
	tflog.Debug(ctx, "Starting request validation", map[string]any{
		"resource": ResourceName,
	})

	if err := validateOfferStartDateTime(ctx, data); err != nil {
		return err
	}

	if err := validateOfferEndDateTime(ctx, data); err != nil {
		return err
	}

	if err := validateOfferIntervalInDays(ctx, data); err != nil {
		return err
	}

	tflog.Debug(ctx, "Request validation completed successfully")
	return nil
}

// validateOfferStartDateTime ensures offer_start_date_time_in_utc is today or in the future
func validateOfferStartDateTime(ctx context.Context, data *WindowsFeatureUpdatePolicyResourceModel) error {

	if data.RolloutSettings.IsNull() || data.RolloutSettings.IsUnknown() {
		tflog.Debug(ctx, "Skipping offer start date validation - rollout_settings is null or unknown")
		return nil
	}

	var rolloutModel RolloutSettingsModel
	if diags := data.RolloutSettings.As(ctx, &rolloutModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return fmt.Errorf("failed to parse rollout_settings: %v", diags.Errors())
	}

	if rolloutModel.OfferStartDateTimeInUTC.IsNull() || rolloutModel.OfferStartDateTimeInUTC.IsUnknown() {
		tflog.Debug(ctx, "Skipping offer start date validation - offer_start_date_time_in_utc is null or unknown")
		return nil
	}

	offerStartDateStr := rolloutModel.OfferStartDateTimeInUTC.ValueString()

	offerStartDate, err := time.Parse(time.RFC3339, offerStartDateStr)
	if err != nil {
		return fmt.Errorf("invalid offer_start_date_time_in_utc format '%s': must be in RFC3339 format (e.g., '2026-05-02T00:00:00Z'): %w", offerStartDateStr, err)
	}

	// Get current time in UTC and truncate to start of day for comparison
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	offerStartDay := time.Date(offerStartDate.Year(), offerStartDate.Month(), offerStartDate.Day(), 0, 0, 0, 0, time.UTC)

	tflog.Debug(ctx, "Validating offer start date", map[string]any{
		"offerStartDate": offerStartDateStr,
		"today":          today.Format(time.RFC3339),
		"offerStartDay":  offerStartDay.Format(time.RFC3339),
	})

	if offerStartDay.Before(today) {
		return fmt.Errorf(
			"offer_start_date_time_in_utc cannot be in the past: provided date '%s' is before today '%s'. "+
				"The offer start date must be today or a future date",
			offerStartDate.Format("2006-01-02"),
			today.Format("2006-01-02"),
		)
	}

	tflog.Debug(ctx, "Offer start date validation passed", map[string]any{
		"offerStartDate": offerStartDateStr,
	})

	return nil
}

// validateOfferEndDateTime ensures offer_end_date_time_in_utc is after offer_start_date_time_in_utc
func validateOfferEndDateTime(ctx context.Context, data *WindowsFeatureUpdatePolicyResourceModel) error {

	if data.RolloutSettings.IsNull() || data.RolloutSettings.IsUnknown() {
		tflog.Debug(ctx, "Skipping offer end date validation - rollout_settings is null or unknown")
		return nil
	}

	var rolloutModel RolloutSettingsModel
	if diags := data.RolloutSettings.As(ctx, &rolloutModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return fmt.Errorf("failed to parse rollout_settings: %v", diags.Errors())
	}

	// Skip validation if either start or end date is null or unknown
	if rolloutModel.OfferStartDateTimeInUTC.IsNull() || rolloutModel.OfferStartDateTimeInUTC.IsUnknown() ||
		rolloutModel.OfferEndDateTimeInUTC.IsNull() || rolloutModel.OfferEndDateTimeInUTC.IsUnknown() {
		tflog.Debug(ctx, "Skipping offer end date validation - start or end date is null or unknown")
		return nil
	}

	startDate, err := time.Parse(time.RFC3339, rolloutModel.OfferStartDateTimeInUTC.ValueString())
	if err != nil {
		return fmt.Errorf("invalid offer_start_date_time_in_utc format: %w", err)
	}

	endDate, err := time.Parse(time.RFC3339, rolloutModel.OfferEndDateTimeInUTC.ValueString())
	if err != nil {
		return fmt.Errorf("invalid offer_end_date_time_in_utc format: %w", err)
	}

	if endDate.Before(startDate) {
		return fmt.Errorf(
			"offer_end_date_time_in_utc ('%s') must be on or after offer_start_date_time_in_utc ('%s')",
			rolloutModel.OfferEndDateTimeInUTC.ValueString(),
			rolloutModel.OfferStartDateTimeInUTC.ValueString(),
		)
	}

	tflog.Debug(ctx, "Rollout date range validation passed", map[string]any{
		"startDate": rolloutModel.OfferStartDateTimeInUTC.ValueString(),
		"endDate":   rolloutModel.OfferEndDateTimeInUTC.ValueString(),
	})

	return nil
}

// validateOfferIntervalInDays ensures offer_interval_in_days is not greater than the days between start and end dates
func validateOfferIntervalInDays(ctx context.Context, data *WindowsFeatureUpdatePolicyResourceModel) error {

	if data.RolloutSettings.IsNull() || data.RolloutSettings.IsUnknown() {
		tflog.Debug(ctx, "Skipping offer interval validation - rollout_settings is null or unknown")
		return nil
	}

	var rolloutModel RolloutSettingsModel
	if diags := data.RolloutSettings.As(ctx, &rolloutModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return fmt.Errorf("failed to parse rollout_settings: %v", diags.Errors())
	}

	if rolloutModel.OfferIntervalInDays.IsNull() || rolloutModel.OfferIntervalInDays.IsUnknown() ||
		rolloutModel.OfferStartDateTimeInUTC.IsNull() || rolloutModel.OfferStartDateTimeInUTC.IsUnknown() ||
		rolloutModel.OfferEndDateTimeInUTC.IsNull() || rolloutModel.OfferEndDateTimeInUTC.IsUnknown() {
		tflog.Debug(ctx, "Skipping offer interval validation - one or more required fields is null or unknown")
		return nil
	}

	startDate, err := time.Parse(time.RFC3339, rolloutModel.OfferStartDateTimeInUTC.ValueString())
	if err != nil {
		return fmt.Errorf("invalid offer_start_date_time_in_utc format: %w", err)
	}

	endDate, err := time.Parse(time.RFC3339, rolloutModel.OfferEndDateTimeInUTC.ValueString())
	if err != nil {
		return fmt.Errorf("invalid offer_end_date_time_in_utc format: %w", err)
	}

	daysDifference := int32(endDate.Sub(startDate).Hours() / 24)
	offerInterval := rolloutModel.OfferIntervalInDays.ValueInt32()

	tflog.Debug(ctx, "Validating offer interval", map[string]any{
		"startDate":      rolloutModel.OfferStartDateTimeInUTC.ValueString(),
		"endDate":        rolloutModel.OfferEndDateTimeInUTC.ValueString(),
		"daysDifference": daysDifference,
		"offerInterval":  offerInterval,
	})

	if offerInterval > daysDifference {
		return fmt.Errorf(
			"offer_interval_in_days (%d) cannot be greater than the number of days between "+
				"offer_start_date_time_in_utc ('%s') and offer_end_date_time_in_utc ('%s') which is %d days",
			offerInterval,
			startDate.Format("01/02/2006"),
			endDate.Format("01/02/2006"),
			daysDifference,
		)
	}

	tflog.Debug(ctx, "Offer interval validation passed", map[string]any{
		"offerInterval":  offerInterval,
		"daysDifference": daysDifference,
	})

	return nil
}
