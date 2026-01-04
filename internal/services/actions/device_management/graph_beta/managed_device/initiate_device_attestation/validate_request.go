package graphBetaInitiateDeviceAttestationManagedDevice

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// DeviceAttestationValidationResult contains the results of device validation
type DeviceAttestationValidationResult struct {
	NonExistentManagedDevices   []string
	NonExistentComanagedDevices []string
	UnsupportedManagedDevices   []string
	UnsupportedComanagedDevices []string
}

// validateRequest performs API validation of devices
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, managedDeviceIDs []string, comanagedDeviceIDs []string) (*DeviceAttestationValidationResult, error) {
	result := &DeviceAttestationValidationResult{
		NonExistentManagedDevices:   make([]string, 0),
		NonExistentComanagedDevices: make([]string, 0),
		UnsupportedManagedDevices:   make([]string, 0),
		UnsupportedComanagedDevices: make([]string, 0),
	}

	// Validate managed devices
	for _, deviceID := range managedDeviceIDs {
		tflog.Debug(ctx, "Validating managed device", map[string]any{"device_id": deviceID})

		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentManagedDevices = append(result.NonExistentManagedDevices, deviceID)
				tflog.Warn(ctx, "Managed device not found", map[string]any{"device_id": deviceID})
				continue
			}
			return nil, fmt.Errorf("failed to validate managed device %s: %w", deviceID, err)
		}

		// Validate device is Windows (attestation requires Windows with TPM)
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices,
					fmt.Sprintf("%s (OS: %s, device attestation requires Windows with TPM)", deviceID, *device.GetOperatingSystem()))
				continue
			}
		} else {
			result.UnsupportedManagedDevices = append(result.UnsupportedManagedDevices,
				fmt.Sprintf("%s (Unknown OS)", deviceID))
			continue
		}

		// Note: We can't reliably check if TPM is available via the device object
		// The API will return an error if TPM is not available when we try to initiate attestation

		tflog.Debug(ctx, "Managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	// Validate co-managed devices - use managedDevices endpoint as comanagedDevices/{id} doesn't support GET
	for _, deviceID := range comanagedDeviceIDs {
		tflog.Debug(ctx, "Validating co-managed device", map[string]any{"device_id": deviceID})

		device, err := client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			graphErr := errors.GraphError(ctx, err)
			if graphErr.StatusCode == 404 {
				result.NonExistentComanagedDevices = append(result.NonExistentComanagedDevices, deviceID)
				tflog.Warn(ctx, "Co-managed device not found", map[string]any{"device_id": deviceID})
				continue
			}
			return nil, fmt.Errorf("failed to validate co-managed device %s: %w", deviceID, err)
		}

		// Validate device is Windows (attestation requires Windows with TPM)
		if device.GetOperatingSystem() != nil {
			os := strings.ToLower(*device.GetOperatingSystem())
			if !strings.Contains(os, "windows") {
				result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices,
					fmt.Sprintf("%s (OS: %s, device attestation requires Windows with TPM)", deviceID, *device.GetOperatingSystem()))
				continue
			}
		} else {
			result.UnsupportedComanagedDevices = append(result.UnsupportedComanagedDevices,
				fmt.Sprintf("%s (Unknown OS)", deviceID))
			continue
		}

		tflog.Debug(ctx, "Co-managed device validated successfully", map[string]any{"device_id": deviceID})
	}

	return result, nil
}
