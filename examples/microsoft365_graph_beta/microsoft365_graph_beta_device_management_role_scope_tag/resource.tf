# Example of creating a basic role scope tag with a group assignment
resource "microsoft365_graph_beta_device_management_role_scope_tag" "helpdesk" {
  display_name = "Helpdesk Support Tag"
  description  = "Role scope tag for helpdesk support staff"

  assignments = ["00000000-0000-0000-0000-000000000001"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}

# Example of creating multiple related role scope tags with assignments
resource "microsoft365_graph_beta_device_management_role_scope_tag" "it_support" {
  display_name = "IT Support Tag"
  description  = "Role scope tag for IT support teams"

  assignments = ["00000000-0000-0000-0000-000000000002"]
}

resource "microsoft365_graph_beta_device_management_role_scope_tag" "device_management" {
  display_name = "Device Management Tag"
  description  = "Role scope tag for device management teams"

  assignments = [
    "00000000-0000-0000-0000-000000000003",
    "00000000-0000-0000-0000-000000000004"
  ]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}