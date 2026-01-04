package graphBetaBypassActivationLockManagedDevice

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ValidationResult holds the results of device validation
type ValidationResult struct {
	NonExistentDevices            []string
	UnsupportedOSDevices          []string
	UnsupervisedIOSDevices        []string
	ActivationLockDisabledDevices []string
	AlreadyBypassedDevices        []string
}

// validateRequest performs API validation of devices before attempting bypass
// This function queries the Microsoft Graph API to validate device existence and compatibility
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, deviceIDs []string) (*ValidationResult, error) {
	tflog.Debug(ctx, fmt.Sprintf("Performing API validation for %d device(s)", len(deviceIDs)))

	result := &ValidationResult{
		NonExistentDevices:            []string{},
		UnsupportedOSDevices:          []string{},
		UnsupervisedIOSDevices:        []string{},
		ActivationLockDisabledDevices: []string{},
		AlreadyBypassedDevices:        []string{},
	}

	for _, deviceID := range deviceIDs {
		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)

			if errorInfo.StatusCode == 404 {
				result.NonExistentDevices = append(result.NonExistentDevices, deviceID)
				tflog.Debug(ctx, fmt.Sprintf("Device %s does not exist (404)", deviceID))
			} else {
				tflog.Error(ctx, fmt.Sprintf("Failed to validate device %s", deviceID), map[string]any{
					"device_id":   deviceID,
					"status_code": errorInfo.StatusCode,
					"error_code":  errorInfo.ErrorCode,
					"error":       errorInfo.ErrorMessage,
				})
				return nil, fmt.Errorf("failed to validate device %s: %s", deviceID, errorInfo.ErrorMessage)
			}
			continue
		}

		// Validate device type
		if device.GetDeviceType() != nil {
			deviceType := *device.GetDeviceType()

			// Activation Lock bypass is supported on iPhone, iPad, and Mac
			supportsActivationLockBypass := deviceType == models.IPHONE_DEVICETYPE ||
				deviceType == models.IPAD_DEVICETYPE ||
				deviceType == models.MAC_DEVICETYPE

			if !supportsActivationLockBypass {
				result.UnsupportedOSDevices = append(result.UnsupportedOSDevices,
					fmt.Sprintf("%s (deviceType: %s)", deviceID, deviceType.String()))
				tflog.Debug(ctx, fmt.Sprintf("Device %s has unsupported deviceType: %s", deviceID, deviceType.String()))
				continue
			}

			// iOS/iPadOS must be supervised
			if (deviceType == models.IPHONE_DEVICETYPE || deviceType == models.IPAD_DEVICETYPE) &&
				(device.GetIsSupervised() == nil || !*device.GetIsSupervised()) {
				result.UnsupervisedIOSDevices = append(result.UnsupervisedIOSDevices,
					fmt.Sprintf("%s (deviceType: %s - not supervised)", deviceID, deviceType.String()))
				tflog.Debug(ctx, fmt.Sprintf("Device %s is iOS/iPadOS but not supervised", deviceID))
				continue
			}

			// Check if device already has bypass code
			if device.GetActivationLockBypassCode() != nil && *device.GetActivationLockBypassCode() != "" {
				result.AlreadyBypassedDevices = append(result.AlreadyBypassedDevices, deviceID)
				tflog.Debug(ctx, fmt.Sprintf("Device %s already has an Activation Lock bypass code", deviceID))
			}

			// Check if Activation Lock may not be enabled (macOS specific)
			if deviceType == models.MAC_DEVICETYPE {
				// For macOS, if not clearly DEP enrolled, note potential issue
				result.ActivationLockDisabledDevices = append(result.ActivationLockDisabledDevices,
					fmt.Sprintf("%s (may not have Activation Lock enabled)", deviceID))
				tflog.Debug(ctx, fmt.Sprintf("Device %s is macOS - Activation Lock status unclear", deviceID))
			}
		} else {
			result.UnsupportedOSDevices = append(result.UnsupportedOSDevices, fmt.Sprintf("%s (Unknown deviceType)", deviceID))
			tflog.Debug(ctx, fmt.Sprintf("Device %s has unknown deviceType", deviceID))
			continue
		}
	}

	tflog.Debug(ctx, "API validation completed", map[string]any{
		"non_existent":            len(result.NonExistentDevices),
		"unsupported_os":          len(result.UnsupportedOSDevices),
		"unsupervised_ios":        len(result.UnsupervisedIOSDevices),
		"already_bypassed":        len(result.AlreadyBypassedDevices),
		"activation_lock_unclear": len(result.ActivationLockDisabledDevices),
	})

	return result, nil
}
