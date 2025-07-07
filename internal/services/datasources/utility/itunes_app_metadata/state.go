package itunes_app_metadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapResponseToState maps the iTunes API response to the Terraform state
func mapResponseToState(ctx context.Context, response ItunesSearchResponse, state *ItunesAppMetadataDataSourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	appResults := make([]ItunesAppResult, 0, response.ResultCount)
	for _, app := range response.Results {
		// Map screenshot URLs
		screenshotUrls, urlDiags := types.ListValueFrom(ctx, types.StringType, app.ScreenshotUrls)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		// Map iPad screenshot URLs
		ipadScreenshotUrls, urlDiags := types.ListValueFrom(ctx, types.StringType, app.IpadScreenshotUrls)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		// Map Apple TV screenshot URLs
		appletvScreenshotUrls, urlDiags := types.ListValueFrom(ctx, types.StringType, app.AppletvScreenshotUrls)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		// Map supported devices
		supportedDevices, urlDiags := types.ListValueFrom(ctx, types.StringType, app.SupportedDevices)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		// Map features
		features, urlDiags := types.ListValueFrom(ctx, types.StringType, app.Features)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		// Map advisories
		advisories, urlDiags := types.ListValueFrom(ctx, types.StringType, app.Advisories)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		// Map language codes
		languageCodes, urlDiags := types.ListValueFrom(ctx, types.StringType, app.LanguageCodesISO2A)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		// Map genres
		genres, urlDiags := types.ListValueFrom(ctx, types.StringType, app.Genres)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		// Map genre IDs
		genreIds, urlDiags := types.ListValueFrom(ctx, types.StringType, app.GenreIds)
		diags.Append(urlDiags...)
		if diags.HasError() {
			return diags
		}

		appResult := ItunesAppResult{
			TrackId:                            types.Int64Value(app.TrackId),
			TrackName:                          types.StringValue(app.TrackName),
			BundleId:                           types.StringValue(app.BundleId),
			ArtworkUrl60:                       types.StringValue(app.ArtworkUrl60),
			ArtworkUrl100:                      types.StringValue(app.ArtworkUrl100),
			ArtworkUrl512:                      types.StringValue(app.ArtworkUrl512),
			SellerName:                         types.StringValue(app.SellerName),
			PrimaryGenre:                       types.StringValue(app.PrimaryGenre),
			Description:                        types.StringValue(app.Description),
			Version:                            types.StringValue(app.Version),
			Price:                              types.Float64Value(app.Price),
			FormattedPrice:                     types.StringValue(app.FormattedPrice),
			ReleaseDate:                        types.StringValue(app.ReleaseDate),
			AverageRating:                      types.Float64Value(app.AverageRating),
			ArtistName:                         types.StringValue(app.ArtistName),
			MinimumOsVersion:                   types.StringValue(app.MinimumOsVersion),
			ContentAdvisoryRating:              types.StringValue(app.ContentAdvisoryRating),
			IsVppDeviceBasedLicensed:           types.BoolValue(app.IsVppDeviceBasedLicensing),
			ReleaseNotes:                       types.StringValue(app.ReleaseNotes),
			Currency:                           types.StringValue(app.Currency),
			UserRatingCount:                    types.Int64Value(app.UserRatingCount),
			TrackViewUrl:                       types.StringValue(app.TrackViewUrl),
			ScreenshotUrls:                     screenshotUrls,
			IpadScreenshotUrls:                 ipadScreenshotUrls,
			AppletvScreenshotUrls:              appletvScreenshotUrls,
			SupportedDevices:                   supportedDevices,
			Features:                           features,
			Advisories:                         advisories,
			Kind:                               types.StringValue(app.Kind),
			SellerUrl:                          types.StringValue(app.SellerUrl),
			IsGameCenterEnabled:                types.BoolValue(app.IsGameCenterEnabled),
			AverageUserRatingForCurrentVersion: types.Float64Value(app.AverageUserRatingForCurrentVersion),
			UserRatingCountForCurrentVersion:   types.Int64Value(app.UserRatingCountForCurrentVersion),
			FileSizeBytes:                      types.StringValue(app.FileSizeBytes),
			LanguageCodesISO2A:                 languageCodes,
			TrackContentRating:                 types.StringValue(app.TrackContentRating),
			ArtistId:                           types.Int64Value(app.ArtistId),
			ArtistViewUrl:                      types.StringValue(app.ArtistViewUrl),
			Genres:                             genres,
			PrimaryGenreId:                     types.Int64Value(int64(app.PrimaryGenreId)),
			GenreIds:                           genreIds,
			TrackCensoredName:                  types.StringValue(app.TrackCensoredName),
			CurrentVersionReleaseDate:          types.StringValue(app.CurrentVersionReleaseDate),
			WrapperType:                        types.StringValue(app.WrapperType),
		}
		appResults = append(appResults, appResult)
	}

	// Convert the slice of ItunesAppResult to types.List
	resultsValue, diags := types.ListValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"track_id":                     types.Int64Type,
			"track_name":                   types.StringType,
			"bundle_id":                    types.StringType,
			"artwork_url_60":               types.StringType,
			"artwork_url_100":              types.StringType,
			"artwork_url_512":              types.StringType,
			"seller_name":                  types.StringType,
			"primary_genre":                types.StringType,
			"description":                  types.StringType,
			"version":                      types.StringType,
			"price":                        types.Float64Type,
			"formatted_price":              types.StringType,
			"release_date":                 types.StringType,
			"average_rating":               types.Float64Type,
			"artist_name":                  types.StringType,
			"minimum_os_version":           types.StringType,
			"content_advisory_rating":      types.StringType,
			"is_vpp_device_based_licensed": types.BoolType,
			"release_notes":                types.StringType,
			"currency":                     types.StringType,
			"user_rating_count":            types.Int64Type,
			"track_view_url":               types.StringType,
			"screenshot_urls":              types.ListType{ElemType: types.StringType},
			"ipad_screenshot_urls":         types.ListType{ElemType: types.StringType},
			"appletv_screenshot_urls":      types.ListType{ElemType: types.StringType},
			"supported_devices":            types.ListType{ElemType: types.StringType},
			"features":                     types.ListType{ElemType: types.StringType},
			"advisories":                   types.ListType{ElemType: types.StringType},
			"kind":                         types.StringType,
			"seller_url":                   types.StringType,
			"is_game_center_enabled":       types.BoolType,
			"average_user_rating_for_current_version": types.Float64Type,
			"user_rating_count_for_current_version":   types.Int64Type,
			"file_size_bytes":                         types.StringType,
			"language_codes_iso2a":                    types.ListType{ElemType: types.StringType},
			"track_content_rating":                    types.StringType,
			"artist_id":                               types.Int64Type,
			"artist_view_url":                         types.StringType,
			"genres":                                  types.ListType{ElemType: types.StringType},
			"primary_genre_id":                        types.Int64Type,
			"genre_ids":                               types.ListType{ElemType: types.StringType},
			"track_censored_name":                     types.StringType,
			"current_version_release_date":            types.StringType,
			"wrapper_type":                            types.StringType,
		},
	}, appResults)

	if diags.HasError() {
		return diags
	}

	state.Results = resultsValue
	return diags
}
