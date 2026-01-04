resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "random_uuid" "test_group" {}

data "azuread_group" "test_exclusion" {
  display_name     = "Test Group"
  security_enabled = true
}

resource "microsoft365_graph_beta_device_management_group_policy_configuration" "maximal_assignment" {
  display_name       = "AccTest-MaxAssign-GPC-${random_string.suffix.result}"
  description        = "Configuration with comprehensive assignments for acceptance testing"
  role_scope_tag_ids = ["0"]

  assignments = [
    {
      type = "allDevicesAssignmentTarget"
    },
    {
      type = "allLicensedUsersAssignmentTarget"
    }
  ]
}

