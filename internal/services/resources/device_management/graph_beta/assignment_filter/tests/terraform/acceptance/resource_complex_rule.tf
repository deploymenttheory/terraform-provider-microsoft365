resource "microsoft365_graph_beta_device_management_assignment_filter" "complex" {
  display_name                        = "Test Complex Rule Assignment Filter"
  platform                           = "windows10AndLater"
  rule                               = <<-EOT
    (device.osVersion -startsWith "10.0") and 
    (device.manufacturer -eq "Microsoft Corporation") and 
    (device.model -notContains "Virtual") and
    (device.enrollmentProfileName -contains "Corporate")
  EOT
  assignment_filter_management_type  = "devices"

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}