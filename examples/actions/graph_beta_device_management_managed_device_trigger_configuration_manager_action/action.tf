# Example 1: Trigger Configuration Manager action on a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "trigger_single" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
        action    = "refreshMachinePolicy"
      }
    ]
  }
}

# Example 2: Trigger multiple Configuration Manager actions
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "trigger_multiple" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
        action    = "refreshMachinePolicy"
      },
      {
        device_id = "87654321-4321-4321-4321-ba9876543210"
        action    = "refreshUserPolicy"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Trigger with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "trigger_maximal" {
  config {
    managed_devices = [
      {
        device_id = "12345678-1234-1234-1234-123456789abc"
        action    = "appEvaluation"
      }
    ]

    comanaged_devices = [
      {
        device_id = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        action    = "refreshMachinePolicy"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Trigger policy refresh on all co-managed devices
data "microsoft365_graph_beta_device_management_managed_device" "comanaged_devices" {
  filter_type  = "odata"
  odata_filter = "managementAgent eq 'configurationManagerClientMdm'"
}

action "microsoft365_graph_beta_device_management_managed_device_trigger_configuration_manager_action" "refresh_all_comanaged" {
  config {
    comanaged_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.comanaged_devices.items : {
        device_id = device.id
        action    = "refreshMachinePolicy"
      }
    ]

    ignore_partial_failures = true

    timeouts = {
      invoke = "20m"
    }
  }
}
