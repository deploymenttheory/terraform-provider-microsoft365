resource "microsoft365_graph_beta_identity_and_access_named_location" "ipv6_only" {
  display_name = "example ipv6 named location"
  is_trusted   = true

  ipv6_ranges = [
    "2001:db8::/32",
    "fe80::/10"
  ]
}

resource "microsoft365_graph_beta_identity_and_access_named_location" "ipv4_only" {
  display_name = "example ipv4 named location"
  is_trusted   = false

  ipv4_ranges = [
    "192.168.1.0/24"
  ]
}

resource "microsoft365_graph_beta_identity_and_access_named_location" "ip_ranges" {
  display_name = "example ip ranges named location"
  is_trusted   = true

  ipv4_ranges = [
    "192.168.0.0/16",
    "172.16.0.0/12"
  ]

  ipv6_ranges = [
    "2001:db8::/32",
    "fe80::/10",
    "2001:4860:4860::/48"
  ]
}

resource "microsoft365_graph_beta_identity_and_access_named_location" "country_client_ip" {
  display_name                          = "example country client ip named location"
  country_lookup_method                 = "clientIpAddress"
  include_unknown_countries_and_regions = false

  countries_and_regions = [
    "US",
    "CA",
    "GB"
  ]
}

resource "microsoft365_graph_beta_identity_and_access_named_location" "country_authenticator_gps" {
  display_name                          = "example country authenticator gps named location"
  country_lookup_method                 = "authenticatorAppGps"
  include_unknown_countries_and_regions = true

  countries_and_regions = [
    "AD",
    "AO",
    "AI",
    "AQ"
  ]
}