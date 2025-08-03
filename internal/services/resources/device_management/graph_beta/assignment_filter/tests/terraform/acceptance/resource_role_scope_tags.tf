resource "microsoft365_graph_beta_device_management_assignment_filter" "role_tags" {
  display_name                      = "Test Role Scope Tags Assignment Filter"
  platform                          = "windows10AndLater"
  rule                              = "(device.osVersion -startsWith \"10.0\")"
  assignment_filter_management_type = "devices"
  role_scope_tags                   = ["0", "1", "2"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}