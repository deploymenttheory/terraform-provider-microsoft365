package graphBetaNamedLocation

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource converts the Terraform resource model to a Kiota SDK model
// Returns a NamedLocationable interface that can be used with the Kiota SDK
func constructResource(ctx context.Context, data *NamedLocationResourceModel) (graphmodels.NamedLocationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, data); err != nil {
		return nil, err
	}

	// Determine the type of named location based on which fields are populated
	isIPLocation := (!data.IPv4Ranges.IsNull() && !data.IPv4Ranges.IsUnknown()) ||
		(!data.IPv6Ranges.IsNull() && !data.IPv6Ranges.IsUnknown()) ||
		(!data.IsTrusted.IsNull() && !data.IsTrusted.IsUnknown())

	var requestBody graphmodels.NamedLocationable

	switch {
	case isIPLocation:
		ipLocation := graphmodels.NewIpNamedLocation()

		// Set display name
		convert.FrameworkToGraphString(data.DisplayName, ipLocation.SetDisplayName)

		// Set isTrusted
		convert.FrameworkToGraphBool(data.IsTrusted, ipLocation.SetIsTrusted)

		// Construct IP ranges
		ipRanges, err := constructIPRanges(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("failed to construct IP ranges: %w", err)
		}

		if len(ipRanges) > 0 {
			ipLocation.SetIpRanges(ipRanges)
		}

		requestBody = ipLocation

	default:
		// Country Named Location (validation ensures exactly one type is configured)
		countryLocation := graphmodels.NewCountryNamedLocation()

		// Set display name
		convert.FrameworkToGraphString(data.DisplayName, countryLocation.SetDisplayName)

		// Set country lookup method
		convert.FrameworkToGraphEnum(
			data.CountryLookupMethod,
			graphmodels.ParseCountryLookupMethodType,
			countryLocation.SetCountryLookupMethod,
		)

		// Set include unknown countries and regions
		convert.FrameworkToGraphBool(
			data.IncludeUnknownCountriesAndRegions,
			countryLocation.SetIncludeUnknownCountriesAndRegions,
		)

		// Set countries and regions
		if err := convert.FrameworkToGraphStringSet(ctx, data.CountriesAndRegions, func(values []string) {
			if len(values) > 0 {
				countryLocation.SetCountriesAndRegions(values)
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert countries and regions: %w", err)
		}

		requestBody = countryLocation
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructIPRanges builds the ipRanges array from IPv4 and IPv6 ranges
func constructIPRanges(ctx context.Context, data *NamedLocationResourceModel) ([]graphmodels.IpRangeable, error) {
	var ipRanges []graphmodels.IpRangeable

	// Add IPv4 ranges
	if err := convert.FrameworkToGraphStringSet(ctx, data.IPv4Ranges, func(values []string) {
		for _, cidr := range values {
			ipv4Range := graphmodels.NewIPv4CidrRange()
			ipv4Range.SetCidrAddress(&cidr)
			ipRanges = append(ipRanges, ipv4Range)
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert IPv4 ranges: %w", err)
	}

	// Add IPv6 ranges
	if err := convert.FrameworkToGraphStringSet(ctx, data.IPv6Ranges, func(values []string) {
		for _, cidr := range values {
			ipv6Range := graphmodels.NewIPv6CidrRange()
			ipv6Range.SetCidrAddress(&cidr)
			ipRanges = append(ipRanges, ipv6Range)
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert IPv6 ranges: %w", err)
	}

	return ipRanges, nil
}
