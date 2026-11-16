# ==============================================================================
# Random Suffix for Unique Resource Names
# ==============================================================================

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

# ==============================================================================
# Lifecycle Workflow Custom Task Extension
# ==============================================================================

# CTX003 step 1: Minimal custom task extension that later steps update in place
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension" "ctx003_update" {
  display_name = "acc-test-ctx003-update-${random_string.suffix.result}"

  endpoint_configuration = {
    subscription_id         = "00000000-0000-0000-0000-000000000000"
    resource_group_name     = "acc-test-rg-placeholder"
    logic_app_workflow_name = "acc-test-la-placeholder"
    url                     = "https://prod-00.westeurope.logic.azure.com:443/workflows/00000000000000000000000000000000/triggers/manual/paths/invoke?api-version=2016-10-01"
  }

  authentication_configuration = {
    type = "azureAdPopTokenAuthentication"
  }
}
