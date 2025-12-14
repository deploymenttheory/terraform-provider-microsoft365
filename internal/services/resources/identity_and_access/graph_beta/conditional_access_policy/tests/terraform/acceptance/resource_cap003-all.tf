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

resource "microsoft365_graph_beta_groups_group" "cap003_exclude" {
  display_name     = "EID_UA_CAP003_EXCLUDE"
  mail_nickname    = "eid-ua-cap003-exclude"
  mail_enabled     = false
  security_enabled = true
  description      = "uexcludeion group for CA policy CAP003_EXCLUDE"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================


# CAP003: Block Device Code Flow
# Blocks device code authentication flow to prevent phishing attacks.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cap003_block_device_code_flow" {
  display_name = "acc-test-cap003-all: Block device code authentication flow ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types = ["all"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cap003_exclude.id
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

    authentication_flows = {
      transfer_methods = "deviceCodeFlow"
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }
}

