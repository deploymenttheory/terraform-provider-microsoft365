package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"
	"strings"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
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

	if data.EntraDeviceObjectIds.IsNull() || data.EntraDeviceObjectIds.IsUnknown() {
		tflog.Debug(ctx, "No Entra ID device IDs to validate")
		return nil
	}

	elements := data.EntraDeviceObjectIds.Elements()
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

		err := r.validateEntraDeviceExists(ctx, entraDeviceID)
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
		return fmt.Errorf("%w: %d Entra ID devices not found", sentinels.ErrEntraDeviceValidationFailed, len(notFoundEntraDeviceIDs))
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
		return fmt.Errorf("%w: %d Entra ID devices could not be validated", sentinels.ErrEntraDeviceValidationFailed, len(invalidEntraDeviceIDs))
	}

	tflog.Debug(ctx, fmt.Sprintf("Entra ID device validation completed successfully for %d devices", len(elements)))
	return nil
}

// validateEntraDeviceExists validates that an Entra ID device (Azure AD device object) exists.
// This queries the /devices endpoint (Entra ID), NOT the /deviceManagement/managedDevices endpoint (Intune).
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) validateEntraDeviceExists(
	ctx context.Context,
	entraDeviceObjectID string,
) error {
	tflog.Debug(ctx, "Validating Entra ID device exists", map[string]any{
		"entraDeviceObjectId": entraDeviceObjectID,
	})

	entraDevice, err := r.client.
		Devices().
		ByDeviceId(entraDeviceObjectID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		tflog.Error(ctx, "Failed to retrieve Entra ID device", map[string]any{
			"entraDeviceObjectId": entraDeviceObjectID,
			"statusCode":          errorInfo.StatusCode,
			"errorCode":           errorInfo.ErrorCode,
		})

		if errorInfo.StatusCode == 404 {
			return fmt.Errorf("%w: %s", sentinels.ErrEntraDeviceNotFound, entraDeviceObjectID)
		}

		return fmt.Errorf("%w %s: %w", sentinels.ErrRetrieveEntraDevice, entraDeviceObjectID, err)
	}

	if entraDevice == nil {
		tflog.Error(ctx, "Entra ID device not found", map[string]any{
			"entraDeviceObjectId": entraDeviceObjectID,
		})
		return fmt.Errorf("%w: %s", sentinels.ErrEntraDeviceNotFound, entraDeviceObjectID)
	}

	tflog.Debug(ctx, "Entra ID device validated successfully", map[string]any{
		"entraDeviceObjectId": entraDeviceObjectID,
		"displayName":         entraDevice.GetDisplayName(),
	})

	return nil
}
