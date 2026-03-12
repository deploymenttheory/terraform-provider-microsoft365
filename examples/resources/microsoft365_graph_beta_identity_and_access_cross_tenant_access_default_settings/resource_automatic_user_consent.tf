# Automatic user consent — controls whether users in this tenant can
# automatically consent to cross-tenant access without admin approval.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  automatic_user_consent_settings = {
    inbound_allowed  = false
    outbound_allowed = false
  }
}
