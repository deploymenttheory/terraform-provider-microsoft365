# Security defaults disabled.
resource "microsoft365_graph_beta_identity_and_access_security_defaults_enforcement_policy" "test" {
  is_enabled = false
}
