package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroup

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	sdkdevices "github.com/microsoftgraph/msgraph-beta-sdk-go/devices"
)

// validateRequest performs pre-flight validation before creating or updating the
// updatable asset group. It verifies that every device object ID in entra_device_object_ids
// resolves to a real Entra ID device object.
func validateRequest(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	data *WindowsUpdatesAutopatchUpdatableAssetGroupResourceModel,
) error {
	tflog.Debug(ctx, fmt.Sprintf("Starting validation for %s resource", ResourceName))

	deviceObjectIDs := extractDeviceObjectIDs(data)
	if len(deviceObjectIDs) > 0 {
		if err := validateDeviceObjectIds(ctx, client, deviceObjectIDs); err != nil {
			return fmt.Errorf("entra_device_object_ids validation failed: %w", err)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validation completed successfully for %s resource", ResourceName))
	return nil
}

// validateDeviceObjectIds verifies that all supplied Entra ID device object IDs exist by
// querying the /devices endpoint individually for each ID using a deviceId eq filter.
// A single error is returned listing any IDs that could not be found.
func validateDeviceObjectIds(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	deviceObjectIDs []string,
) error {
	if len(deviceObjectIDs) == 0 {
		return nil
	}

	tflog.Debug(ctx, "Validating Entra ID device object IDs exist", map[string]any{
		"count": len(deviceObjectIDs),
	})

	var missing []string

	for _, deviceID := range deviceObjectIDs {
		filter := fmt.Sprintf("deviceId eq '%s'", deviceID)

		requestConfig := &sdkdevices.DevicesRequestBuilderGetRequestConfiguration{
			QueryParameters: &sdkdevices.DevicesRequestBuilderGetQueryParameters{
				Filter: &filter,
				Select: []string{"deviceId"},
			},
		}

		resp, err := client.Devices().Get(ctx, requestConfig)
		if err != nil {
			return fmt.Errorf("%w: failed to query Entra ID device '%s': %w",
				sentinels.ErrInvalidEntraDeviceIDs, deviceID, err)
		}

		if resp == nil || len(resp.GetValue()) == 0 {
			missing = append(missing, deviceID)
		}
	}

	if len(missing) > 0 {
		tflog.Error(ctx, "One or more Entra ID device object IDs not found", map[string]any{
			"missingDeviceObjectIds": missing,
		})
		return fmt.Errorf("%w: %s", sentinels.ErrInvalidEntraDeviceIDs, strings.Join(missing, ", "))
	}

	tflog.Debug(ctx, "All Entra ID device object IDs validated successfully", map[string]any{
		"count": len(deviceObjectIDs),
	})
	return nil
}
