package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

type WindowsAutopatchDeviceRegistrationTestResource struct{}

// Exists checks whether any of the devices from state are still enrolled in Windows Updates
// for the given update category.
//
// Returns:
//   - (*true, nil)  – at least one device is still enrolled → resource still exists
//   - (*false, nil) – no devices enrolled → resource has been destroyed
//   - (nil, err)    – API call or client setup failed
func (r WindowsAutopatchDeviceRegistrationTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create graph client: %w", err)
	}

	updateCategory := state.Attributes["update_category"]

	deviceIDsCount := 0
	for key := range state.Attributes {
		if key == "entra_device_object_ids.#" {
			fmt.Sscanf(state.Attributes[key], "%d", &deviceIDsCount)
			break
		}
	}

	if deviceIDsCount == 0 {
		f := false
		return &f, nil
	}

	stateDeviceIDs := make(map[string]bool)
	for i := 0; i < deviceIDsCount; i++ {
		deviceID := state.Attributes[fmt.Sprintf("entra_device_object_ids.%d", i)]
		if deviceID != "" {
			stateDeviceIDs[deviceID] = true
		}
	}

	filterQuery := "isof('microsoft.graph.windowsUpdates.azureADDevice')"
	assetsResp, err := graphClient.
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
		return nil, fmt.Errorf("failed to get updatable assets: %w", err)
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
		t := true
		return &t, nil
	}

	f := false
	return &f, nil
}
