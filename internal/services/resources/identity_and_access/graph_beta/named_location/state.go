package graphBetaNamedLocation

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote named location to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *NamedLocationResourceModel, remoteResource graphmodels.NamedLocationable) {
	if remoteResource == nil {
		return
	}

	// Basic properties using helpers
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.ModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetModifiedDateTime())

	// Determine the type of named location based on the concrete type
	switch location := remoteResource.(type) {
	case *graphmodels.IpNamedLocation:
		mapIPNamedLocationFields(ctx, data, location)
	case *graphmodels.CountryNamedLocation:
		mapCountryNamedLocationFields(ctx, data, location)
	default:
		// For base NamedLocation or unknown types, check AdditionalData for @odata.type
		if additionalData := remoteResource.GetAdditionalData(); additionalData != nil {
			if odataType, ok := additionalData["@odata.type"].(string); ok {
				switch odataType {
				case "#microsoft.graph.ipNamedLocation":
					// Try to cast or handle as IP location
					if ipLoc, ok := remoteResource.(*graphmodels.IpNamedLocation); ok {
						mapIPNamedLocationFields(ctx, data, ipLoc)
					}
				case "#microsoft.graph.countryNamedLocation":
					// Try to cast or handle as country location
					if countryLoc, ok := remoteResource.(*graphmodels.CountryNamedLocation); ok {
						mapCountryNamedLocationFields(ctx, data, countryLoc)
					}
				}
			}
		}
	}
}

// mapIPNamedLocationFields maps IP named location specific fields
func mapIPNamedLocationFields(ctx context.Context, data *NamedLocationResourceModel, ipLocation *graphmodels.IpNamedLocation) {
	data.IsTrusted = convert.GraphToFrameworkBool(ipLocation.GetIsTrusted())

	// Parse IP ranges from the API response
	ipRanges := ipLocation.GetIpRanges()
	if ipRanges != nil {
		var ipv4Ranges []string
		var ipv6Ranges []string

		for _, rangeItem := range ipRanges {
			switch ipRange := rangeItem.(type) {
			case *graphmodels.IPv4CidrRange:
				if cidr := ipRange.GetCidrAddress(); cidr != nil {
					ipv4Ranges = append(ipv4Ranges, *cidr)
				}
			case *graphmodels.IPv6CidrRange:
				if cidr := ipRange.GetCidrAddress(); cidr != nil {
					ipv6Ranges = append(ipv6Ranges, *cidr)
				}
			}
		}

		// Convert to Terraform sets
		data.IPv4Ranges = convert.GraphToFrameworkStringSet(ctx, ipv4Ranges)
		data.IPv6Ranges = convert.GraphToFrameworkStringSet(ctx, ipv6Ranges)
	} else {
		// Set empty sets if no IP ranges are present
		data.IPv4Ranges = types.SetNull(types.StringType)
		data.IPv6Ranges = types.SetNull(types.StringType)
	}

	// Ensure country fields are null for IP locations
	data.CountryLookupMethod = types.StringNull()
	data.IncludeUnknownCountriesAndRegions = types.BoolNull()
	data.CountriesAndRegions = types.SetNull(types.StringType)
}

// mapCountryNamedLocationFields maps country named location specific fields
func mapCountryNamedLocationFields(ctx context.Context, data *NamedLocationResourceModel, countryLocation *graphmodels.CountryNamedLocation) {
	data.CountryLookupMethod = convert.GraphToFrameworkEnum(
		countryLocation.GetCountryLookupMethod(),
	)
	data.IncludeUnknownCountriesAndRegions = convert.GraphToFrameworkBool(
		countryLocation.GetIncludeUnknownCountriesAndRegions(),
	)

	// Parse countries and regions from the API response
	countriesAndRegions := countryLocation.GetCountriesAndRegions()
	if countriesAndRegions != nil {
		data.CountriesAndRegions = convert.GraphToFrameworkStringSet(ctx, countriesAndRegions)
	} else {
		data.CountriesAndRegions = types.SetNull(types.StringType)
	}

	// Ensure IP fields are null for country locations
	data.IsTrusted = types.BoolNull()
	data.IPv4Ranges = types.SetNull(types.StringType)
	data.IPv6Ranges = types.SetNull(types.StringType)
}
