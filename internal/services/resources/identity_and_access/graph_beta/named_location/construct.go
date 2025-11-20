package graphBetaNamedLocation

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource converts the Terraform resource model to a plain map for JSON marshaling
// Returns a map[string]any that can be directly JSON marshaled by the HTTP client
func constructResource(ctx context.Context, data *NamedLocationResourceModel) (map[string]any, error) {

	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, data); err != nil {
		return nil, err
	}

	requestBody := make(map[string]any)

	convert.FrameworkToGraphString(data.DisplayName, func(val *string) {
		if val != nil {
			requestBody["displayName"] = *val
		}
	})

	// Determine the type of named location based on which fields are populated
	isIPLocation := (!data.IPv4Ranges.IsNull() && !data.IPv4Ranges.IsUnknown()) ||
		(!data.IPv6Ranges.IsNull() && !data.IPv6Ranges.IsUnknown()) ||
		(!data.IsTrusted.IsNull() && !data.IsTrusted.IsUnknown())

	switch {
	case isIPLocation:
		requestBody["@odata.type"] = "#microsoft.graph.ipNamedLocation"

		convert.FrameworkToGraphBool(data.IsTrusted, func(val *bool) {
			if val != nil {
				requestBody["isTrusted"] = *val
			}
		})

		ipRanges, err := constructIPRanges(ctx, data)
		if err != nil {
			return nil, fmt.Errorf("failed to construct IP ranges: %w", err)
		}

		if len(ipRanges) > 0 {
			requestBody["ipRanges"] = ipRanges
		}

	default:
		// Country Named Location (validation ensures exactly one type is configured)
		requestBody["@odata.type"] = "#microsoft.graph.countryNamedLocation"

		convert.FrameworkToGraphString(data.CountryLookupMethod, func(val *string) {
			if val != nil {
				requestBody["countryLookupMethod"] = *val
			}
		})

		convert.FrameworkToGraphBool(data.IncludeUnknownCountriesAndRegions, func(val *bool) {
			if val != nil {
				requestBody["includeUnknownCountriesAndRegions"] = *val
			}
		})

		if err := convert.FrameworkToGraphStringSet(ctx, data.CountriesAndRegions, func(values []string) {
			if len(values) > 0 {
				requestBody["countriesAndRegions"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert countries and regions: %w", err)
		}
	}

	if debugJSON, err := json.MarshalIndent(requestBody, "", "    "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]any{
			"json": "\n" + string(debugJSON),
		})
	} else {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructIPRanges builds the ipRanges array from IPv4 and IPv6 ranges
func constructIPRanges(ctx context.Context, data *NamedLocationResourceModel) ([]map[string]any, error) {
	var ipRanges []map[string]any

	if err := convert.FrameworkToGraphStringSet(ctx, data.IPv4Ranges, func(values []string) {
		for _, cidr := range values {
			ipRange := map[string]any{
				"@odata.type": "#microsoft.graph.iPv4CidrRange",
				"cidrAddress": cidr,
			}
			ipRanges = append(ipRanges, ipRange)
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert IPv4 ranges: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.IPv6Ranges, func(values []string) {
		for _, cidr := range values {
			ipRange := map[string]any{
				"@odata.type": "#microsoft.graph.iPv6CidrRange",
				"cidrAddress": cidr,
			}
			ipRanges = append(ipRanges, ipRange)
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert IPv6 ranges: %w", err)
	}

	return ipRanges, nil
}
