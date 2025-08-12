# Assignment Filters for assignments acceptance test
# These have special dependency management to ensure proper destroy order
# Assignment filters must be destroyed AFTER the configuration policy that uses them

resource "random_string" "assignment_filter_suffix" {
  length  = 8
  special = false
  upper   = false
}

# Create multiple unique assignment filters
resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_assignment_filter_1" {
  display_name                      = "acc-test-assignment-filter-settings-catalog-configuration-policy-${random_string.assignment_filter_suffix.result}"
  description                       = "Assignment filter for group 1 acceptance testing"
  platform                          = "windows10AndLater"
  rule                              = <<-EOT
    (device.osVersion -startsWith "10.0") and 
    (device.manufacturer -eq "Microsoft Corporation") and 
    (device.model -notContains "Virtual")
  EOT
  assignment_filter_management_type = "devices"
  role_scope_tags                   = ["0"]

  timeouts = {
    create = "10s"   
    read   = "10s"   
    update = "10s"   
    delete = "300s"  # Very long timeout to allow server-side propagation after policy deletion
  }
}