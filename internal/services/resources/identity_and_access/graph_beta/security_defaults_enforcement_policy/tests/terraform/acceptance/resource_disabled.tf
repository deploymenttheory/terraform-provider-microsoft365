# Security defaults disabled. The acceptance tests only apply the tenant's current value (false)
# and never enable security defaults in the shared test tenant.
resource "microsoft365_graph_beta_identity_and_access_security_defaults_enforcement_policy" "test" {
  is_enabled = false
}
