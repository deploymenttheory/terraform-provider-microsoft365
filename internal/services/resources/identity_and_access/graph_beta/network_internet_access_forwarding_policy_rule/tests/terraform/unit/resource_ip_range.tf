resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "test" {
  forwarding_policy_id = "11111111-2222-3333-4444-555555555555"
  name                 = "unit-test-forwarding-rule"
  action               = "forward"
  rule_type              = "ip_range"
  client_fallback_action = "block"
  ports                  = ["443"]
  protocol               = "tcp"

  destinations = [
    {
      type          = "ip_range"
      begin_address = "192.0.2.10"
      end_address   = "192.0.2.20"
    },
  ]
}
