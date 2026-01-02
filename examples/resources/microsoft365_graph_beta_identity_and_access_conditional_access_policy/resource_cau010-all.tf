# CAU010: Require Terms of Use
# Requires acceptance of terms of use for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau010_terms_of_use" {
  display_name = "CAU010-All: Grant Require ToU for All Users when Browser and Modern Auth Clients-v1.2"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cau010_exclude.id
      ]
      include_roles = []
      exclude_roles = []

      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications = ["All"]
      exclude_applications = [
        "0000000a-0000-0000-c000-000000000000", # Microsoft Intune
        "d4ebce55-015a-49b5-a083-c84d1797ae8c"  # Microsoft Intune Enrollment
      ]
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    # include your terms of use ID here
    terms_of_use = [microsoft365_graph_identity_and_access_conditional_access_terms_of_use.consent_on_every_device.id]
  }
}

