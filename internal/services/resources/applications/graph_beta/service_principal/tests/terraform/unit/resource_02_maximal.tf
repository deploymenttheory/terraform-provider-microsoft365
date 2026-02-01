# Service Principal configuration for unit testing - Maximal
resource "microsoft365_graph_beta_applications_service_principal" "test_maximal" {
  app_id                        = "22222222-2222-2222-2222-222222222222"
  account_enabled               = true
  app_role_assignment_required  = true
  tags                          = ["HideApp", "WindowsAzureActiveDirectoryIntegratedApp"]
}
