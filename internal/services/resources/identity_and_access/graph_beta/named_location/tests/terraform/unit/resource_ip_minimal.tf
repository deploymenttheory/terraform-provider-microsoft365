resource "microsoft365_graph_beta_identity_and_access_named_location" "ip_minimal" {
  display_name = "unit-test-ip-named-location-minimal"
  is_trusted   = false

  ipv4_ranges = [
    "192.168.1.0/24"
  ]
}