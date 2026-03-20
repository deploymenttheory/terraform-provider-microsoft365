package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

type WindowsAutopatchDeviceRegistrationTestResource struct{}

func (r WindowsAutopatchDeviceRegistrationTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		updateCategory := state.Attributes["update_category"]
		
		deviceIDsCount := 0
		for key := range state.Attributes {
			if key == "entra_device_object_ids.#" {
				fmt.Sscanf(state.Attributes[key], "%d", &deviceIDsCount)
				break
			}
		}
		
		if deviceIDsCount == 0 {
			return fmt.Errorf("%w: no device IDs in state", sentinels.ErrNoDevicesEnrolled)
		}
		
		stateDeviceIDs := make(map[string]bool)
		for i := 0; i < deviceIDsCount; i++ {
			deviceID := state.Attributes[fmt.Sprintf("entra_device_object_ids.%d", i)]
			if deviceID != "" {
				stateDeviceIDs[deviceID] = true
			}
		}

		filterQuery := "isof('microsoft.graph.windowsUpdates.azureADDevice')"
		assetsResp, err := client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			Get(ctx, &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetRequestConfiguration{
				QueryParameters: &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetQueryParameters{
					Filter: &filterQuery,
				},
			})

		if err != nil {
			return fmt.Errorf("failed to get updatable assets: %w", err)
		}

		devices := assetsResp.GetValue()
		enrolledCount := 0
		
		for _, asset := range devices {
			if asset == nil {
				continue
			}

			azureDevice, ok := asset.(windowsupdates.AzureADDeviceable)
			if !ok {
				continue
			}

			deviceID := azureDevice.GetId()
			if deviceID == nil || !stateDeviceIDs[*deviceID] {
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
						enrolledCount++
					}
				}
			}
		}

		if enrolledCount > 0 {
			return fmt.Errorf("found %d devices from state still enrolled for category %s", enrolledCount, updateCategory)
		}

		return nil
	})
}
