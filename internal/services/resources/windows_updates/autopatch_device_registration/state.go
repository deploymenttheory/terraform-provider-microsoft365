package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsUpdatesAutopatchDeviceRegistrationResourceModel, devices []graphmodelswindowsupdates.UpdatableAssetable) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state")

	updateCategory := data.UpdateCategory.ValueString()
	enrolledDeviceIDs := make([]attr.Value, 0)

	plannedDeviceIDs := make(map[string]bool)
	filterByPlanned := false
	if !data.EntraDeviceObjectIds.IsNull() && !data.EntraDeviceObjectIds.IsUnknown() {
		elements := data.EntraDeviceObjectIds.Elements()
		if len(elements) > 0 {
			filterByPlanned = true
			for _, elem := range elements {
				if strVal, ok := elem.(types.String); ok {
					plannedDeviceIDs[strVal.ValueString()] = true
				}
			}
			tflog.Debug(ctx, fmt.Sprintf("Filtering by %d planned device IDs", len(plannedDeviceIDs)))
		}
	} else {
		tflog.Debug(ctx, "No planned device IDs to filter by - will include all enrolled devices")
	}

	for _, asset := range devices {
		if asset == nil {
			continue
		}

		azureDevice, ok := asset.(graphmodelswindowsupdates.AzureADDeviceable)
		if !ok {
			continue
		}

		deviceID := azureDevice.GetId()
		if deviceID == nil {
			continue
		}

		if filterByPlanned && !plannedDeviceIDs[*deviceID] {
			tflog.Debug(ctx, fmt.Sprintf("Skipping device %s - not in planned device IDs", *deviceID))
			continue
		}

		enrollment := azureDevice.GetEnrollment()
		if enrollment == nil {
			tflog.Debug(ctx, fmt.Sprintf("Skipping device %s - no enrollment information", *deviceID))
			continue
		}

		var categoryEnrollment graphmodelswindowsupdates.UpdateCategoryEnrollmentInformationable
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
					tflog.Debug(ctx, fmt.Sprintf("Including device %s - enrollmentState=%s for %s", *deviceID, stateStr, updateCategory))
					enrolledDeviceIDs = append(enrolledDeviceIDs, types.StringValue(*deviceID))
				} else {
					tflog.Debug(ctx, fmt.Sprintf("Skipping device %s - enrollmentState=%s for %s", *deviceID, stateStr, updateCategory))
				}
			} else {
				tflog.Debug(ctx, fmt.Sprintf("Skipping device %s - enrollmentState is nil for %s", *deviceID, updateCategory))
			}
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Skipping device %s - no category enrollment object for %s", *deviceID, updateCategory))
		}
	}

	if len(enrolledDeviceIDs) > 0 {
		data.EntraDeviceObjectIds = types.SetValueMust(types.StringType, enrolledDeviceIDs)
	} else {
		data.EntraDeviceObjectIds = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.ID = types.StringValue(fmt.Sprintf("%s", data.UpdateCategory.ValueString()))

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"enrolledDeviceCount": len(enrolledDeviceIDs),
	})
}
