resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "audit_mode" {
  name        = "acc-test-app-control-audit-mode-${random_string.test_suffix.result}"
  description = "acc-test-app-control-audit-mode"
  
  enable_app_control = "audit"
  role_scope_tag_ids = ["0"]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
