package graphBetaCloudPcGalleryImage

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.CloudPcGalleryImageable) CloudPcGalleryImageItemModel {
	model := CloudPcGalleryImageItemModel{
		ID:              convert.GraphToFrameworkString(data.GetId()),
		DisplayName:     convert.GraphToFrameworkString(data.GetDisplayName()),
		OSVersionNumber: convert.GraphToFrameworkString(data.GetOsVersionNumber()),
		PublisherName:   convert.GraphToFrameworkString(data.GetPublisherName()),
		OfferName:       convert.GraphToFrameworkString(data.GetOfferName()),
		SkuName:         convert.GraphToFrameworkString(data.GetSkuName()),
	}

	// Convert dates to strings
	model.StartDate = convert.GraphToFrameworkDateOnly(data.GetStartDate())
	model.EndDate = convert.GraphToFrameworkDateOnly(data.GetEndDate())
	model.ExpirationDate = convert.GraphToFrameworkDateOnly(data.GetExpirationDate())

	// Convert Int32 to Int64 for Terraform
	if sizeInGB := data.GetSizeInGB(); sizeInGB != nil {
		model.SizeInGB = convert.GraphToFrameworkInt32AsInt64(sizeInGB)
	}

	// Enum: Status
	model.Status = convert.GraphToFrameworkEnum(data.GetStatus())

	return model
}
