resource "microsoft365_graph_beta_device_management_macos_platform_script_assignment" "example" {
  macos_platform_script_id = "00000000-0000-0000-0000-000000000001"
  target = {
    target_type = "allDevices"
  }
}

resource "microsoft365_graph_beta_device_management_macos_platform_script_assignment" "all_users" {
  macos_platform_script_id = "00000000-0000-0000-0000-000000000001"
  target = {
    target_type = "allLicensedUsers"
  }
}

resource "microsoft365_graph_beta_device_management_macos_platform_script_assignment" "group" {
  macos_platform_script_id = "00000000-0000-0000-0000-000000000001"
  target = {
    target_type                                      = "groupAssignment"
    group_id                                         = "11111111-1111-1111-1111-111111111111"
    device_and_app_management_assignment_filter_type = "include"
    device_and_app_management_assignment_filter_id   = "11111111-1111-1111-1111-111111111111"
  }
} 