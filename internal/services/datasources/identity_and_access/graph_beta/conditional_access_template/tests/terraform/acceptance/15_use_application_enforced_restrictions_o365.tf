# Test 15: Create a conditional access policy from template - 
# Use application enforced restrictions for O365 apps

# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Datasource
# ==============================================================================

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "app_enforced_restrictions_o365" {
  name = "Use application enforced restrictions for O365 apps"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_app_enforced_restrictions_o365" {
  display_name = "acc-test-ca-policy-template-use-application-enforced-restrictions-for-o365-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.conditions.users.exclude_roles
    }
  }

  # Template has null grant_controls, but we need to provide an empty structure to avoid inconsistency
  grant_controls = {
    operator                      = "OR"
    built_in_controls             = []
    custom_authentication_factors = []
    terms_of_use                  = []
  }

  session_controls = {
    application_enforced_restrictions = {
      is_enabled = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.app_enforced_restrictions_o365.details.session_controls.application_enforced_restrictions.is_enabled
    }
  }
}
