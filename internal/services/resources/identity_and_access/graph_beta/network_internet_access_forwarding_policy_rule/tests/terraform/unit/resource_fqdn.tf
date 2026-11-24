resource "microsoft365_graph_beta_identity_and_access_network_internet_access_forwarding_policy_rule" "test" {
  forwarding_policy_id = "11111111-2222-3333-4444-555555555555"
  name                 = "unit-test-forwarding-rule"
  action               = "forward"
  rule_type              = "fqdn"
  client_fallback_action = "block"
  ports                  = ["80", "443"]
  protocol               = "tcp"

  destinations = [
    {
      type  = "fqdn"
      value = "www.example.com"
    },
  ]
}
