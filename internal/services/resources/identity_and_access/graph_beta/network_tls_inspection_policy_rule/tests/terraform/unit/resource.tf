resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy" "test" {
  name           = "tf-api-probe-policy"
  description    = "Initial description"
  default_action = "inspect"
}
resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy_rule" "test" {
  tls_inspection_policy_id = microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy.test.id
  name                     = "tf-api-probe-rule"
  description              = "Initial description"
  action                   = "inspect"
  priority                 = 1000
  enabled                  = true
  destinations             = [{ type = "fqdn", values = ["example.com", "*.example.org"] }]
}
