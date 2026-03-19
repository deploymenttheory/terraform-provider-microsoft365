package graphBetaWindowsUpdatesDeviceEnrollment

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToDataSource(ctx context.Context, device graphmodelswindowsupdates.UpdatableAssetable) ([]UpdateManagementEnrollment, []UpdatableAssetError) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to data source for %s", DataSourceName))

	var enrollments []UpdateManagementEnrollment
	var assetErrors []UpdatableAssetError

	if azureADDevice, ok := device.(graphmodelswindowsupdates.AzureADDeviceable); ok {
		if enrollment := azureADDevice.GetEnrollment(); enrollment != nil {
			if driverEnrollment := enrollment.GetDriver(); driverEnrollment != nil {
				enrollments = append(enrollments, UpdateManagementEnrollment{
					UpdateCategory: types.StringValue("driver"),
				})
			}
			if featureEnrollment := enrollment.GetFeature(); featureEnrollment != nil {
				enrollments = append(enrollments, UpdateManagementEnrollment{
					UpdateCategory: types.StringValue("feature"),
				})
			}
			if qualityEnrollment := enrollment.GetQuality(); qualityEnrollment != nil {
				enrollments = append(enrollments, UpdateManagementEnrollment{
					UpdateCategory: types.StringValue("quality"),
				})
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to data source for %s", DataSourceName))
	return enrollments, assetErrors
}
