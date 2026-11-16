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

# CTX002: Maximal custom task extension with token authentication, client settings and a callback
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension" "ctx002_maximal" {
  display_name = "acc-test-ctx002-maximal-${random_string.suffix.result}"
  description  = "Calls a Logic App and waits for it to resume the workflow"

  endpoint_configuration = {
    subscription_id         = "11111111-1111-1111-1111-111111111111"
    resource_group_name     = "acc-test-rg-placeholder-2"
    logic_app_workflow_name = "acc-test-la-placeholder-2"
    url                     = "https://prod-01.westeurope.logic.azure.com:443/workflows/11111111111111111111111111111111/triggers/manual/paths/invoke?api-version=2016-10-01"
  }

  authentication_configuration = {
    type        = "azureAdTokenAuthentication"
    resource_id = "00000003-0000-0000-c000-000000000000"
  }

  client_configuration = {
    maximum_retries         = 0
    timeout_in_milliseconds = 2000
  }

  callback_configuration = {
    timeout_duration   = "PT1H"
    authorized_app_ids = ["00000003-0000-0000-c000-000000000000"]
  }
}
