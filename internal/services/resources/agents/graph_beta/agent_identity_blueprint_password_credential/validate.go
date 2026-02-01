package graphBetaAgentIdentityBlueprintPasswordCredential

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateRequest validates the request body during creation or update operations.
// Validates:
// - Start date is not in the past
// - End date is after start date
func validateRequest(ctx context.Context, data *AgentIdentityBlueprintPasswordCredentialResourceModel) error {
	tflog.Debug(ctx, "Starting validation of agent identity blueprint password credential request")

	if err := validateCredentialStartDateIsNotInThePast(ctx, data); err != nil {
		return err
	}

	if err := validateCredentialEndDateIsAfterStartDate(ctx, data); err != nil {
		return err
	}

	tflog.Debug(ctx, "Successfully validated agent identity blueprint password credential request")
	return nil
}

// validateCredentialStartDateIsNotInThePast validates that the start date is not in the past.
// Microsoft Graph API requires that credential start dates are either now or in the future.
func validateCredentialStartDateIsNotInThePast(ctx context.Context, data *AgentIdentityBlueprintPasswordCredentialResourceModel) error {

	if data.StartDateTime.IsNull() || data.StartDateTime.IsUnknown() {
		tflog.Debug(ctx, "No start_date_time provided, skipping past date validation")
		return nil
	}

	startDateStr := data.StartDateTime.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Validating start_date_time: %s", startDateStr))

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return fmt.Errorf("invalid start_date_time format '%s': %w. Expected ISO 8601 format (e.g., 2026-01-01T00:00:00Z)", startDateStr, err)
	}

	now := time.Now().UTC()

	// Allow a small buffer (5 seconds) to account for clock skew and processing time
	bufferDuration := -5 * time.Second
	allowedStartTime := now.Add(bufferDuration)

	if startDate.Before(allowedStartTime) {
		return fmt.Errorf("start_date_time '%s' is in the past. Microsoft Graph requires that credential start dates are either now or in the future. Current time (UTC): %s",
			startDateStr,
			now.Format(time.RFC3339))
	}

	tflog.Debug(ctx, fmt.Sprintf("start_date_time validation passed: %s is not in the past", startDateStr))
	return nil
}

// validateCredentialEndDateIsAfterStartDate validates that the end date is after the start date.
func validateCredentialEndDateIsAfterStartDate(ctx context.Context, data *AgentIdentityBlueprintPasswordCredentialResourceModel) error {

	if data.StartDateTime.IsNull() || data.StartDateTime.IsUnknown() ||
		data.EndDateTime.IsNull() || data.EndDateTime.IsUnknown() {
		tflog.Debug(ctx, "Start or end date not provided, skipping date order validation")
		return nil
	}

	startDateStr := data.StartDateTime.ValueString()
	endDateStr := data.EndDateTime.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Validating end_date_time (%s) is after start_date_time (%s)", endDateStr, startDateStr))

	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return fmt.Errorf("invalid start_date_time format '%s': %w", startDateStr, err)
	}

	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return fmt.Errorf("invalid end_date_time format '%s': %w", endDateStr, err)
	}

	if !endDate.After(startDate) {
		return fmt.Errorf("end_date_time '%s' must be after start_date_time '%s'", endDateStr, startDateStr)
	}

	tflog.Debug(ctx, "Date order validation passed: end_date_time is after start_date_time")
	return nil
}
