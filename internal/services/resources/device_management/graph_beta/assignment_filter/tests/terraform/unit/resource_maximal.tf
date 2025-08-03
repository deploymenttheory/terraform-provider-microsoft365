resource "microsoft365_graph_beta_device_management_assignment_filter" "maximal" {
  display_name                        = "Test Maximal Assignment Filter - Unique"
  description                         = "Maximal assignment filter for testing with all features"
  platform                           = "windows10AndLater"
  rule                               = <<-EOT
    (device.osVersion -startsWith "10.0") and 
    (device.manufacturer -eq "Microsoft Corporation") and 
    (device.model -notContains "Virtual")
  EOT
  assignment_filter_management_type  = "devices"
  role_scope_tags                    = ["0", "1"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}