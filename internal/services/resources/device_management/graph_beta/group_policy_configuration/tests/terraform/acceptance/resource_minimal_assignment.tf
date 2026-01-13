resource "random_string" "test_suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "minimal_assignment" {
  display_name = "acc-test-003-minimal-assignment-${random_string.test_suffix.result}"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}
