resource "random_uuid" "acc_test_named_location_ip_minimal" {}

resource "microsoft365_graph_beta_identity_and_access_named_location" "ip_minimal" {
  display_name = "acc-test-named-location-ip-minimal-${random_uuid.acc_test_named_location_ip_minimal.result}"
  is_trusted   = false
  
  ipv4_ranges = [
    "192.168.1.0/24"
  ]
}