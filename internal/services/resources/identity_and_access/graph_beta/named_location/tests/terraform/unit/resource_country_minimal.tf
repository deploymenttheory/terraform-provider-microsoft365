
resource "microsoft365_graph_beta_identity_and_access_named_location" "country_minimal" {
  display_name                          = "unit-test-country-named-location-minimal"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false

  countries_and_regions = [
    "US"
  ]
}