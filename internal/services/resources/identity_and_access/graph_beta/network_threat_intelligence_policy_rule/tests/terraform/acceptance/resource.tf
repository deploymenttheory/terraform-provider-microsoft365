resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy" "test" {
  name           = "tf-api-probe-parent"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy_rule" "test" {
  threat_intelligence_policy_id = microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy.test.id
  name                          = "tf-api-probe-rule"
  description                   = "Initial description"
  action                        = "allow"
  priority                      = 1000
  enabled                       = true
  severity                      = "high"
  destinations = [{
    type   = "fqdn"
    values = ["Example.COM", "example.org", "Example.COM"]
  }]
}
