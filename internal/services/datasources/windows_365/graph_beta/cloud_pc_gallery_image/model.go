// REF: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcgalleryimage?view=graph-rest-beta

package graphBetaCloudPcGalleryImage

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CloudPcGalleryImageDataSourceModel represents the Terraform data source model for gallery images
type CloudPcGalleryImageDataSourceModel struct {
	FilterType  types.String                   `tfsdk:"filter_type"`
	FilterValue types.String                   `tfsdk:"filter_value"`
	Items       []CloudPcGalleryImageItemModel `tfsdk:"items"`
	Timeouts    timeouts.Value                 `tfsdk:"timeouts"`
}

// CloudPcGalleryImageItemModel represents an individual gallery image
type CloudPcGalleryImageItemModel struct {
	ID              types.String `tfsdk:"id"`
	DisplayName     types.String `tfsdk:"display_name"`
	StartDate       types.String `tfsdk:"start_date"`
	EndDate         types.String `tfsdk:"end_date"`
	ExpirationDate  types.String `tfsdk:"expiration_date"`
	OSVersionNumber types.String `tfsdk:"os_version_number"`
	PublisherName   types.String `tfsdk:"publisher_name"`
	OfferName       types.String `tfsdk:"offer_name"`
	SkuName         types.String `tfsdk:"sku_name"`
	SizeInGB        types.Int64  `tfsdk:"size_in_gb"`
	Status          types.String `tfsdk:"status"`
}
