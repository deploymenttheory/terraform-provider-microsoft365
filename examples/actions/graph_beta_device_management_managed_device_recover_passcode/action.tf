# Example 1: Recover passcode for a single iOS device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "single_device" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Recover passcodes for multiple supervised iOS devices
action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "multiple_devices" {
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

# Example 3: Recover passcodes with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_with_validation" {
  config {
    device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Recover passcodes for supervised iOS devices using data source
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ios" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_supervised_ios" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ios.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 5: Recover passcodes for supervised iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "supervised_ipads" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iPadOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_classroom_ipads" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.supervised_ipads.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "15m"
    }
  }
}

# Example 6: Recover passcode for specific user's supervised iOS device
data "microsoft365_graph_beta_device_management_managed_device" "user_ios_device" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'user@example.com') and (operatingSystem eq 'iOS') and (isSupervised eq true)"
}

action "microsoft365_graph_beta_device_management_managed_device_recover_passcode" "recover_user_device" {
  config {
    device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.user_ios_device.items : device.id]

    timeouts = {
      invoke = "5m"
    }
  }
}
