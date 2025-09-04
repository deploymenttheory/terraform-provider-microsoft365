resource "random_uuid" "acc_test_named_location_country_authenticator_gps" {}

resource "microsoft365_graph_beta_identity_and_access_named_location" "country_authenticator_gps" {
  display_name                          = "acc-test-named-location-country-authenticator-gps-${random_uuid.acc_test_named_location_country_authenticator_gps.result}"
  country_lookup_method                 = "authenticatorAppGps"
  include_unknown_countries_and_regions = true

  countries_and_regions = [
    "AD",
    "AO",
    "AI",
    "AQ"
  ]
}