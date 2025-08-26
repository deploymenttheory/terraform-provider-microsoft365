# Advanced App Control for Business configuration with multiple assignments and filters
resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "advanced" {
   name        = "unit-test-app-control-for-business-built-in-controls-maximal"
  description = "unit-test-app-control-for-business-built-in-controls-maximal"
  
  # App Control settings
  enable_app_control = "audit"  # audit = logs but allows, enforce = blocks untrusted apps
  additional_rules_for_trusting_apps = ["trust_apps_with_good_reputation", "trust_apps_from_managed_installers"]  # Both Microsoft Store and managed installer apps
  
  # Role scope tags for specific departments
  role_scope_tag_ids = ["0", "1", "2"]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type        = "groupAssignmentTarget"
      group_id    = "33333333-3333-3333-3333-333333333333"
      filter_id   = "44444444-4444-4444-4444-444444444444"
      filter_type = "include"
    },
    {
      type        = "allDevicesAssignmentTarget"
      filter_id   = "55555555-5555-5555-5555-555555555555"
      filter_type = "exclude"
    }
  ]
  
  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}