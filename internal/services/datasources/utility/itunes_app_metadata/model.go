package itunes_app_metadata

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ItunesAppMetadataDataSourceModel represents the data source model for iTunes app metadata
type ItunesAppMetadataDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	SearchTerm  types.String `tfsdk:"search_term"`
	CountryCode types.String `tfsdk:"country_code"`
	Results     types.List   `tfsdk:"results"`
}

// ItunesAppResult represents an individual app result from the iTunes API
type ItunesAppResult struct {
	TrackId                            types.Int64   `tfsdk:"track_id"`
	TrackName                          types.String  `tfsdk:"track_name"`
	BundleId                           types.String  `tfsdk:"bundle_id"`
	ArtworkUrl60                       types.String  `tfsdk:"artwork_url_60"`
	ArtworkUrl100                      types.String  `tfsdk:"artwork_url_100"`
	ArtworkUrl512                      types.String  `tfsdk:"artwork_url_512"`
	SellerName                         types.String  `tfsdk:"seller_name"`
	PrimaryGenre                       types.String  `tfsdk:"primary_genre"`
	Description                        types.String  `tfsdk:"description"`
	Version                            types.String  `tfsdk:"version"`
	Price                              types.Float64 `tfsdk:"price"`
	FormattedPrice                     types.String  `tfsdk:"formatted_price"`
	ReleaseDate                        types.String  `tfsdk:"release_date"`
	AverageRating                      types.Float64 `tfsdk:"average_rating"`
	ArtistName                         types.String  `tfsdk:"artist_name"`
	MinimumOsVersion                   types.String  `tfsdk:"minimum_os_version"`
	ContentAdvisoryRating              types.String  `tfsdk:"content_advisory_rating"`
	IsVppDeviceBasedLicensed           types.Bool    `tfsdk:"is_vpp_device_based_licensed"`
	ReleaseNotes                       types.String  `tfsdk:"release_notes"`
	Currency                           types.String  `tfsdk:"currency"`
	UserRatingCount                    types.Int64   `tfsdk:"user_rating_count"`
	TrackViewUrl                       types.String  `tfsdk:"track_view_url"`
	ScreenshotUrls                     types.List    `tfsdk:"screenshot_urls"`
	IpadScreenshotUrls                 types.List    `tfsdk:"ipad_screenshot_urls"`
	AppletvScreenshotUrls              types.List    `tfsdk:"appletv_screenshot_urls"`
	SupportedDevices                   types.List    `tfsdk:"supported_devices"`
	Features                           types.List    `tfsdk:"features"`
	Advisories                         types.List    `tfsdk:"advisories"`
	Kind                               types.String  `tfsdk:"kind"`
	SellerUrl                          types.String  `tfsdk:"seller_url"`
	IsGameCenterEnabled                types.Bool    `tfsdk:"is_game_center_enabled"`
	AverageUserRatingForCurrentVersion types.Float64 `tfsdk:"average_user_rating_for_current_version"`
	UserRatingCountForCurrentVersion   types.Int64   `tfsdk:"user_rating_count_for_current_version"`
	FileSizeBytes                      types.String  `tfsdk:"file_size_bytes"`
	LanguageCodesISO2A                 types.List    `tfsdk:"language_codes_iso2a"`
	TrackContentRating                 types.String  `tfsdk:"track_content_rating"`
	ArtistId                           types.Int64   `tfsdk:"artist_id"`
	ArtistViewUrl                      types.String  `tfsdk:"artist_view_url"`
	Genres                             types.List    `tfsdk:"genres"`
	PrimaryGenreId                     types.Int64   `tfsdk:"primary_genre_id"`
	GenreIds                           types.List    `tfsdk:"genre_ids"`
	TrackCensoredName                  types.String  `tfsdk:"track_censored_name"`
	CurrentVersionReleaseDate          types.String  `tfsdk:"current_version_release_date"`
	WrapperType                        types.String  `tfsdk:"wrapper_type"`
}

// ItunesSearchResponse represents the response from the iTunes Search API
type ItunesSearchResponse struct {
	ResultCount int             `json:"resultCount"`
	Results     []ItunesAppData `json:"results"`
}

// ItunesAppData represents the raw data from the iTunes API
type ItunesAppData struct {
	TrackId                            int64    `json:"trackId"`
	TrackName                          string   `json:"trackName"`
	BundleId                           string   `json:"bundleId"`
	ArtworkUrl60                       string   `json:"artworkUrl60"`
	ArtworkUrl100                      string   `json:"artworkUrl100"`
	ArtworkUrl512                      string   `json:"artworkUrl512"`
	SellerName                         string   `json:"sellerName"`
	PrimaryGenre                       string   `json:"primaryGenreName"`
	Description                        string   `json:"description"`
	Version                            string   `json:"version"`
	Price                              float64  `json:"price"`
	FormattedPrice                     string   `json:"formattedPrice"`
	ReleaseDate                        string   `json:"releaseDate"`
	AverageRating                      float64  `json:"averageUserRating"`
	ArtistName                         string   `json:"artistName"`
	MinimumOsVersion                   string   `json:"minimumOsVersion"`
	ContentAdvisoryRating              string   `json:"contentAdvisoryRating"`
	IsVppDeviceBasedLicensing          bool     `json:"isVppDeviceBasedLicensingEnabled"`
	ReleaseNotes                       string   `json:"releaseNotes"`
	Currency                           string   `json:"currency"`
	UserRatingCount                    int64    `json:"userRatingCount"`
	TrackViewUrl                       string   `json:"trackViewUrl"`
	ScreenshotUrls                     []string `json:"screenshotUrls"`
	IpadScreenshotUrls                 []string `json:"ipadScreenshotUrls"`
	AppletvScreenshotUrls              []string `json:"appletvScreenshotUrls"`
	SupportedDevices                   []string `json:"supportedDevices"`
	Features                           []string `json:"features"`
	Advisories                         []string `json:"advisories"`
	Kind                               string   `json:"kind"`
	SellerUrl                          string   `json:"sellerUrl"`
	IsGameCenterEnabled                bool     `json:"isGameCenterEnabled"`
	AverageUserRatingForCurrentVersion float64  `json:"averageUserRatingForCurrentVersion"`
	UserRatingCountForCurrentVersion   int64    `json:"userRatingCountForCurrentVersion"`
	FileSizeBytes                      string   `json:"fileSizeBytes"`
	LanguageCodesISO2A                 []string `json:"languageCodesISO2A"`
	TrackContentRating                 string   `json:"trackContentRating"`
	ArtistId                           int64    `json:"artistId"`
	ArtistViewUrl                      string   `json:"artistViewUrl"`
	Genres                             []string `json:"genres"`
	PrimaryGenreId                     int      `json:"primaryGenreId"`
	GenreIds                           []string `json:"genreIds"`
	TrackCensoredName                  string   `json:"trackCensoredName"`
	CurrentVersionReleaseDate          string   `json:"currentVersionReleaseDate"`
	WrapperType                        string   `json:"wrapperType"`
}
