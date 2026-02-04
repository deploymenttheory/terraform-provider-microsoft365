# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Service Principal Dependencies
# ==============================================================================

# Use a built-in Microsoft service principal for testing
data "microsoft365_graph_beta_applications_service_principal" "windows_azure_service_management_api" {
  display_name = "Windows Azure Service Management API"
}

# Wait for service principal to propagate
resource "time_sleep" "wait_for_sp" {
  depends_on = [data.microsoft365_graph_beta_applications_service_principal.windows_azure_service_management_api]

  create_duration = "10s"
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
    user_risk_levels = []
    sign_in_risk_levels = []
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
        data.microsoft365_graph_beta_applications_service_principal.windows_azure_service_management_api.id
      ]
      exclude_service_principals = []
    }

    
  }

  grant_controls = {
    operator                      = "OR"
    built_in_controls             = ["block"]
    custom_authentication_factors = []
    terms_of_use = []
    authentication_strength = null
  }

  timeouts = {
    create = "150s"
    read   = "150s"
    update = "150s"
    delete = "150s"
  }
}