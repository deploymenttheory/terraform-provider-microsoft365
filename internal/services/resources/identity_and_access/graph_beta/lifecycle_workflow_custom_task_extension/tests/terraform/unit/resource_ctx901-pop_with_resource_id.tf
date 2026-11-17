# CTX901: azureAdPopTokenAuthentication with resource_id
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension" "ctx901_invalid" {
  display_name = "unit-test-ctx901-invalid"

  endpoint_configuration = {
    subscription_id         = "00000000-0000-0000-0000-000000000000"
    resource_group_name     = "acc-test-rg-placeholder"
    logic_app_workflow_name = "acc-test-la-placeholder"
    url                     = "https://prod-00.westeurope.logic.azure.com:443/workflows/00000000000000000000000000000000/triggers/manual/paths/invoke?api-version=2016-10-01"
  }

  authentication_configuration = {
    type        = "azureAdPopTokenAuthentication"
    resource_id = "00000003-0000-0000-c000-000000000000"
  }
}
