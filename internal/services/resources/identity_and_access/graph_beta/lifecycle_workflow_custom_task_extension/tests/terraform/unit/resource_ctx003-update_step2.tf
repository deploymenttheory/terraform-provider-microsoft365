# CTX003 step 2: Every property changed and a callback added
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension" "ctx003_update" {
  display_name = "unit-test-ctx003-update-renamed"
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
