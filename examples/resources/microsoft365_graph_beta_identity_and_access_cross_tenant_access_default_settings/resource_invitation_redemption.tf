# Invitation redemption — controls the order in which identity providers are
# tried when a B2B guest redeems an invitation, and the fallback provider if
# none of the primary options succeed.

resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "example" {
  restore_defaults_on_destroy = true

  invitation_redemption_identity_provider_configuration = {
    primary_identity_provider_precedence_order = [
      "azureActiveDirectory",
      "externalFederation",
      "socialIdentityProviders"
    ]
    fallback_identity_provider = "emailOneTimePasscode"
  }
}
