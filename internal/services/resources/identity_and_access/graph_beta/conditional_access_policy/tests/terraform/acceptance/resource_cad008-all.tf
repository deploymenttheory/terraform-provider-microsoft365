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
}

resource "microsoft365_graph_beta_groups_group" "cad008_exclude" {
  display_name     = "EID_UA_CAD008_EXCLUDE"
  mail_nickname    = "eid-ua-cad008-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAD008_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAD008: Sign-in Frequency for Browser on Non-Compliant Devices
# Sets sign-in frequency to 1 hour for all apps accessed via browser on non-compliant devices.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad008_browser_signin_frequency" {
  display_name = "acc-test-cad008-all: Session set Sign-in Frequency for All users when Browser and Non-Compliant ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["browser"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad008_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications                            = ["All"]
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

  session_controls = {
    sign_in_frequency = {
      value               = 1
      type                = "days"
      authentication_type = "primaryAndSecondaryAuthentication"
      frequency_interval  = "timeBased"
      is_enabled          = true
    }
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
  }
}

