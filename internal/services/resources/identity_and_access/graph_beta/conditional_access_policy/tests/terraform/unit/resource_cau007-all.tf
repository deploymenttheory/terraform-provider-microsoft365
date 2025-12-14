# CAU007: Password Change for Medium/High User Risk
# Requires password change for medium and high user risk.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau007_user_risk_password_change" {
  display_name = "CAU007-All: Grant access for Medium and High Risk Users for All Users when Browser and Modern Auth Clients require PWD reset-v1.3"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]
    user_risk_levels = ["high", "medium"]

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

      exclude_guests_or_external_users = {
        guest_or_external_user_types = ["internalGuest", "b2bCollaborationGuest", "b2bCollaborationMember", "b2bDirectConnectUser", "otherExternalUser", "serviceProvider"]
        external_tenants = {
          membership_kind = "all"
        }
      }
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "AND"
    built_in_controls             = ["mfa", "passwordChange"]
    custom_authentication_factors = []
  }

  session_controls = {
    sign_in_frequency = {
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "everyTime"
      is_enabled          = true
    }
  }
}

