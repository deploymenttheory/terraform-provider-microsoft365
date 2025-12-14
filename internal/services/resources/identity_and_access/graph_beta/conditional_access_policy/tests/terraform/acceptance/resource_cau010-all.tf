# ==============================================================================
# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

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

resource "microsoft365_graph_beta_groups_group" "cau010_exclude" {
  display_name     = "EID_UA_CAU010_EXCLUDE"
  mail_nickname    = "eid-ua-cau010-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAU010_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAU010: Require Terms of Use
# Requires acceptance of terms of use for all users.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau010_terms_of_use" {
  display_name = "acc-test-cau010-all: Grant Require ToU for All Users when Browser and Modern Auth Clients ${random_string.suffix.result}"
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
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    # include your terms of use ID here
    terms_of_use = [microsoft365_graph_identity_and_access_conditional_access_terms_of_use.consent_on_every_device.id]
  }
}

