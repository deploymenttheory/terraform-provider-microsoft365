# ============================================================================
# Example 1: Quick scan on managed devices
# ============================================================================
# Use case: Routine security check on selected Windows devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "quick_scan_managed" {

  managed_devices = [
    {
      device_id  = "12345678-1234-1234-1234-123456789abc"
      quick_scan = true
    },
    {
      device_id  = "87654321-4321-4321-4321-ba9876543210"
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 2: Full scan on specific device
# ============================================================================
# Use case: Comprehensive scan after security incident
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "full_scan_incident" {

  managed_devices = [
    {
      device_id  = "12345678-1234-1234-1234-123456789abc"
      quick_scan = false # Full comprehensive scan
    }
  ]

  timeouts = {
    invoke = "5m"
  }
}

# ============================================================================
# Example 3: Mixed scan types on different devices
# ============================================================================
# Use case: Quick scan for most, full scan for suspicious devices
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "mixed_scan_types" {

  managed_devices = [
    {
      device_id  = "12345678-1234-1234-1234-123456789abc"
      quick_scan = true # Routine check
    },
    {
      device_id  = "87654321-4321-4321-4321-ba9876543210"
      quick_scan = false # Suspected malware - full scan
    },
    {
      device_id  = "abcdef12-3456-7890-abcd-ef1234567890"
      quick_scan = true # Routine check
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 4: Scan co-managed devices
# ============================================================================
# Use case: Scan Windows devices managed by both Intune and ConfigMgr
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_comanaged" {

  comanaged_devices = [
    {
      device_id  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      quick_scan = true
    },
    {
      device_id  = "11111111-2222-3333-4444-555555555555"
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 5: Scan both managed and co-managed devices
# ============================================================================
# Use case: Mixed environment with different management types
action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_mixed_management" {

  managed_devices = [
    {
      device_id  = "12345678-1234-1234-1234-123456789abc"
      quick_scan = true
    }
  ]

  comanaged_devices = [
    {
      device_id  = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

# ============================================================================
# Example 6: Quick scan all Windows devices using datasource
# ============================================================================
# Use case: Routine security scan across entire Windows fleet
data "microsoft365_graph_beta_device_management_managed_device" "all_windows" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_all_windows" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows.items : {
      device_id  = device.id
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 7: Full scan on non-compliant Windows devices
# ============================================================================
# Use case: Thorough scan on devices that failed compliance
data "microsoft365_graph_beta_device_management_managed_device" "non_compliant_windows" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (complianceState eq 'noncompliant')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_non_compliant" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.non_compliant_windows.items : {
      device_id  = device.id
      quick_scan = false # Full scan for non-compliant devices
    }
  ]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 8: Scan Windows devices by naming convention
# ============================================================================
# Use case: Scan specific department or location devices
data "microsoft365_graph_beta_device_management_managed_device" "finance_windows" {
  filter_type  = "device_name"
  filter_value = "FIN-WS-"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_finance_dept" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.finance_windows.items : {
      device_id  = device.id
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "15m"
  }
}

# ============================================================================
# Example 9: After-hours full scan on workstations
# ============================================================================
# Use case: Comprehensive scan during off-hours to avoid performance impact
data "microsoft365_graph_beta_device_management_managed_device" "workstations" {
  filter_type  = "device_name"
  filter_value = "WKSTN-"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "after_hours_full_scan" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.workstations.items : {
      device_id  = device.id
      quick_scan = false # Full scan during off-hours
    }
  ]

  timeouts = {
    invoke = "30m"
  }
}

# ============================================================================
# Example 10: Conditional scan based on last scan time
# ============================================================================
# Use case: Full scan on devices that haven't been scanned recently
locals {
  # Devices needing full scan (example logic)
  devices_need_full_scan = [
    "12345678-1234-1234-1234-123456789abc",
    "87654321-4321-4321-4321-ba9876543210"
  ]

  # Devices needing quick scan
  devices_need_quick_scan = [
    "abcdef12-3456-7890-abcd-ef1234567890",
    "fedcba98-7654-3210-fedc-ba9876543210"
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "conditional_scan" {

  managed_devices = concat(
    [
      for device_id in local.devices_need_full_scan : {
        device_id  = device_id
        quick_scan = false
      }
    ],
    [
      for device_id in local.devices_need_quick_scan : {
        device_id  = device_id
        quick_scan = true
      }
    ]
  )

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 11: Emergency threat response scan
# ============================================================================
# Use case: Immediate full scan after threat intel indicates new malware
data "microsoft365_graph_beta_device_management_managed_device" "all_windows_devices" {
  filter_type  = "odata"
  odata_filter = "operatingSystem eq 'Windows'"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "emergency_threat_scan" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.all_windows_devices.items : {
      device_id  = device.id
      quick_scan = false # Full scan for threat response
    }
  ]

  timeouts = {
    invoke = "60m"
  }
}

# ============================================================================
# Example 12: Scan Windows servers only
# ============================================================================
# Use case: Security scan on Windows Server infrastructure
data "microsoft365_graph_beta_device_management_managed_device" "windows_servers" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (contains(deviceName, 'SRV'))"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_servers" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.windows_servers.items : {
      device_id  = device.id
      quick_scan = true # Quick scan for servers to minimize impact
    }
  ]

  timeouts = {
    invoke = "20m"
  }
}

# ============================================================================
# Example 13: Scan by user assignment
# ============================================================================
# Use case: Scan all Windows devices assigned to specific user
data "microsoft365_graph_beta_device_management_managed_device" "user_devices" {
  filter_type  = "odata"
  odata_filter = "(operatingSystem eq 'Windows') and (userPrincipalName eq 'john.doe@company.com')"
}

action "microsoft365_graph_beta_device_management_managed_device_windows_defender_scan" "scan_user_devices" {

  managed_devices = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.user_devices.items : {
      device_id  = device.id
      quick_scan = true
    }
  ]

  timeouts = {
    invoke = "10m"
  }
}

