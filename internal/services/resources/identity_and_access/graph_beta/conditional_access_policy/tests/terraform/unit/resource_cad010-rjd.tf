# CAD010: Require MFA for Device Registration/Join
# Requires MFA when users register or join devices to Entra ID.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad010_device_registration_mfa" {
  display_name = "CAD010-RJD: Require MFA for device join or registration when Browser and Modern Auth Clients-v1.1"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = []
      exclude_applications                            = []
      include_user_actions                            = ["urn:user:registerdevice"]
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["mfa"]
    custom_authentication_factors = []
  }
}

