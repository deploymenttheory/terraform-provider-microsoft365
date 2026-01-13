resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "enforce_mode" {
  name        = "unit-test-app-control-enforce-mode"
  description = "unit-test-app-control-enforce-mode"

  enable_app_control = "enforce"
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
