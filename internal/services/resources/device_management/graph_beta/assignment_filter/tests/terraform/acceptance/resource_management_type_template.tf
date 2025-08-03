resource "microsoft365_graph_beta_device_management_assignment_filter" "{{.ManagementType}}" {
  display_name                        = "Test {{.ManagementType}} Management Assignment Filter"
  {{- if eq .ManagementType "apps" }}
  platform                           = "androidMobileApplicationManagement"
  rule                               = "(app.osVersion -startsWith \"10.0\")"
  {{- else }}
  platform                           = "windows10AndLater"
  rule                               = "(device.osVersion -startsWith \"10.0\")"
  {{- end }}
  assignment_filter_management_type  = "{{.ManagementType}}"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}