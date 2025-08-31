resource "random_uuid" "acc_test_named_location_ipv6_only" {}

resource "microsoft365_graph_beta_identity_and_access_named_location" "ipv6_only" {
  display_name = "acc-test-named-location-ipv6-only-${random_uuid.acc_test_named_location_ipv6_only.result}"
  is_trusted   = true

  ipv6_ranges = [
    "2001:db8::/32",
    "fe80::/10"
  ]
}