# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Test Application and Service Principal
# ==============================================================================

# Create a test application
resource "microsoft365_graph_beta_applications_application" "test" {
  display_name = "acc-test-cau014-sp-${random_string.suffix.result}"
  description  = "Test application for CAU014 conditional access policy"
  hard_delete  = true
}

# Wait for application to be fully created
resource "time_sleep" "wait_for_app" {
  depends_on      = [microsoft365_graph_beta_applications_application.test]
  create_duration = "15s"
}

# Create service principal from the test application
resource "microsoft365_graph_beta_applications_service_principal" "test" {
  app_id = microsoft365_graph_beta_applications_application.test.app_id

  depends_on = [time_sleep.wait_for_app]
}

# Wait for service principal to be fully created
resource "time_sleep" "wait_for_sp" {
  depends_on      = [microsoft365_graph_beta_applications_service_principal.test]
  create_duration = "15s"
}

# ==============================================================================
# Conditional Access Policy
# ==============================================================================

# CAU014: Block Managed Identity for Medium/High Sign-in Risk
# Blocks managed identity (service principal) access when sign-in risk is medium or high.
resource "microsoft365_graph_beta_identity_and_access_conditional_access_policy" "cau014_block_managed_identity_risk" {
  depends_on = [time_sleep.wait_for_sp]

  display_name = "acc-test-cau014-all: Block Managed Identity when Sign in Risk is Medium or High ${random_string.suffix.result}"
  state        = "enabledForReportingButNotEnforced"

  conditions = {
    user_risk_levels              = []
    sign_in_risk_levels           = []
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
        microsoft365_graph_beta_applications_service_principal.test.id
      ]
      exclude_service_principals = []
    }


  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
    terms_of_use                  = []
    authentication_strength       = null
  }

  timeouts = {
    create = "150s"
    read   = "150s"
    update = "150s"
    delete = "150s"
  }
}