
resource "microsoft365_graph_beta_identity_and_access_named_location" "country_authenticator_gps" {
  display_name                           = "unit-test-country-named-location-authenticator-gps"
  country_lookup_method                  = "authenticatorAppGps"
  include_unknown_countries_and_regions = true
  
  countries_and_regions = [
    "AD",
    "AO", 
    "AI",
    "AQ"
  ]
}