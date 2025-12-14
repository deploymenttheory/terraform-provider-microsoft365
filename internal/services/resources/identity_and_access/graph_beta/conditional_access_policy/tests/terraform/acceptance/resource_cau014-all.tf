# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Application Dependencies
# ==============================================================================

# Test Application for CAU014
resource "azuread_application" "cau014_test_app" {
  display_name = "CAU014-TestApp-${random_string.suffix.result}"
}

# Service Principal for Test Application
resource "azuread_service_principal" "cau014_test_app" {
  client_id = azuread_application.cau014_test_app.client_id
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================

# CAU014: Block Managed Identity for Medium/High Sign-in Risk
# Blocks managed identity (service principal) access when sign-in risk is medium or high.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau014_block_managed_identity_risk" {
  display_name = "acc-test-cau014-all: Block Managed Identity when Sign in Risk is Medium or High ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    client_app_types              = ["all"]
    service_principal_risk_levels = ["high", "medium"]

    users = {
      include_users  = ["None"]
      exclude_users  = []
      include_groups = []
      exclude_groups = []
      include_roles  = []
      exclude_roles  = []
    }

    applications = {
      include_applications                            = ["All"]
      exclude_applications                            = []
      include_user_actions                            = []
      include_authentication_context_class_references = []
    }

    client_applications = {
      include_service_principals = [
        azuread_service_principal.cau014_test_app.object_id
      ]
      exclude_service_principals = []
    }

    sign_in_risk_levels = []
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
  }

  timeouts = {
    create = "150s"
    read   = "150s"
    update = "150s"
    delete = "150s"
  }
}