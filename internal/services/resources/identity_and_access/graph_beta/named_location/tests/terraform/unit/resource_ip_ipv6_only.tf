resource "microsoft365_graph_beta_identity_and_access_named_location" "ip_ipv6_only" {
  display_name = "unit-test-ip-named-location-ipv6-only"
  is_trusted   = true

  ipv6_ranges = [
    "2001:db8::/32",
    "fe80::/10"
  ]
}