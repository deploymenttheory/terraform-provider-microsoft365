# Test 18: Create a conditional access policy from template - 
# Block access to Office365 apps for users with insider risk

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

data "microsoft365_graph_beta_identity_and_access_conditional_access_template" "block_o365_insider_risk" {
  name = "Block access to Office365 apps for users with insider risk"
}

# ==============================================================================
# Validation
# ==============================================================================

resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "from_template_block_o365_insider_risk" {
  display_name = "acc-test-ca-policy-template-block-access-o365-insider-risk-${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.client_app_types
    user_risk_levels              = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.user_risk_levels
    sign_in_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.sign_in_risk_levels
    service_principal_risk_levels = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.service_principal_risk_levels
    insider_risk_levels           = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.insider_risk_levels

    applications = {
      include_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.applications.include_applications
      exclude_applications                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.applications.exclude_applications
      include_user_actions                            = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.applications.include_user_actions
      include_authentication_context_class_references = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.applications.include_authentication_context_class_references
    }

    users = {
      include_users  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.include_users
      exclude_users  = [] # Placeholder string "Current administrator will be excluded". In prod define real users to exclude administrators from this policy.
      include_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.include_groups
      exclude_groups = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.exclude_groups
      include_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.include_roles
      exclude_roles  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.exclude_roles

      exclude_guests_or_external_users = {
        guest_or_external_user_types = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.conditions.users.exclude_guests_or_external_users.guest_or_external_user_types
        external_tenants = {
          membership_kind = "all" # Template has externalTenants: null, but external_tenants is required in the resource
        }
      }
    }
  }

  grant_controls = {
    operator                      = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.grant_controls.operator
    built_in_controls             = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.grant_controls.built_in_controls
    custom_authentication_factors = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.grant_controls.custom_authentication_factors
    terms_of_use                  = data.microsoft365_graph_beta_identity_and_access_conditional_access_template.block_o365_insider_risk.details.grant_controls.terms_of_use
  }
}
