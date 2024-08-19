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

