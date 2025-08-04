resource "microsoft365_graph_beta_device_management_assignment_filter" "minimal" {
  display_name                      = "Test Minimal Assignment Filter - Unique"
  platform                          = "windows10AndLater"
  rule                              = "(device.osVersion -startsWith \"10.0\")"
  assignment_filter_management_type = "devices"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}