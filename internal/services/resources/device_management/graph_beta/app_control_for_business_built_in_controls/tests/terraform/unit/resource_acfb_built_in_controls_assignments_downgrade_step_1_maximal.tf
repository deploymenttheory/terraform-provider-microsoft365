resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "assignments_downgrade" {
  name        = "unit-test-app-control-assignments-downgrade"
  description = "Assignments downgrade test - Step 1: Maximal assignments"

  enable_app_control = "audit"
  role_scope_tag_ids = ["0"]

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
