resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_app_control_for_business_built_in_controls" "assignments_lifecycle" {
  name        = "acc-test-app-control-assignments-lifecycle-${random_string.test_suffix.result}"
  description = "Assignments lifecycle test - Step 1: Minimal assignments"

  enable_app_control = "audit"
  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]

  timeouts = {
    create = "15m"
    read   = "5m"
    update = "15m"
    delete = "10m"
  }
}
