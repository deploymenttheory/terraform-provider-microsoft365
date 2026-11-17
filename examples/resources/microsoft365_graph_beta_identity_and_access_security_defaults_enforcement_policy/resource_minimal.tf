# Keep Microsoft Entra ID security defaults disabled, for example because the tenant
# is protected with Conditional Access policies instead.
resource "microsoft365_graph_beta_identity_and_access_security_defaults_enforcement_policy" "example" {
  is_enabled = false
}
