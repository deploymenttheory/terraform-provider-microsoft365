# Security defaults enabled. Only used against mocks: enabling security defaults in a real tenant
# enforces MFA registration for every user and conflicts with enabled Conditional Access policies.
resource "microsoft365_graph_beta_identity_and_access_security_defaults_enforcement_policy" "test" {
  is_enabled = true
}
