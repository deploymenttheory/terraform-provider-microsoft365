package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

// verifyEnrollmentComplete polls the API to verify that devices have been successfully enrolled.
// This is necessary because enrollment has eventual consistency (typically 2-4 seconds).
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) verifyEnrollmentComplete(
	ctx context.Context,
	data *WindowsUpdatesAutopatchDeviceRegistrationResourceModel,
) error {
	updateCategory := data.UpdateCategory.ValueString()
	elements := data.EntraDeviceObjectIds.Elements()

	deviceIDs := make(map[string]bool)
	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok {
			deviceIDs[strVal.ValueString()] = true
		}
	}

	if len(deviceIDs) == 0 {
		tflog.Debug(ctx, "No devices to verify enrollment for")
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Verifying enrollment for %d devices", len(deviceIDs)))

	maxAttempts := 20
	waitInterval := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		tflog.Debug(ctx, fmt.Sprintf("Enrollment verification attempt %d/%d", attempt, maxAttempts))

		filter := "isof('microsoft.graph.windowsUpdates.azureADDevice')"
		result, err := r.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			Get(ctx, &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetRequestConfiguration{
				QueryParameters: &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetQueryParameters{
					Filter: &filter,
				},
			})

		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Error fetching devices for enrollment verification (attempt %d/%d): %v", attempt, maxAttempts, err))
			time.Sleep(waitInterval)
			continue
		}

		devices := result.GetValue()
		enrolledCount := 0

		for _, device := range devices {
			azureADDevice, ok := device.(windowsupdates.AzureADDeviceable)
			if !ok {
				continue
			}

			deviceID := azureADDevice.GetId()
			if deviceID == nil {
				continue
			}

			if !deviceIDs[*deviceID] {
				continue
			}

			enrollment := azureADDevice.GetEnrollment()
			if enrollment == nil {
				continue
			}

			var categoryEnrollment windowsupdates.UpdateCategoryEnrollmentInformationable

			switch updateCategory {
			case "feature":
				categoryEnrollment = enrollment.GetFeature()
			case "quality":
				categoryEnrollment = enrollment.GetQuality()
			case "driver":
				categoryEnrollment = enrollment.GetDriver()
			}

			if categoryEnrollment != nil {
				enrollmentState := categoryEnrollment.GetEnrollmentState()
				if enrollmentState != nil {
					stateStr := enrollmentState.String()
					if stateStr == "enrolled" || stateStr == "enrolledWithPolicy" {
						enrolledCount++
						tflog.Debug(ctx, fmt.Sprintf("Device %s is enrolled (state=%s)", *deviceID, stateStr))
					}
				}
			}
		}

		if enrolledCount == len(deviceIDs) {
			tflog.Debug(ctx, fmt.Sprintf("All %d devices successfully enrolled", enrolledCount))
			return nil
		}

		tflog.Debug(ctx, fmt.Sprintf("Enrollment verification: %d/%d devices enrolled, waiting %v", enrolledCount, len(deviceIDs), waitInterval))
		time.Sleep(waitInterval)
	}

	return fmt.Errorf("enrollment did not complete within expected time (40 seconds)")
}

// verifyUnenrollmentComplete polls the API to verify that devices have been successfully unenrolled.
// This is necessary because unenrollment has eventual consistency (typically 8-10 seconds).
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) verifyUnenrollmentComplete(
	ctx context.Context,
	data *WindowsUpdatesAutopatchDeviceRegistrationResourceModel,
) error {
	updateCategory := data.UpdateCategory.ValueString()
	elements := data.EntraDeviceObjectIds.Elements()

	deviceIDs := make(map[string]bool)
	for _, elem := range elements {
		if strVal, ok := elem.(types.String); ok {
			deviceIDs[strVal.ValueString()] = true
		}
	}

	if len(deviceIDs) == 0 {
		tflog.Debug(ctx, "No devices to verify unenrollment for")
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Verifying unenrollment for %d devices", len(deviceIDs)))

	maxAttempts := 25 // typical unenrollment time is 8-10 seconds, max 75 seconds
	waitInterval := 3 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		tflog.Debug(ctx, fmt.Sprintf("Unenrollment verification attempt %d/%d", attempt, maxAttempts))

		filter := "isof('microsoft.graph.windowsUpdates.azureADDevice')"
		result, err := r.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			Get(ctx, &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetRequestConfiguration{
				QueryParameters: &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetQueryParameters{
					Filter: &filter,
				},
			})

		if err != nil {
			tflog.Error(ctx, "Failed to query devices for unenrollment verification", map[string]any{
				"attempt": attempt,
				"error":   err.Error(),
			})
			return fmt.Errorf("failed to verify unenrollment: %w", err)
		}

		stillEnrolledCount := 0
		devices := result.GetValue()

		for _, asset := range devices {
			if asset == nil {
				continue
			}

			azureDevice, ok := asset.(windowsupdates.AzureADDeviceable)
			if !ok {
				continue
			}

			deviceID := azureDevice.GetId()
			if deviceID == nil || !deviceIDs[*deviceID] {
				continue
			}

			enrollment := azureDevice.GetEnrollment()
			if enrollment == nil {
				continue
			}

			var categoryEnrollment windowsupdates.UpdateCategoryEnrollmentInformationable
			switch updateCategory {
			case "driver":
				categoryEnrollment = enrollment.GetDriver()
			case "feature":
				categoryEnrollment = enrollment.GetFeature()
			case "quality":
				categoryEnrollment = enrollment.GetQuality()
			}

			if categoryEnrollment != nil {
				enrollmentState := categoryEnrollment.GetEnrollmentState()
				if enrollmentState != nil {
					stateStr := enrollmentState.String()
					if stateStr == "enrolled" || stateStr == "enrolledWithPolicy" {
						stillEnrolledCount++
						tflog.Debug(ctx, fmt.Sprintf("Device %s still enrolled (state: %s)", *deviceID, stateStr))
					}
				}
			}
		}

		if stillEnrolledCount == 0 {
			tflog.Debug(ctx, fmt.Sprintf("All devices successfully unenrolled after %d attempts", attempt))
			return nil
		}

		tflog.Debug(ctx, fmt.Sprintf("Still waiting for %d devices to unenroll", stillEnrolledCount))

		if attempt < maxAttempts {
			time.Sleep(waitInterval)
		}
	}

	return fmt.Errorf("unenrollment did not complete within expected time (75 seconds)")
}
