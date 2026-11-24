resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "test" {
  forwarding_policy_id = "11111111-2222-3333-4444-555555555555"
  name                 = "unit-test-forwarding-rule"
  action               = "forward"
  rule_type = "ip_subnet"
  ports     = ["443"]
  protocol  = "tcp"

  destinations = [
    {
      type  = "ip_subnet"
      value = "192.0.2.0/24"
    },
    {
      type  = "ip_subnet"
      value = "198.51.100.0/24"
    },
  ]
}
