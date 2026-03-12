# Inbound trust — configures which claims from external tenants are trusted
# when evaluating Conditional Access policies for inbound B2B users.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  inbound_trust = {
    is_mfa_accepted                           = true
    is_compliant_device_accepted              = true
    is_hybrid_azure_ad_joined_device_accepted = true
  }
}
