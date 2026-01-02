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
  display_name     = "EID_UA_ConAcc-Breakglass-${random_string.suffix.result}"
  mail_nickname    = "eid-ua-conacc-breakglass-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Group containing Break Glass Accounts"
}

resource "microsoft365_graph_beta_groups_group" "cad013_exclude" {
  display_name     = "EID_UA_CAD013_Exclude-${random_string.suffix.result}"
  mail_nickname    = "eid-ua-cad013-exclude-${random_string.suffix.result}"
  mail_enabled     = false
  security_enabled = true
  description      = "Exclusion group for CA policy CAD013"
}

# ==============================================================================
# CAD013: Selected Apps - Compliant Device Requirement
# ==============================================================================

# Requires compliant device for access to selected applications.
# Using well-known Microsoft service App IDs:
# - Azure DevOps: 499b84ac-1321-427f-aa17-267ca6975798
# - Office 365 Management: c5393580-f805-4401-95e8-94b7a6ef2fc2
# - Dynamics 365 Business Central: 996def3d-b36c-4153-8607-a6fd3c01b89f
# - Windows Azure Service Management API: 797f4846-ba00-4fd7-ba43-dac1f8f63013

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cad013_selected_apps_compliant" {
  display_name = "acc-test-cad013-selected: Grant access for All users when Browser and Modern Auth Clients and Compliant ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"


  conditions = {
    client_app_types = ["browser", "mobileAppsAndDesktopClients"]

    users = {
      include_users  = ["All"]
      exclude_users  = []
      include_groups = []
      exclude_groups = [
        microsoft365_graph_beta_groups_group.breakglass.id,
        microsoft365_graph_beta_groups_group.cad013_exclude.id
      ]
      include_roles = []
      exclude_roles = []
    }

    applications = {
      include_applications = [
        "c5393580-f805-4401-95e8-94b7a6ef2fc2", # Office 365 Management APIs
        "499b84ac-1321-427f-aa17-267ca6975798", # Azure DevOps
        "996def3d-b36c-4153-8607-a6fd3c01b89f", # Dynamics 365 Business Central
        "797f4846-ba00-4fd7-ba43-dac1f8f63013"  # Windows Azure Service Management API
      ]

      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    platforms = {
      include_platforms = ["all"]
      exclude_platforms = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator          = "OR"
    built_in_controls = ["compliantDevice", "domainJoinedDevice"]

    custom_authentication_factors = []
  }

  timeouts = {
    create = "150s"
    read   = "150s"
    update = "150s"
    delete = "150s"
  }
}

