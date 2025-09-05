# Minimal App Control for Business configuration with multiple assignments and filters
resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "minimal" {
  name        = "acc-test-app-control-for-business-built-in-controls-minimal"
  description = "acc-test-app-control-for-business-built-in-controls-minimal"

  # App Control settings
  enable_app_control = "audit" # audit = logs but allows, enforce = blocks untrusted apps

  # Role scope tags for specific departments
  role_scope_tag_ids = ["0", "1", "2"]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = microsoft365_graph_beta_groups_group.acc_test_group_2.id
    },
    {
      type = "allDevicesAssignmentTarget"
    }
  ]


  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}