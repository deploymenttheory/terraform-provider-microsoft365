resource "microsoft365_graph_beta_device_management_group_policy_configuration" "minimal_assignment" {
  display_name = "Minimal Assignment Group Policy Configuration"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}

