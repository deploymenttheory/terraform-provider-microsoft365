# Example 1: Revoke Apple VPP licenses from a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_single" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc"
    ]
  }
}

# Example 2: Revoke Apple VPP licenses from multiple devices
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_multiple" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210",
      "abcdef12-3456-7890-abcd-ef1234567890"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Revoke Apple VPP licenses with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_with_validation" {
  config {
    managed_device_ids = [
      "12345678-1234-1234-1234-123456789abc",
      "87654321-4321-4321-4321-ba9876543210"
    ]

    comanaged_device_ids = [
      "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Revoke VPP licenses from departing user's devices
data "microsoft365_graph_beta_device_management_managed_device" "departing_user_ios" {
  filter_type  = "odata"
  odata_filter = "(userPrincipalName eq 'departing.user@example.com') and ((operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS'))"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_departing_user" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.departing_user_ios.items : device.id]

    validate_device_exists = true

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 5: Revoke VPP licenses from all iOS/iPadOS devices
data "microsoft365_graph_beta_device_management_managed_device" "all_apple_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'iOS') or (operatingSystem eq 'iPadOS')"
}

action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_all_ios" {
  config {
    managed_device_ids = [for device in data.microsoft365_graph_beta_device_management_managed_device.all_apple_devices.items : device.id]

    ignore_partial_failures = true

    timeouts = {
      invoke = "30m"
    }
  }
}

# Example 6: Revoke VPP licenses for co-managed devices
action "microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses" "revoke_comanaged" {
  config {
    comanaged_device_ids = [
      "11111111-1111-1111-1111-111111111111",
      "22222222-2222-2222-2222-222222222222"
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Output examples
output "revoked_vpp_licenses_count" {
  value       = length(action.microsoft365_graph_beta_device_management_managed_device_revoke_apple_vpp_licenses.revoke_multiple.config.managed_device_ids)
  description = "Number of devices that had VPP licenses revoked"
}
