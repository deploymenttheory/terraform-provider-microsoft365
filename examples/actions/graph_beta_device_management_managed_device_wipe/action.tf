# Example 1: Wipe a single device (factory reset, removes all data) - Minimal
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_single" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 2: Wipe multiple devices with validation and failure handling
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_batch" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    validate_device_exists  = true # Validate devices before wiping
    ignore_partial_failures = true # Continue if some wipes fail

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Selective wipe - keep user data, remove only company data
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_company_data_only" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    keep_user_data = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Wipe with enrollment data preserved (for automatic re-enrollment)
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_keep_enrollment" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    keep_enrollment_data = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 5: Wipe macOS device with Activation Lock
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_macos" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    macos_unlock_code = "123456" # 6-digit PIN for Activation Lock bypass

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 6: Wipe macOS with obliteration behavior control
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_macos_always_obliterate" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    obliteration_behavior = "always" # Always obliterate on T2+ Macs

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 7: Wipe Windows device with protected wipe (preserves UEFI licenses)
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_windows" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    use_protected_wipe = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 8: Wipe devices with eSIM, preserving data plan
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_keep_esim" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    persist_esim_data_plan = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 9: Comprehensive wipe with multiple options and robust validation
action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_comprehensive" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]

    keep_enrollment_data    = true
    keep_user_data          = true
    persist_esim_data_plan  = true
    obliteration_behavior   = "doNotObliterate"
    validate_device_exists  = true
    ignore_partial_failures = false

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 10: Wipe non-compliant devices from data source
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_non_compliant_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_devices.items : device.id]

    # Wipe but keep enrollment data for automatic re-enrollment after compliance
    keep_enrollment_data    = true
    validate_device_exists  = true
    ignore_partial_failures = true # Continue even if some wipes fail

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 11: Wipe old devices by OS version
data "microsoft365_graph_beta_device_management_managed_device" "old_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows' and osVersion startsWith '10.0'"
}

action "microsoft365_graph_beta_device_management_managed_device_wipe" "wipe_old_windows_devices" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.old_devices.items : device.id]

    use_protected_wipe = true # Preserve Windows licenses

    timeouts = {
      invoke = "20m"
    }
  }
}

# Output examples
output "wiped_device_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_wipe.wipe_batch.config.device_ids)
  description = "Number of devices wiped in batch operation"
}

output "non_compliant_devices_to_wipe" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_wipe.wipe_non_compliant_devices.config.device_ids)
  description = "Number of non-compliant devices being wiped"
}

