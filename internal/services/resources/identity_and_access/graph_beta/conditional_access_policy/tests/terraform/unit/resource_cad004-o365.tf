# CAD004: Require MFA on Non-Compliant Devices via Browser
# Requires MFA for all users accessing Office 365 via browser when device is non-compliant.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad004_browser_noncompliant_mfa" {
  display_name = "CAD004-O365: Grant Require MFA for All users when Browser and Non-Compliant-v1.3"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

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
      include_applications                            = ["Office365"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    devices = {
      device_filter = {
        mode = "exclude"
        rule = "device.isCompliant -eq True -or device.trustType -eq \"ServerAD\""
      }
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    # Note: Source uses custom authentication strength ID "eaedd457-3e01-413b-a02e-417489193d1d" (Custom MFA)
    # Replace with your tenant-specific custom authentication strength ID or use built-in:
    authentication_strength = {
      id = "00000000-0000-0000-0000-000000000002" # Multifactor authentication (Changed from custom MFA ID)
    }
  }
}

