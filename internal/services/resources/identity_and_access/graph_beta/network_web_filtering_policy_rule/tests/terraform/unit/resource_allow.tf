resource "microsoft365_graph_beta_identity_and_access_network_web_filtering_policy_rule" "test" {
  web_filtering_policy_id = "11111111-2222-3333-4444-555555555555"
  name                    = "unit-test-allow-rule"
  description             = "Allow rule without headers"
  priority                = 300
  action                  = "allow"
  status                  = "enabled"
  web_categories          = ["Gambling", "Hacking"]
}
