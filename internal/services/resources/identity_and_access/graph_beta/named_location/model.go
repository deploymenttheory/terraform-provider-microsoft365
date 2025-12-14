// REF: https://learn.microsoft.com/en-us/graph/api/resources/namedlocation?view=graph-rest-beta
package graphBetaNamedLocation

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NamedLocationResourceModel represents the schema for the Named Location resource
type NamedLocationResourceModel struct {
	ID               types.String `tfsdk:"id"`
	DisplayName      types.String `tfsdk:"display_name"`
	CreatedDateTime  types.String `tfsdk:"created_date_time"`
	ModifiedDateTime types.String `tfsdk:"modified_date_time"`
	// IP Named Location fields
	IsTrusted  types.Bool `tfsdk:"is_trusted"`
	IPv4Ranges types.Set  `tfsdk:"ipv4_ranges"`
	IPv6Ranges types.Set  `tfsdk:"ipv6_ranges"`
	// Country Named Location fields
	CountryLookupMethod               types.String `tfsdk:"country_lookup_method"`
	IncludeUnknownCountriesAndRegions types.Bool   `tfsdk:"include_unknown_countries_and_regions"`
	CountriesAndRegions               types.Set    `tfsdk:"countries_and_regions"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}
