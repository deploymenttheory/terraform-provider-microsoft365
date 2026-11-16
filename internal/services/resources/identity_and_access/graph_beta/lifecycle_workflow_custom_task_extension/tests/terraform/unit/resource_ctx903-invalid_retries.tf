# CTX903: maximum_retries outside 0..1
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension" "ctx903_invalid" {
  display_name = "unit-test-ctx903-invalid"

  endpoint_configuration = {
    subscription_id         = "00000000-0000-0000-0000-000000000000"
    resource_group_name     = "acc-test-rg-placeholder"
    logic_app_workflow_name = "acc-test-la-placeholder"
    url                     = "https://prod-00.westeurope.logic.azure.com:443/workflows/00000000000000000000000000000000/triggers/manual/paths/invoke?api-version=2016-10-01"
  }

  authentication_configuration = {
    type = "azureAdPopTokenAuthentication"
  }

  client_configuration = {
    maximum_retries = 3
  }
}
