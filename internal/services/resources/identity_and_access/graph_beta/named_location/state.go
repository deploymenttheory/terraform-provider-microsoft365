package graphBetaNamedLocation

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapRemoteResourceStateToTerraform maps the remote named location to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *NamedLocationResourceModel, remoteResource map[string]any) {
	// Basic properties using helpers
	if id, ok := remoteResource["id"].(string); ok {
		data.ID = types.StringValue(id)
	}

	data.DisplayName = convert.GraphToFrameworkString(getStringPtr(remoteResource, "displayName"))
	data.CreatedDateTime = convert.GraphToFrameworkString(getStringPtr(remoteResource, "createdDateTime"))
	data.ModifiedDateTime = convert.GraphToFrameworkString(getStringPtr(remoteResource, "modifiedDateTime"))

	// Determine the type of named location based on @odata.type
	if odataType, ok := remoteResource["@odata.type"].(string); ok {
		switch odataType {
		case "#microsoft.graph.ipNamedLocation":
			mapIPNamedLocationFields(ctx, data, remoteResource)
		case "#microsoft.graph.countryNamedLocation":
			mapCountryNamedLocationFields(ctx, data, remoteResource)
		}
	}
}

// mapIPNamedLocationFields maps IP named location specific fields
func mapIPNamedLocationFields(ctx context.Context, data *NamedLocationResourceModel, remoteResource map[string]any) {
	data.IsTrusted = convert.GraphToFrameworkBool(getBoolPtr(remoteResource, "isTrusted"))

	// Parse IP ranges from the API response
	if ipRanges, ok := remoteResource["ipRanges"].([]interface{}); ok {
		var ipv4Ranges []string
		var ipv6Ranges []string

		for _, rangeItem := range ipRanges {
			if rangeMap, ok := rangeItem.(map[string]any); ok {
				if odataType, typeOk := rangeMap["@odata.type"].(string); typeOk {
					if cidrAddress, addrOk := rangeMap["cidrAddress"].(string); addrOk {
						switch odataType {
						case "#microsoft.graph.iPv4CidrRange":
							ipv4Ranges = append(ipv4Ranges, cidrAddress)
						case "#microsoft.graph.iPv6CidrRange":
							ipv6Ranges = append(ipv6Ranges, cidrAddress)
						}
					}
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
func mapCountryNamedLocationFields(ctx context.Context, data *NamedLocationResourceModel, remoteResource map[string]any) {
	data.CountryLookupMethod = convert.GraphToFrameworkString(getStringPtr(remoteResource, "countryLookupMethod"))
	data.IncludeUnknownCountriesAndRegions = convert.GraphToFrameworkBool(getBoolPtr(remoteResource, "includeUnknownCountriesAndRegions"))

	// Parse countries and regions from the API response
	if countriesAndRegions, ok := remoteResource["countriesAndRegions"].([]interface{}); ok {
		var countries []string
		for _, countryItem := range countriesAndRegions {
			if country, ok := countryItem.(string); ok {
				countries = append(countries, country)
			}
		}
		data.CountriesAndRegions = convert.GraphToFrameworkStringSet(ctx, countries)
	} else {
		data.CountriesAndRegions = types.SetNull(types.StringType)
	}

	// Ensure IP fields are null for country locations
	data.IsTrusted = types.BoolNull()
	data.IPv4Ranges = types.SetNull(types.StringType)
	data.IPv6Ranges = types.SetNull(types.StringType)
}

// Helper function to get string pointer from map
func getStringPtr(data map[string]any, key string) *string {
	if value, ok := data[key].(string); ok {
		return &value
	}
	return nil
}

// Helper function to get bool pointer from map
func getBoolPtr(data map[string]any, key string) *bool {
	if value, ok := data[key].(bool); ok {
		return &value
	}
	return nil
}
