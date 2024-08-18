resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "example" {
  display_name                      = "new filter"
  description                       = "This is an example assignment filter"
  platform                          = "iOS" 
  rule                              = "(device.manufacturer -eq \"thing\")"
  assignment_filter_management_type = "devices"

  role_scope_tags = [8,9]

  timeouts = {
    create = "10s"
    read   = "10s"
    update = "10s"
    delete = "10s"
  }
}

# Output: Display information from the data source
output "existing_filter_id" {
  value = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.existing_filter.id
}

output "existing_filter_description" {
  value = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.existing_filter.description
}

output "existing_filter_platform" {
  value = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.existing_filter.platform
}

output "existing_filter_rule" {
  value = data.microsoft365_graph_beta_device_and_app_management_assignment_filter.existing_filter.rule
}