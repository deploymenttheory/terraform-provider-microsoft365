resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy" "test" {
  name           = "tf-api-probe-parent"
  default_action = "allow"
}

resource "microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy_rule" "test" {
  threat_intelligence_policy_id = microsoft365_graph_beta_identity_and_access_network_threat_intelligence_policy.test.id
  name                          = "tf-api-probe-rule"
  description                   = ""
  action                        = "allow"
  priority                      = 65001
  enabled                       = true
  severity                      = "high"
  destinations = [{
    type   = "fqdn"
    values = ["example.org", "Example.COM"]
  }]
}
