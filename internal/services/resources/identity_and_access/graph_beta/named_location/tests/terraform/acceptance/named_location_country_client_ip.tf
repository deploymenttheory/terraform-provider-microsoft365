resource "random_uuid" "acc_test_named_location_country_client_ip" {}

resource "microsoft365_graph_beta_identity_and_access_named_location" "country_client_ip" {
  display_name                          = "acc-test-named-location-country-client-ip-${random_uuid.acc_test_named_location_country_client_ip.result}"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false

  countries_and_regions = [
    "US",
    "CA",
    "GB"
  ]
}