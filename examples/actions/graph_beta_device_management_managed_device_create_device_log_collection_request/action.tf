# Example 1: Create device log collection request for a single device - Minimal
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "create_single" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      }
    ]
  }
}

# Example 2: Create log collection requests for multiple devices
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "create_multiple" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      },
      {
        device_id     = "87654321-4321-4321-4321-ba9876543210"
        template_type = "custom"
      }
    ]

    timeouts = {
      invoke = "10m"
    }
  }
}

# Example 3: Create with validation - Maximal
action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "create_maximal" {
  config {
    managed_devices = [
      {
        device_id     = "12345678-1234-1234-1234-123456789abc"
        template_type = "predefined"
      }
    ]

    comanaged_devices = [
      {
        device_id     = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
        template_type = "predefined"
      }
    ]

    ignore_partial_failures = true
    validate_device_exists  = true

    timeouts = {
      invoke = "5m"
    }
  }
}

# Example 4: Create log collection requests for non-compliant devices
data "microsoft365_graph_beta_device_management_managed_device" "noncompliant_devices" {
  filter_type  = "odata"
  odata_filter = "complianceState eq 'noncompliant'"
}

action "microsoft365_graph_beta_device_management_managed_device_create_device_log_collection_request" "create_noncompliant" {
  config {
    managed_devices = [
      for device in data.microsoft365_graph_beta_device_management_managed_device.noncompliant_devices.items : {
        device_id     = device.id
        template_type = "predefined"
      }
    ]

    validate_device_exists = true

    timeouts = {
      invoke = "20m"
    }
  }
}
