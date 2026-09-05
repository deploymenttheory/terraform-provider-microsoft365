resource "random_id" "suffix" { byte_length = 4 }
resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy" "test" {
  name = "tf-api-probe-policy-${random_id.suffix.hex}-updated"

  default_action = "bypass"
}
resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy_rule" "test" {
  tls_inspection_policy_id = microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy.test.id
  name                     = "tf-api-probe-rule-${random_id.suffix.hex}-updated"

  action       = "bypass"
  priority     = 65001
  enabled      = false
  destinations = [{ type = "fqdn", values = ["example.net"] }, { type = "web_category", values = ["Finance", "HealthAndMedicine"] }]
}
