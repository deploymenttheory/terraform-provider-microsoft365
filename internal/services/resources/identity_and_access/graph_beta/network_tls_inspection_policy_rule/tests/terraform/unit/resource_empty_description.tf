resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy" "test" {
  description = ""
  name        = "tf-api-probe-policy-updated"

  default_action = "bypass"
}
resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy_rule" "test" {
  tls_inspection_policy_id = microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy.test.id
  description              = ""
  name                     = "tf-api-probe-rule-updated"

  action       = "bypass"
  priority     = 65001
  enabled      = false
  destinations = [{ type = "fqdn", values = ["example.net"] }, { type = "web_category", values = ["Finance", "HealthAndMedicine"] }]
}
