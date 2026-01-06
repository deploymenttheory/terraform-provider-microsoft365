# Example 1: Enable lost mode for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_single_lost_device" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "This device has been lost"
        phone_number = "+1234567890"
      }
    ]
  }
}

# Example 2: Enable lost mode for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_multiple_lost_devices" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "Lost iPhone - Please call John at IT to return"
        phone_number = "+1-555-0123"
        footer       = "Reward available for return"
      },
      {
        device_id    = "87654321-4321-4321-4321-ba9876543210"
        message      = "Lost iPad - Contact Mary in HR to return"
        phone_number = "+1-555-0456"
        footer       = "Property of Contoso"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Maximal configuration with validation
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_maximal" {
  config {
    managed_devices = [
      {
        device_id    = "12345678-1234-1234-1234-123456789abc"
        message      = "This device has been lost"
        phone_number = "+1234567890"
        footer       = "Please return to owner"
      }
    ]

    comanaged_devices = [
      {
        device_id    = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        message      = "Lost device"
        phone_number = "+0987654321"
      }
    ]

    ignore_partial_failures = false
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Enable lost mode for supervised iOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_for_supervised_ios" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items : {
        device_id    = device.id
        message      = "This iOS device has been lost. Please contact IT."
        phone_number = "+1-555-0100"
        footer       = "Company Property"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Enable lost mode for user's devices
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'user@example.com') and ((operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS'))"
}

action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_user_lost_devices" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : {
        device_id    = device.id
        message      = format("Lost device belonging to %s", device.userDisplayName)
        phone_number = "+1-555-0200"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Enable lost mode for co-managed device
action "microsoft365_graph_beta_device_management_managed_device_enable_lost_mode" "enable_comanaged_lost" {
  config {
    comanaged_devices = [
      {
        device_id    = "abcdef12-3456-7890-abcd-ef1234567890"
        message      = "Lost co-managed device"
        phone_number = "+1-555-0300"
      }
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}