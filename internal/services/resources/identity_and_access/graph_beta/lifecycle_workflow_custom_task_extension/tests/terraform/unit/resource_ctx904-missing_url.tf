# CTX904: endpoint without url, which Graph rejects for an app-only caller
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension" "ctx904_missing_url" {
  display_name = "unit-test-ctx904-missing-url"

  endpoint_configuration = {
    subscription_id         = "00000000-0000-0000-0000-000000000000"
    resource_group_name     = "acc-test-rg-placeholder"
    logic_app_workflow_name = "acc-test-la-placeholder"
  }

  authentication_configuration = {
    type = "azureAdPopTokenAuthentication"
  }
}
