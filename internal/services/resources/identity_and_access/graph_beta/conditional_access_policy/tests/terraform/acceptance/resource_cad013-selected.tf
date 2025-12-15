# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Azure AD Applications
# ==============================================================================

resource "azuread_application" "cad013_app_01" {
  display_name = "acc-test-cad013-app-01-${random_string.suffix.result}"
}

resource "azuread_application" "cad013_app_02" {
  display_name = "acc-test-cad013-app-02-${random_string.suffix.result}"
}

# ==============================================================================
# CAD013: Selected Apps - Compliant Device Requirement
# ==============================================================================

# Requires compliant device for access to selected applications.
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
        azuread_application.cad013_app_01.application_id,
        azuread_application.cad013_app_02.application_id
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
    operator                      = "OR"
    built_in_controls             = ["compliantDevice", "domainJoinedDevice"]
    custom_authentication_factors = []
  }

  timeouts = {
    create = "150s"
    read   = "150s"
    update = "150s"
    delete = "150s"
  }
}

