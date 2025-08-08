# Assignment Filters used for acceptance testing
# These assignment filters serve as dependencies.

resource "random_string" "assignment_filter_suffix" {
  length  = 8
  special = false
  upper   = false
}


resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_assignment_filter_1" {
  display_name                      = "acc-test-assignment-filter-1-${random_string.assignment_filter_suffix.result}"
  description                       = "Updated description for acceptance testing"
  platform                          = "windows10AndLater"
  rule                              = <<-EOT
    (device.osVersion -startsWith "10.0") and 
    (device.manufacturer -eq "Microsoft Corporation") and 
    (device.model -notContains "Virtual")
  EOT
  assignment_filter_management_type = "devices"
  role_scope_tags                   = ["0", "1"]

  timeouts = {
    create = "180s"
    read   = "180s"
    update = "180s"
    delete = "180s"
  }
}