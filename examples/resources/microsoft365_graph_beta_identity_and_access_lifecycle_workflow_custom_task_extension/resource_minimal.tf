# Fire-and-forget custom task extension: the workflow calls the Logic App and continues without waiting.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension" "minimal" {
  display_name = "Notify HR system"

  endpoint_configuration = {
    subscription_id         = "00000000-0000-0000-0000-000000000000"
    resource_group_name     = "rg-identity-governance"
    logic_app_workflow_name = "la-notify-hr"
    # Required when the provider authenticates as an application (client credentials).
    url = "https://prod-00.westeurope.logic.azure.com:443/workflows/00000000000000000000000000000000/triggers/manual/paths/invoke?api-version=2016-10-01"
  }

  authentication_configuration = {
    type = "azureAdPopTokenAuthentication"
  }
}
