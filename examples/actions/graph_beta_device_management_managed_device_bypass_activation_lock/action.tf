# Example 1: Bypass Activation Lock for a single device
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_single_device" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 2: Bypass Activation Lock for multiple devices (batch processing)
action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_batch" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Bypass Activation Lock for supervised iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS') and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_supervised_ios" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items : device.id]

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 4: Bypass Activation Lock for DEP-enrolled macOS devices
data "microsoft365_graph_beta_device_management_managed_device" "dep_macos" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'macOS') and (deviceEnrollmentType eq 'deviceEnrollmentProgram')"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_dep_macos" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.dep_macos.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Bypass Activation Lock for departing employee's Apple devices
data "microsoft365_graph_beta_device_management_managed_device" "departing_user_apple_devices" {
  filter_type  = "odata"
  odata_filter = "userId eq 'user@example.com' and (operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS' or operatingSystem eq 'macOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_departing_user" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_user_apple_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 6: Bypass Activation Lock for corporate-owned Apple devices
data "microsoft365_graph_beta_device_management_managed_device" "corporate_apple" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS' or operatingSystem eq 'iPadOS' or operatingSystem eq 'macOS') and managedDeviceOwnerType eq 'company'"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_corporate_apple" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.corporate_apple.items : device.id]

    timeouts = {
      invoke = "20m"
    }
  }
}

# Example 7: Bypass Activation Lock for devices with specific model (e.g., iPhone 13)
data "microsoft365_graph_beta_device_management_managed_device" "iphone_13_devices" {
  filter_type  = "odata"
  odata_filter = "model eq 'iPhone 13' and isSupervised eq true"
}

action "microsoft365_graph_beta_device_management_managed_device_bypass_activation_lock" "bypass_iphone_13" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.iphone_13_devices.items : device.id]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "bypassed_device_count" {
  value       = length(action.bypass_batch.config.device_ids)
  description = "Number of devices that received Activation Lock bypass command"
}

output "bypassed_corporate_count" {
  value       = length(action.bypass_corporate_apple.config.device_ids)
  description = "Number of corporate Apple devices with bypass codes generated"
}
