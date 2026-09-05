resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy" "example" {
  name           = "Example TLS inspection policy"
  description    = "Initial description"
  default_action = "inspect"
}
resource "microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy_rule" "example" {
  tls_inspection_policy_id = microsoft365_graph_beta_identity_and_access_network_tls_inspection_policy.example.id
  name                     = "Inspect selected destinations"
  description              = "Initial description"
  action                   = "inspect"
  priority                 = 1000
  enabled                  = true
  destinations             = [{ type = "fqdn", values = ["example.com", "*.example.org"] }]
}
