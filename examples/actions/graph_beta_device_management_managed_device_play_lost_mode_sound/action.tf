# Example 1: Play lost mode sound on a single device with default duration
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_single" {

  managed_devices {
    device_id = "12345678-1234-1234-1234-123456789abc"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Play lost mode sound with specific duration
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_with_duration" {

  managed_devices {
    device_id           = "12345678-1234-1234-1234-123456789abc"
    duration_in_minutes = "5"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 3: Play sound on multiple devices with different durations
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_multiple" {

  managed_devices {
    device_id           = "12345678-1234-1234-1234-123456789abc"
    duration_in_minutes = "3"
  }

  managed_devices {
    device_id           = "87654321-4321-4321-4321-ba9876543210"
    duration_in_minutes = "10"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 4: Play sound on all devices currently in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "devices_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "lostModeState ne 'disabled' and operatingSystem eq 'iOS'"
}

action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_all_lost_mode" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.devices_in_lost_mode.items
    content {
      device_id           = managed_devices.value.id
      duration_in_minutes = "5"
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 5: Play sound for specific user's devices in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "user_lost_devices" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and lostModeState ne 'disabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_user_devices" {

  dynamic "managed_devices" {
    for_each = data.microsoft365_graph_beta_device_management_managed_device.user_lost_devices.items
    content {
      device_id           = managed_devices.value.id
      duration_in_minutes = "3"
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 6: Play sound on co-managed device
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_comanaged" {

  comanaged_devices {
    device_id           = "abcdef12-3456-7890-abcd-ef1234567890"
    duration_in_minutes = "5"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Play sound to locate device nearby
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "locate_nearby_device" {

  managed_devices {
    device_id           = "12345678-1234-1234-1234-123456789abc"
    duration_in_minutes = "2"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Output examples
output "devices_with_sound" {
  value       = length(action.play_sound_multiple.managed_devices)
  description = "Number of devices that had lost mode sound played"
}

output "lost_mode_devices_count" {
  value       = length(action.play_sound_all_lost_mode.managed_devices)
  description = "Number of devices in lost mode that received sound command"
}