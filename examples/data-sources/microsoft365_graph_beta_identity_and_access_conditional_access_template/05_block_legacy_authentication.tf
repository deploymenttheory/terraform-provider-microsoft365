# Test 05: Create a conditional access policy from template - 
# Block legacy authentication

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

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "block_legacy_auth" {
  name = "Block legacy authentication"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_block_legacy_auth" {
  display_name = "acc-test-ca-policy-template-block-legacy-authentication-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_legacy_auth.details.grant_controls.terms_of_use
  }
}
