package graphBetaWindowsUpdatesAutopatchUpdatableAssetGroupAssignment

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	sdkdevices "github.com/microsoftgraph/msgraph-beta-sdk-go/devices"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// validateRequest performs pre-flight validation before creating or updating the
// updatable asset group assignment. It verifies that the supplied updatable_asset_group_id
// references an existing Windows Updates updatable asset group, and that every
// device ID in entra_device_ids resolves to a real Entra ID device object.
func validateRequest(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	data *WindowsUpdatesAutopatchUpdatableAssetGroupAssignmentResourceModel,
) error {
	tflog.Debug(ctx, fmt.Sprintf("Starting validation for %s resource", ResourceName))

	if err := validateValidUpdatableAssetGroupId(ctx, client, data.UpdatableAssetGroupId.ValueString()); err != nil {
		return fmt.Errorf("updatable_asset_group_id validation failed: %w", err)
	}

	deviceIDs := extractDeviceIDs(data)
	if len(deviceIDs) > 0 {
		if err := validateValidDeviceId(ctx, client, deviceIDs); err != nil {
			return fmt.Errorf("entra_device_ids validation failed: %w", err)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validation completed successfully for %s resource", ResourceName))
	return nil
}

// validateValidUpdatableAssetGroupId verifies that the supplied updatable_asset_group_id exists as a Windows Updates
// updatable asset group by performing a paginated GET across all updatable assets and
// checking for a matching ID.
func validateValidUpdatableAssetGroupId(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	groupID string,
) error {
	if groupID == "" {
		return fmt.Errorf("%w: updatable_asset_group_id cannot be empty", sentinels.ErrUpdatableAssetGroupValidationFailed)
	}

	tflog.Debug(ctx, "Validating updatable asset group exists via paginated GET", map[string]any{
		"groupId": groupID,
	})

	requestConfig := &graphadmin.WindowsUpdatesUpdatableAssetsRequestBuilderGetRequestConfiguration{
		QueryParameters: &graphadmin.WindowsUpdatesUpdatableAssetsRequestBuilderGetQueryParameters{
			Select: []string{"id"},
		},
	}

	assetsResponse, err := client.
		Admin().
		Windows().
		Updates().
		UpdatableAssets().
		Get(ctx, requestConfig)

	if err != nil {
		return fmt.Errorf("%w: failed to query updatable assets: %w",
			sentinels.ErrUpdatableAssetGroupValidationFailed, err)
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodelswindowsupdates.UpdatableAssetable](
		assetsResponse,
		client.GetAdapter(),
		graphmodelswindowsupdates.CreateUpdatableAssetCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return fmt.Errorf("%w: failed to create page iterator for updatable assets: %w",
			sentinels.ErrUpdatableAssetGroupValidationFailed, err)
	}

	found := false
	err = pageIterator.Iterate(ctx, func(item graphmodelswindowsupdates.UpdatableAssetable) bool {
		if item != nil && item.GetId() != nil && *item.GetId() == groupID {
			found = true
			return false // stop iterating
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("%w: error iterating updatable assets pages: %w",
			sentinels.ErrUpdatableAssetGroupValidationFailed, err)
	}

	if !found {
		tflog.Error(ctx, "Updatable asset group not found", map[string]any{
			"groupId": groupID,
		})
		return fmt.Errorf("%w with ID '%s'", sentinels.ErrUpdatableAssetGroupNotFound, groupID)
	}

	tflog.Debug(ctx, "Updatable asset group validated successfully", map[string]any{
		"groupId": groupID,
	})
	return nil
}

// validateValidDeviceId verifies that all supplied Entra ID device IDs exist by
// querying the /devices endpoint individually for each ID using a deviceId eq filter.
// A single error is returned listing any IDs that could not be found.
func validateValidDeviceId(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	deviceIDs []string,
) error {
	if len(deviceIDs) == 0 {
		return nil
	}

	tflog.Debug(ctx, "Validating Entra ID device IDs exist", map[string]any{
		"count": len(deviceIDs),
	})

	var missing []string

	for _, deviceID := range deviceIDs {
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
		tflog.Error(ctx, "One or more Entra ID device IDs not found", map[string]any{
			"missingDeviceIds": missing,
		})
		return fmt.Errorf("%w: %s", sentinels.ErrInvalidEntraDeviceIDs, strings.Join(missing, ", "))
	}

	tflog.Debug(ctx, "All Entra ID device IDs validated successfully", map[string]any{
		"count": len(deviceIDs),
	})
	return nil
}
