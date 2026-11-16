# Custom task extension that waits up to one hour for the Logic App to resume the workflow.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_custom_task_extension" "maximal" {
  display_name = "Provision mailbox and wait for completion"
  description  = "Calls the provisioning Logic App and waits for it to report back"

  endpoint_configuration = {
    subscription_id         = "00000000-0000-0000-0000-000000000000"
    resource_group_name     = "rg-identity-governance"
    logic_app_workflow_name = "la-provision-mailbox"
    url                     = "https://prod-00.westeurope.logic.azure.com:443/workflows/00000000000000000000000000000000/triggers/manual/paths/invoke?api-version=2016-10-01"
  }

  authentication_configuration = {
    type        = "azureAdTokenAuthentication"
    resource_id = "11111111-1111-1111-1111-111111111111"
  }

  client_configuration = {
    maximum_retries         = 0
    timeout_in_milliseconds = 2000
  }

  callback_configuration = {
    timeout_duration = "PT1H"
    # Client ID of the Logic App's managed identity, allowed to resume the task.
    authorized_app_ids = ["22222222-2222-2222-2222-222222222222"]
  }

  timeouts = {
    create = "5m"
    read   = "5m"
    update = "5m"
    delete = "5m"
  }
}
