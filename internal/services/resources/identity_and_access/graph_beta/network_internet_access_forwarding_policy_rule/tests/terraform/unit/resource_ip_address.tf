resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "test" {
  forwarding_policy_id = "11111111-2222-3333-4444-555555555555"
  name                 = "unit-test-forwarding-rule"
  action               = "forward"
  rule_type = "ip_address"
  ports     = ["53"]
  protocol  = "udp"

  destinations = [
    {
      type  = "ip_address"
      value = "192.0.2.53"
    },
  ]
}
