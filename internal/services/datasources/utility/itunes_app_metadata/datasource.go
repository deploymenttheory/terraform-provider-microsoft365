package itunes_app_metadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	// Basic datasource interface (Read operations)
	_ datasource.DataSource = &itunesAppMetadataDataSource{}

	// Allows the resource to be configured with the provider client
	_ datasource.DataSourceWithConfigure = &itunesAppMetadataDataSource{}
)

// NewItunesAppMetadataDataSource is a helper function to create the data source
func NewItunesAppMetadataDataSource() datasource.DataSource {
	return &itunesAppMetadataDataSource{}
}

// itunesAppMetadataDataSource is the data source implementation
type itunesAppMetadataDataSource struct {
	// No client needed as this is a public API
}

// Metadata returns the data source type name
func (d *itunesAppMetadataDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_utility_itunes_app_metadata"
}

// Configure adds the provider configured client to the data source
func (d *itunesAppMetadataDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// No client needed as this is a public API
	if req.ProviderData == nil {
		return
	}
}

// Schema defines the schema for the data source
func (d *itunesAppMetadataDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to query the iTunes App Store API for app metadata. " +
			"This data source allows you to search for apps by name and country code, returning details like bundle ID and artwork URLs.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this resource.",
			},
			"search_term": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The search term to use when querying the iTunes App Store API.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"country_code": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The two-letter country code for the App Store regionto search (e.g., 'us', 'gb', 'jp').",
				Validators: []validator.String{
					stringvalidator.LengthBetween(2, 2),
				},
			},
			"results": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of app results returned from the iTunes App Store API.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"track_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the app.",
						},
						"track_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the app.",
						},
						"bundle_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The bundle identifier of the app.",
						},
						"artwork_url_60": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "URL for the 60x60 app icon.",
						},
						"artwork_url_100": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "URL for the 100x100 app icon.",
						},
						"artwork_url_512": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "URL for the 512x512 app icon.",
						},
						"seller_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the app's seller/developer.",
						},
						"primary_genre": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The primary genre of the app.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the app.",
						},
						"version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The version of the app.",
						},
						"price": schema.Float64Attribute{
							Computed:            true,
							MarkdownDescription: "The price of the app in the local currency.",
						},
						"formatted_price": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The formatted price of the app (e.g., 'Free', '$0.99').",
						},
						"release_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The release date of the app.",
						},
						"average_rating": schema.Float64Attribute{
							Computed:            true,
							MarkdownDescription: "The average user rating of the app.",
						},
						"artist_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the artist/developer.",
						},
						"minimum_os_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The minimum OS version required to run the app.",
						},
						"content_advisory_rating": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The content advisory rating (e.g., '4+', '12+', '17+').",
						},
						"is_vpp_device_based_licensed": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether the app supports VPP device-based licensing.",
						},
						"release_notes": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Notes about the latest release of the app.",
						},
						"currency": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The currency code for the price (e.g., 'USD', 'GBP', 'EUR').",
						},
						"user_rating_count": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of user ratings for the app.",
						},
						"track_view_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL to view the app in the App Store.",
						},
						"screenshot_urls": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of screenshot URLs for iPhone.",
						},
						"ipad_screenshot_urls": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of screenshot URLs for iPad.",
						},
						"appletv_screenshot_urls": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of screenshot URLs for Apple TV.",
						},
						"supported_devices": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of device models that support the app.",
						},
						"features": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of app features (e.g., 'iosUniversal').",
						},
						"advisories": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of content advisories.",
						},
						"kind": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of content (e.g., 'software').",
						},
						"seller_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "URL to the seller's website.",
						},
						"is_game_center_enabled": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Whether Game Center is enabled for the app.",
						},
						"average_user_rating_for_current_version": schema.Float64Attribute{
							Computed:            true,
							MarkdownDescription: "The average user rating for the current version of the app.",
						},
						"user_rating_count_for_current_version": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The number of user ratings for the current version of the app.",
						},
						"file_size_bytes": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The size of the app in bytes.",
						},
						"language_codes_iso2a": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of supported language codes in ISO 2A format.",
						},
						"track_content_rating": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The content rating for the app (e.g., '4+', '12+', '17+').",
						},
						"artist_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The unique identifier for the artist/developer.",
						},
						"artist_view_url": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The URL to view the artist/developer in the App Store.",
						},
						"genres": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of genres for the app.",
						},
						"primary_genre_id": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The ID of the primary genre.",
						},
						"genre_ids": schema.ListAttribute{
							ElementType:         types.StringType,
							Computed:            true,
							MarkdownDescription: "List of genre IDs for the app.",
						},
						"track_censored_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The censored name of the app.",
						},
						"current_version_release_date": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The release date of the current version.",
						},
						"wrapper_type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The type of wrapper (e.g., 'software').",
						},
					},
				},
			},
		},
	}
}
