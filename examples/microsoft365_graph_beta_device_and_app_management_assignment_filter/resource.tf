resource "microsoft365_device_and_app_management_assignment_filter" "example" {
  display_name                    = "Example Assignment Filter"
  description                     = "This is an example assignment filter"
  platform                        = "android"
  rule                            = "device.os == 'Android'"
  assignment_filter_management_type = "devices"

  role_scope_tags = ["tag1", "tag2"]

  payloads {
    payload_id               = "payload1"
    payload_type             = "type1"
    group_id                 = "group1"
    assignment_filter_type   = "include"
  }

  payloads {
    payload_id               = "payload2"
    payload_type             = "type2"
    group_id                 = "group2"
    assignment_filter_type   = "exclude"
  }

  timeouts {
    create = "30m"
    read   = "30m"
    update = "30m"
    delete = "30m"
  }
}
