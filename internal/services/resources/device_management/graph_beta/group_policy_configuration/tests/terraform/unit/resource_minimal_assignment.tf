resource "microsoft365_graph_beta_device_management_group_policy_configuration" "minimal_assignment" {
  display_name = "unit-test-003-minimal-assignment"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}
