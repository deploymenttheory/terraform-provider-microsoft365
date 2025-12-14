# CAU001A: Require MFA for Guests - Windows Azure AD
# Requires MFA for guest/external users accessing Windows Azure Active Directory.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau001a_guest_mfa_azure_ad" {
  display_name = "CAU001A-Windows Azure Active Directory: Grant Require MFA for guests when Browser and Modern Auth Clients-v1.0"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
      ]
      include_roles = []
      exclude_roles = []

      include_guests_or_external_users = {
        guest_or_external_user_types = ["b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "internalGuest", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["66666666-6666-6666-6666-666666666666"] # Windows Azure Active Directory
      exclude_applications                            = []
      include_user_actions                            = []
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

