# Example 1: Run remediation script on single device
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "single_device" {
  managed_devices {
    device_id        = "12345678-1234-1234-1234-123456789abc"
    script_policy_id = "87654321-4321-4321-4321-ba9876543210"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 2: Run same remediation on multiple devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "multiple_same_script" {
  managed_devices {
    device_id        = "device1-1234-1234-1234-123456789abc"
    script_policy_id = "disk-cleanup-script-guid-here"
  }

  managed_devices {
    device_id        = "device2-1234-1234-1234-123456789abc"
    script_policy_id = "disk-cleanup-script-guid-here"
  }

  managed_devices {
    device_id        = "device3-1234-1234-1234-123456789abc"
    script_policy_id = "disk-cleanup-script-guid-here"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 3: Run different scripts on different devices
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "different_scripts" {
  managed_devices {
    device_id        = "device1-1234-1234-1234-123456789abc"
    script_policy_id = "disk-cleanup-script-guid"
  }

  managed_devices {
    device_id        = "device2-1234-1234-1234-123456789abc"
    script_policy_id = "network-fix-script-guid"
  }

  managed_devices {
    device_id        = "device3-1234-1234-1234-123456789abc"
    script_policy_id = "printer-repair-script-guid"
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 4: Urgent remediation from variable
variable "urgent_remediation" {
  description = "Devices requiring urgent remediation"
  type = map(object({
    device_id        = string
    script_policy_id = string
  }))
  default = {
    "critical1" = {
      device_id        = "aaaa1111-1111-1111-1111-111111111111"
      script_policy_id = "emergency-fix-script-guid"
    }
    "critical2" = {
      device_id        = "bbbb2222-2222-2222-2222-222222222222"
      script_policy_id = "emergency-fix-script-guid"
    }
  }
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "urgent_fix" {
  dynamic "managed_devices" {
    for_each = var.urgent_remediation
    content {
      device_id        = managed_devices.value.device_id
      script_policy_id = managed_devices.value.script_policy_id
    }
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 5: Post-incident remediation
locals {
  incident_devices = [
    {
      device_id        = "incident1-1111-1111-1111-111111111111"
      script_policy_id = "security-hardening-script-guid"
    },
    {
      device_id        = "incident2-2222-2222-2222-222222222222"
      script_policy_id = "security-hardening-script-guid"
    }
  ]
}

action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "incident_remediation" {
  dynamic "managed_devices" {
    for_each = local.incident_devices
    content {
      device_id        = managed_devices.value.device_id
      script_policy_id = managed_devices.value.script_policy_id
    }
  }

  timeouts = {
    invoke = "15m"
  }
}

# Example 6: Co-managed device remediation
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "comanaged" {
  comanaged_devices {
    device_id        = "comanaged-1234-1234-1234-123456789abc"
    script_policy_id = "sccm-integration-fix-script-guid"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Example 7: Troubleshooting specific issue
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "troubleshoot_vpn" {
  managed_devices {
    device_id        = "vpn-issue-device-1234-123456789abc"
    script_policy_id = "vpn-troubleshoot-script-guid"
  }

  timeouts = {
    invoke = "10m"
  }
}

# Example 8: Testing new remediation script
action "microsoft365_graph_beta_device_management_managed_device_initiate_on_demand_proactive_remediation" "test_new_script" {
  managed_devices {
    device_id        = "test-device-1234-1234-1234-123456789abc"
    script_policy_id = "new-remediation-test-script-guid"
  }

  timeouts = {
    invoke = "5m"
  }
}

# Output examples
output "remediation_summary" {
  value = {
    managed_count   = length([for d in action.multiple_same_script.managed_devices : d])
    comanaged_count = length([for d in action.comanaged.comanaged_devices : d])
  }
  description = "Count of devices with remediation initiated"
}

