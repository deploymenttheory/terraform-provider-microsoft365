package graphBetaWindowsUpdatesApplicableContent

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodelswindowsupdates "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToDataSource(ctx context.Context, remoteContent graphmodelswindowsupdates.ApplicableContentable) ApplicableContent {
	tflog.Debug(ctx, fmt.Sprintf("Mapping remote state to data source for %s", DataSourceName))

	content := ApplicableContent{
		CatalogEntryId: convert.GraphToFrameworkString(remoteContent.GetCatalogEntryId()),
	}

	if catalogEntry := remoteContent.GetCatalogEntry(); catalogEntry != nil {
		if driverEntry, ok := catalogEntry.(graphmodelswindowsupdates.DriverUpdateCatalogEntryable); ok {
			content.CatalogEntry = &CatalogEntry{
				ID:                      convert.GraphToFrameworkString(driverEntry.GetId()),
				DisplayName:             convert.GraphToFrameworkString(driverEntry.GetDisplayName()),
				ReleaseDateTime:         convert.GraphToFrameworkTime(driverEntry.GetReleaseDateTime()),
				DeployableUntilDateTime: convert.GraphToFrameworkTime(driverEntry.GetDeployableUntilDateTime()),
				Description:             convert.GraphToFrameworkString(driverEntry.GetDescription()),
				DriverClass:             convert.GraphToFrameworkString(driverEntry.GetDriverClass()),
				Provider:                convert.GraphToFrameworkString(driverEntry.GetProvider()),
				Manufacturer:            convert.GraphToFrameworkString(driverEntry.GetManufacturer()),
				Version:                 convert.GraphToFrameworkString(driverEntry.GetVersion()),
				VersionDateTime:         convert.GraphToFrameworkTime(driverEntry.GetVersionDateTime()),
			}
		}
	}

	if matchedDevices := remoteContent.GetMatchedDevices(); matchedDevices != nil {
		for _, device := range matchedDevices {
			matchedDevice := MatchedDevice{
				DeviceId: convert.GraphToFrameworkString(device.GetDeviceId()),
			}

			if recommendedBy := device.GetRecommendedBy(); recommendedBy != nil {
				matchedDevice.RecommendedBy = convert.GraphToFrameworkStringSlice(recommendedBy)
			}

			content.MatchedDevices = append(content.MatchedDevices, matchedDevice)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to data source for %s", DataSourceName))
	return content
}
