package graphBetaWindowsAutopatchDeviceRegistration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToTerraform(ctx context.Context, data *WindowsAutopatchDeviceRegistrationResourceModel, devices []graphmodelswindowsupdates.UpdatableAssetable) {
	tflog.Debug(ctx, "Starting to map remote state to Terraform state")

	updateCategory := data.UpdateCategory.ValueString()
	enrolledDeviceIDs := make([]attr.Value, 0)

	plannedDeviceIDs := make(map[string]bool)
	if !data.DeviceIds.IsNull() && !data.DeviceIds.IsUnknown() {
		elements := data.DeviceIds.Elements()
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				plannedDeviceIDs[strVal.ValueString()] = true
			}
		}
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

		if !plannedDeviceIDs[*deviceID] {
			continue
		}

		enrollment := azureDevice.GetEnrollment()
		if enrollment == nil {
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
			enrolledDeviceIDs = append(enrolledDeviceIDs, types.StringValue(*deviceID))
		}
	}

	if len(enrolledDeviceIDs) > 0 {
		data.DeviceIds = types.SetValueMust(types.StringType, enrolledDeviceIDs)
	} else {
		data.DeviceIds = types.SetValueMust(types.StringType, []attr.Value{})
	}

	data.ID = types.StringValue(fmt.Sprintf("%s", data.UpdateCategory.ValueString()))

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"enrolledDeviceCount": len(enrolledDeviceIDs),
	})
}
