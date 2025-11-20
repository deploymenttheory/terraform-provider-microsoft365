package graphBetaNamedLocation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateRequest validates the named location request to ensure that exactly one
// type of location (IP or Country) is configured
func validateRequest(ctx context.Context, data *NamedLocationResourceModel) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating %s resource configuration", ResourceName))

	// Determine the type of named location based on which fields are populated
	// We check both IsNull() and IsUnknown() because Computed fields can be Unknown during planning
	isIPLocation := (!data.IPv4Ranges.IsNull() && !data.IPv4Ranges.IsUnknown()) ||
		(!data.IPv6Ranges.IsNull() && !data.IPv6Ranges.IsUnknown()) ||
		(!data.IsTrusted.IsNull() && !data.IsTrusted.IsUnknown())
	isCountryLocation := (!data.CountryLookupMethod.IsNull() && !data.CountryLookupMethod.IsUnknown()) ||
		(!data.CountriesAndRegions.IsNull() && !data.CountriesAndRegions.IsUnknown()) ||
		(!data.IncludeUnknownCountriesAndRegions.IsNull() && !data.IncludeUnknownCountriesAndRegions.IsUnknown())

	switch {
	case isIPLocation && isCountryLocation:
		return fmt.Errorf("cannot specify both IP location fields and country location fields in the same resource")
	case !isIPLocation && !isCountryLocation:
		return fmt.Errorf("must specify either IP location fields (ipv4_ranges, ipv6_ranges, is_trusted) or country location fields (country_lookup_method, countries_and_regions, include_unknown_countries_and_regions)")
	default:
		tflog.Debug(ctx, fmt.Sprintf("Validation passed for %s resource", ResourceName))
		return nil
	}
}
