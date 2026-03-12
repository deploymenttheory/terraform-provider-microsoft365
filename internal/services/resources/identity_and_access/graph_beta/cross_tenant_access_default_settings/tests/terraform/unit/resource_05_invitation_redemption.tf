resource "microsoft365_graph_beta_identity_and_access_cross_tenant_access_default_settings" "test" {
  restore_defaults_on_destroy = true

  b2b_collaboration_outbound = {
    users_and_groups = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllUsers"
          target_type = "user"
        }
      ]
    }
    applications = {
      access_type = "allowed"
      targets = [
        {
          target      = "AllApplications"
          target_type = "application"
        }
      ]
    }
  }

  invitation_redemption_identity_provider_configuration = {
    primary_identity_provider_precedence_order = [
      "azureActiveDirectory",
      "externalFederation",
      "socialIdentityProviders"
    ]
    fallback_identity_provider = "emailOneTimePasscode"
  }
}
