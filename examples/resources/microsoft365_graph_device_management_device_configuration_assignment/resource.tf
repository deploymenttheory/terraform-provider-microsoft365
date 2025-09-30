resource "microsoft365_graph_device_management_device_configuration_assignment" "example" {
  device_configuration_id = "00000000-0000-0000-0000-000000000000"
  target_type             = "groupAssignment"
  group_id                = "00000000-0000-0000-0000-000000000000"
}

resource "microsoft365_graph_device_management_device_configuration_assignment" "all_devices_example" {
  device_configuration_id = "00000000-0000-0000-0000-000000000000"
  target_type             = "allDevices"
}

resource "microsoft365_graph_device_management_device_configuration_assignment" "filtered_example" {
  device_configuration_id = "00000000-0000-0000-0000-000000000000"
  target_type             = "groupAssignment"
  group_id                = "00000000-0000-0000-0000-000000000000"
} 