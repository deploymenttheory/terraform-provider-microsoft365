# Test 16: Create a conditional access policy from template - 
# Require phishing-resistant multifactor authentication for admins

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

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "phishing_resistant_mfa_admins" {
  name = "Require phishing-resistant multifactor authentication for admins"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_phishing_resistant_mfa_admins" {
  display_name = "acc-test-ca-policy-template-require-phishing-resistant-mfa-for-admins-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.service_principal_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.conditions.users.exclude_roles
    }
  }

  grant_controls = {
    operator                      = "OR" # Template specifies AND, but API normalizes to OR when only authentication_strength is used
    built_in_controls             = []
    custom_authentication_factors = []
    authentication_strength = {
      id = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.phishing_resistant_mfa_admins.details.grant_controls.authentication_strength.id
    }
  }
}
