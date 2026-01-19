# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Group Dependencies
# ==============================================================================

# Break Glass Emergency Access Accounts
resource "microsoft365_graph_beta_groups_group" "breakglass" {
  display_name     = "EID_UA_ConAcc-Breakglass"
  mail_nickname    = "eid-ua-conacc-breakglass"
  mail_enabled     = false
  security_enabled = true
  description      = "Group containing Break Glass Accounts"
  hard_delete      = true
}

resource "microsoft365_graph_beta_groups_group" "cad004_exclude" {
  display_name     = "EID_UA_CAD004_EXCLUDE"
  mail_nickname    = "eid-ua-cad004-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "exclusion group for CA policy CAD004_EXCLUDE"
  hard_delete      = true
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD004: Require MFA on Non-Compliant Devices via Browser
# Requires MFA for all users accessing Office 365 via browser when device is non-compliant.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad004_browser_noncompliant_mfa" {
  display_name = "acc-test-cad004-o365: Grant Require MFA for All users when Browser and Non-Compliant ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad004_exclude.id
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

