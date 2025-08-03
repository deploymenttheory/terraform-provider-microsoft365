resource "microsoft365_graph_beta_device_management_assignment_filter" "{{.Platform}}" {
  display_name                        = "Test {{.Platform}} Assignment Filter"
  platform                           = "{{.Platform}}"
  {{- if or (eq .Platform "androidMobileApplicationManagement") (eq .Platform "iOSMobileApplicationManagement") (eq .Platform "windowsMobileApplicationManagement") }}
  rule                               = "(app.osVersion -eq \"10.0\")"
  assignment_filter_management_type  = "apps"
  {{- else }}
  rule                               = "(device.osVersion -eq \"10.0\")"
  assignment_filter_management_type  = "devices"
  {{- end }}

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}