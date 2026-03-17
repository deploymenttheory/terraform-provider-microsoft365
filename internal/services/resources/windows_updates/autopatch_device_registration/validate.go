package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateRequest performs validation before enrolling or unenrolling devices.
// It verifies that all Entra ID device IDs (Azure AD device object IDs) exist in Entra ID.
// Note: These are NOT Intune managed device IDs - they are Entra ID device object IDs.
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) validateRequest(
	ctx context.Context,
	data *WindowsUpdatesAutopatchDeviceRegistrationResourceModel,
	diagnostics *diag.Diagnostics,
) error {
	tflog.Debug(ctx, fmt.Sprintf("Starting Entra ID device validation for %s resource", ResourceName))

	if data.DeviceIds.IsNull() || data.DeviceIds.IsUnknown() {
		tflog.Debug(ctx, "No Entra ID device IDs to validate")
		return nil
	}

	elements := data.DeviceIds.Elements()
	if len(elements) == 0 {
		tflog.Debug(ctx, "Empty Entra ID device IDs set")
		return nil
	}

	var invalidEntraDeviceIDs []string
	var notFoundEntraDeviceIDs []string

	for _, elem := range elements {
		strVal, ok := elem.(types.String)
		if !ok {
			continue
		}

		entraDeviceID := strVal.ValueString()
		if entraDeviceID == "" {
			continue
		}

		err := r.validateEntraDeviceExists(ctx, entraDeviceID, diagnostics)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				notFoundEntraDeviceIDs = append(notFoundEntraDeviceIDs, entraDeviceID)
			} else {
				invalidEntraDeviceIDs = append(invalidEntraDeviceIDs, entraDeviceID)
			}
		}
	}

	if len(notFoundEntraDeviceIDs) > 0 {
		tflog.Error(ctx, "Some Entra ID device IDs do not exist in Entra ID", map[string]any{
			"notFoundEntraDeviceIDs": notFoundEntraDeviceIDs,
		})
		diagnostics.AddError(
			"Invalid Entra ID Device IDs",
			fmt.Sprintf("The following Entra ID device IDs (Azure AD device object IDs) do not exist in Entra ID: %s. "+
				"Please ensure all IDs are valid Entra ID device object IDs, not Intune managed device IDs.",
				strings.Join(notFoundEntraDeviceIDs, ", ")),
		)
		return fmt.Errorf("entra device validation failed: %d Entra ID devices not found", len(notFoundEntraDeviceIDs))
	}

	if len(invalidEntraDeviceIDs) > 0 {
		tflog.Error(ctx, "Failed to validate some Entra ID device IDs", map[string]any{
			"invalidEntraDeviceIDs": invalidEntraDeviceIDs,
		})
		diagnostics.AddError(
			"Entra ID Device Validation Failed",
			fmt.Sprintf("Failed to validate the following Entra ID device IDs: %s. "+
				"Please check the Entra ID device IDs and try again.",
				strings.Join(invalidEntraDeviceIDs, ", ")),
		)
		return fmt.Errorf("entra device validation failed: %d Entra ID devices could not be validated", len(invalidEntraDeviceIDs))
	}

	tflog.Debug(ctx, fmt.Sprintf("Entra ID device validation completed successfully for %d devices", len(elements)))
	return nil
}

// validateEntraDeviceExists validates that an Entra ID device (Azure AD device object) exists.
// This queries the /devices endpoint (Entra ID), NOT the /deviceManagement/managedDevices endpoint (Intune).
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) validateEntraDeviceExists(
	ctx context.Context,
	entraDeviceID string,
	diagnostics *diag.Diagnostics,
) error {
	tflog.Debug(ctx, "Validating Entra ID device exists", map[string]any{
		"entraDeviceId": entraDeviceID,
	})

	entraDevice, err := r.client.
		Devices().
		ByDeviceId(entraDeviceID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		tflog.Error(ctx, "Failed to retrieve Entra ID device", map[string]any{
			"entraDeviceId": entraDeviceID,
			"statusCode":    errorInfo.StatusCode,
			"errorCode":     errorInfo.ErrorCode,
		})

		if errorInfo.StatusCode == 404 {
			return fmt.Errorf("entra device %s not found", entraDeviceID)
		}

		return fmt.Errorf("failed to retrieve Entra ID device %s: %s", entraDeviceID, err.Error())
	}

	if entraDevice == nil {
		tflog.Error(ctx, "Entra ID device not found", map[string]any{
			"entraDeviceId": entraDeviceID,
		})
		return fmt.Errorf("entra device %s not found", entraDeviceID)
	}

	tflog.Debug(ctx, "Entra ID device validated successfully", map[string]any{
		"entraDeviceId": entraDeviceID,
		"displayName":   entraDevice.GetDisplayName(),
	})

	return nil
}
