resource "random_uuid" "acc_test_named_location_ip_maximal" {}

resource "microsoft365_graph_beta_identity_and_access_named_location" "ip_maximal" {
  display_name = "acc-test-named-location-ip-maximal-${random_uuid.acc_test_named_location_ip_maximal.result}"
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