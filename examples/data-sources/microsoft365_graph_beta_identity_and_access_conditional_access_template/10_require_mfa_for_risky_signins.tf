# Test 10: Create a conditional access policy from template - 
# Require multifactor authentication for risky sign-ins

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

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "mfa_risky_signins" {
  name = "Require multifactor authentication for risky sign-ins"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_mfa_risky_signins" {
  display_name = "acc-test-ca-policy-template-require-mfa-for-risky-signins-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.include_users
      exclude_users  = []
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.mfa_risky_signins.details.grant_controls.terms_of_use
  }
}
