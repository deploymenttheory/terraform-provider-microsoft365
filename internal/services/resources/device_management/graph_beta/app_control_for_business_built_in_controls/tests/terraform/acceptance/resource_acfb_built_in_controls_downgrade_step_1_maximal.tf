resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "downgrade" {
  name        = "acc-test-app-control-downgrade-${random_string.test_suffix.result}"
  description = "Downgrade test - Step 1: Maximal configuration"

  enable_app_control = "audit"
  role_scope_tag_ids = ["0", "1", "2"]

  additional_rules_for_trusting_apps = [
    "trust_apps_with_good_reputation",
    "trust_apps_from_managed_installers"
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
