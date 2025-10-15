
# Create multiple unique assignment filters to avoid dependency issues
resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_assignment_filter_1" {
  display_name                      = "acc-test-assignment-filter-1-${random_string.assignment_filter_suffix.result}"
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
    create = "60s"
    read   = "60s"
    update = "60s"
    delete = "180s"
  }
}

resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_assignment_filter_2" {
  display_name                      = "acc-test-assignment-filter-2-${random_string.assignment_filter_suffix.result}"
  description                       = "Assignment filter for group 2 acceptance testing"
  platform                          = "windows10AndLater"
  rule                              = <<-EOT
    (device.osVersion -startsWith "10.0") and 
    (device.deviceCategory -eq "Corporate")
  EOT
  assignment_filter_management_type = "devices"
  role_scope_tags                   = ["0"]

}

resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_assignment_filter_3" {
  display_name                      = "acc-test-assignment-filter-3-${random_string.assignment_filter_suffix.result}"
  description                       = "Assignment filter for all users acceptance testing"
  platform                          = "windows10AndLater"
  rule                              = "(device.deviceTrustType -eq \"Azure AD joined\")"
  assignment_filter_management_type = "devices"
  role_scope_tags                   = ["0"]

}

resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_assignment_filter_4" {
  display_name                      = "acc-test-assignment-filter-4-${random_string.assignment_filter_suffix.result}"
  description                       = "Assignment filter for all devices acceptance testing"
  platform                          = "windows10AndLater"
  rule                              = "(device.deviceTrustType -ne \"Azure AD registered\")"
  assignment_filter_management_type = "devices"
  role_scope_tags                   = ["0"]

}

resource "microsoft365_graph_beta_device_management_assignment_filter" "acc_test_assignment_filter_5" {
  display_name                      = "acc-test-assignment-filter-5-${random_string.assignment_filter_suffix.result}"
  description                       = "Assignment filter for exclusion group acceptance testing"
  platform                          = "windows10AndLater"
  rule                              = "(device.deviceTrustType -in [\"Hybrid Azure AD joined\",\"Azure AD joined\"])"
  assignment_filter_management_type = "devices"
  role_scope_tags                   = ["0"]

}