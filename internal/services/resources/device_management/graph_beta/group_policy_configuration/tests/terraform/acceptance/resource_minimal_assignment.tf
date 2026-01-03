resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "minimal_assignment" {
  display_name = "AccTest-MinAssign-GPC-${random_string.suffix.result}"

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    }
  ]
}

