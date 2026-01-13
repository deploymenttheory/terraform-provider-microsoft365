resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "lifecycle" {
  name        = "unit-test-app-control-lifecycle"
  description = "Lifecycle test - Step 1: Minimal configuration"
  
  enable_app_control = "audit"
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
