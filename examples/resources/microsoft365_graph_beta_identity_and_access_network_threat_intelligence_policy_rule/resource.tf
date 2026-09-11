resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy" "example" {
  name           = "Example threat intelligence policy"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy_rule" "example" {
  threat_intelligence_policy_id = microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy.example.id
  name                          = "Example threat intelligence rule"
  description                   = "Initial description"
  action                        = "allow"
  priority                      = 1000
  enabled                       = true
  severity                      = "high"
  destinations = [
    {
      type   = "fqdn"
      values = ["example.com", "example.org"]
    }
  ]
}
