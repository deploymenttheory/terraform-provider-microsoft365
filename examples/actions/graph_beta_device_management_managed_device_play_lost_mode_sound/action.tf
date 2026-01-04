# Example 1: Play lost mode sound for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_single" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
      }
    ]
  }
}

# Example 2: Play lost mode sound for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_multiple" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
      },
      {
        device_id = "87654321-4321-4321-4321-ba9876543210"
      },
      {
        device_id = "abcdef12-3456-7890-abcd-ef1234567890"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_maximal" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
      },
      {
        device_id = "87654321-4321-4321-4321-ba9876543210"
      }
    ]

    comanaged_devices = [
      {
        device_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Play sound for all devices in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "devices_in_lost_mode" {
  filter_type  = "odata"
  odata_filter = "lostModeState eq 'enabled'"
}

action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_all_lost_mode" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.devices_in_lost_mode.items : {
        device_id = device.id
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Play sound for user's devices in lost mode
data "microsoft365_graph_beta_device_management_managed_device" "user_lost_devices" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'user@example.com') and (lostModeState eq 'enabled')"
}

action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_user_devices" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.user_lost_devices.items : {
        device_id = device.id
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 6: Play sound for co-managed device
action "microsoft365_graph_beta_device_management_managed_device_play_lost_mode_sound" "play_sound_comanaged" {
  config {
    comanaged_devices = [
      {
        device_id = "abcdef12-3456-7890-abcd-ef1234567890"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

}