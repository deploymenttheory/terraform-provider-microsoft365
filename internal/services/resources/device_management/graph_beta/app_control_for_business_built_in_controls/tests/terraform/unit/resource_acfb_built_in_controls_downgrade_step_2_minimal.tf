resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "downgrade" {
  name        = "unit-test-app-control-downgrade"
  description = "Downgrade test - Step 2: Minimal configuration"

  enable_app_control = "audit"
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
