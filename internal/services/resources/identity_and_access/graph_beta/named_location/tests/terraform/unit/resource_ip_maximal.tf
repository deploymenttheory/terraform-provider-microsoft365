resource "microsoft365_graph_beta_identity_and_access_named_location" "ip_maximal" {
  display_name = "unit-test-ip-named-location-maximal"
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