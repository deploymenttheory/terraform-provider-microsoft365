resource "microsoft365_graph_beta_device_management_assignment_filter" "{{.Platform}}" {
  display_name                      = "Test {{.Platform}} Assignment Filter"
  platform                          = "{{.Platform}}"
  rule                              = "{{.Rule}}"
  assignment_filter_management_type = "{{.ManagementType}}"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}