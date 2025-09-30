# Example of creating a basic role scope tag with a group assignment
resource "microsoft365_graph_beta_device_management_role_scope_tag" "helpdesk" {
  display_name = "Helpdesk Support Tag"
  description  = "Role scope tag for helpdesk support staff"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "device_management" {
  display_name = "Device Management Tag"
  description  = "Role scope tag for device management teams"

  assignments = [
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000001"
    },
    {
      type     = "groupAssignmentTarget"
      group_id = "00000000-0000-0000-0000-000000000002"
    },
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}