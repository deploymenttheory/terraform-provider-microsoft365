# CAD016: Token Protection for EXO/SPO/CloudPC on Windows
# Requires token protection for Exchange Online, SharePoint Online, and Cloud PC on Windows.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad016_token_protection_windows" {
  display_name = "CAD016-EXO_SPO_CloudPC: Require token protection when Modern Auth Clients on Windows-v1.2"
  state        = "enabledForReportingButNotEnforced"
  hard_delete  = true

  conditions = {
    client_app_types = ["mobileAppsAndDesktopClients"]

    users = {
      include_users  = []
      exclude_users  = []
      include_groups = ["77777777-7777-7777-7777-777777777777"]
      exclude_groups = [
        "22222222-2222-2222-2222-222222222222",
        "33333333-3333-3333-3333-333333333333"
      ]
      include_roles = []
      exclude_roles = []

      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications = [
        "9cdead84-a844-4324-93f2-b2e6bb768d07",
        "0af06dc6-e4b5-4f28-818e-e78e62d137a5",
        "270efc09-cd0d-444b-a71f-39af4910ec45",
        "00000002-0000-0ff1-ce00-000000000000",
        "00000003-0000-0ff1-ce00-000000000000",
      ]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["windows"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

