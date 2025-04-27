resource "microsoft365_graph_beta_device_and_app_management_windows_quality_update_policy" "quality_update_policy_example" {
  display_name       = "Windows Quality Update Policy"
  description        = "Monthly quality updates for Windows devices"
  hotpatch_enabled   = true
  role_scope_tag_ids = [
    "9", "8"
  ]
  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}