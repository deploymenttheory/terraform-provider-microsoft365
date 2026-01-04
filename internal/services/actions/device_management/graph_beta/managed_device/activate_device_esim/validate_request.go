package graphBetaActivateDeviceEsimManagedDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// EsimValidationResult holds the results of device validation for eSIM activation
type EsimValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	UnsupportedManagedDevices   []string
	UnsupportedComanagedDevices []string
}

// validateRequest performs API validation of devices before attempting eSIM activation
// This function queries the Microsoft Graph API to validate device existence and compatibility
// GET https://graph.microsoft.com/beta/deviceManagement/managedDevices/{managedDeviceId}
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDevices []ManagedDeviceActivateEsim, comanagedDevices []ComanagedDeviceActivateEsim) (*EsimValidationResult, error) {
	tflog.Debug(ctx, fmt.Sprintf("Performing API validation for %d managed and %d co-managed device(s)",
		len(managedDevices), len(comanagedDevices)))

	result := &EsimValidationResult{
		NonExistentManagedDevices:   []string{},
		NonExistentComanagedDevices: []string{},
		UnsupportedManagedDevices:   []string{},
		UnsupportedComanagedDevices: []string{},
	}

	// Validate managed devices
	for _, managedDevice := range managedDevices {
		deviceID := managedDevice.DeviceID.ValueString()
		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 {
				result.NonExistentManagedDevices = append(result.NonExistentManagedDevices, deviceID)
				tflog.Debug(ctx, fmt.Sprintf("Managed device %s does not exist (404)", deviceID))
			} else {
				tflog.Error(ctx, fmt.Sprintf("Failed to validate managed device %s", deviceID), map[string]any{
					"device_id":   deviceID,
					"status_code": errorInfo.StatusCode,
					"error_code":  errorInfo.ErrorCode,
					"error":       errorInfo.ErrorMessage,
				})
				return nil, fmt.Errorf("failed to validate managed device %s: %s", deviceID, errorInfo.ErrorMessage)
			}
			continue
		}

		// Validate device type supports eSIM activation (iOS/iPadOS only)
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			supportsESIM := deviceType == models.IPHONE_DEVICETYPE ||
				deviceType == models.IPAD_DEVICETYPE

			if !supportsESIM {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
				tflog.Debug(ctx, fmt.Sprintf("Managed device %s has unsupported deviceType: %s", deviceID, deviceType.String()))
			} else {
				tflog.Debug(ctx, fmt.Sprintf("Managed device %s (deviceType: %s) supports eSIM", deviceID, deviceType.String()))
			}
		} else {
			result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
			tflog.Debug(ctx, fmt.Sprintf("Managed device %s has unknown deviceType", deviceID))
		}
	}

	// Validate co-managed devices using managedDevices endpoint
	// Note: Co-managed devices are still managed devices in Intune, they're just also managed by another system
	// The comanagedDevices endpoint doesn't support GET by ID, so we use managedDevices for validation
	for _, comanagedDevice := range comanagedDevices {
		deviceID := comanagedDevice.DeviceID.ValueString()
		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 {
				result.NonExistentComanagedDevices = append(result.NonExistentComanagedDevices, deviceID)
				tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s does not exist (404)", deviceID))
			} else {
				tflog.Error(ctx, fmt.Sprintf("Failed to validate co-managed device %s", deviceID), map[string]any{
					"device_id":   deviceID,
					"status_code": errorInfo.StatusCode,
					"error_code":  errorInfo.ErrorCode,
					"error":       errorInfo.ErrorMessage,
				})
				return nil, fmt.Errorf("failed to validate co-managed device %s: %s", deviceID, errorInfo.ErrorMessage)
			}
			continue
		}

		// Validate device type supports eSIM activation (iOS/iPadOS only)
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			supportsESIM := deviceType == models.IPHONE_DEVICETYPE ||
				deviceType == models.IPAD_DEVICETYPE

			if !supportsESIM {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
				tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s has unsupported deviceType: %s", deviceID, deviceType.String()))
			} else {
				tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s (deviceType: %s) supports eSIM", deviceID, deviceType.String()))
			}
		} else {
			result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
			tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s has unknown deviceType", deviceID))
		}
	}

	tflog.Debug(ctx, "API validation completed", map[string]any{
		"non_existent_managed":   len(result.NonExistentManagedDevices),
		"non_existent_comanaged": len(result.NonExistentComanagedDevices),
		"unsupported_managed":    len(result.UnsupportedManagedDevices),
		"unsupported_comanaged":  len(result.UnsupportedComanagedDevices),
	})

	return result, nil
}
