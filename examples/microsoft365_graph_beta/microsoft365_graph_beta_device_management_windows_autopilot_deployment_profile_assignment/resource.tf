# Assignment to all devices
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile_assignment" "all_devices" {
  windows_autopilot_deployment_profile_id = microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.entra_joined.id
  source = "direct"
  
  target = {
    target_type = "allDevices"
  }

# Optional timeouts block
  timeouts = {
    create = "3m"
    read   = "3m"
    update = "3m"
    delete = "3m"
  }
}

# Assignment to a specific Entra ID group
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile_assignment" "corporate_devices_1" {
  windows_autopilot_deployment_profile_id = microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.entra_joined.id
  source = "direct"
  
  target = {
    target_type = "groupAssignment"
    group_id = "11111111-2222-3333-4444-555555555555"
  }
}

# Exclusion assignment
resource "microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile_assignment" "exclude_group" {
  windows_autopilot_deployment_profile_id = microsoft365_graph_beta_device_management_windows_autopilot_deployment_profile.entra_joined.id
  source = "direct"
  
  target = {
    target_type = "exclusionGroupAssignment"
    group_id = "11111111-2222-3333-4444-555555555555"
  }
}