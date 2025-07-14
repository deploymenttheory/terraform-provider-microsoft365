# Example: Windows Autopilot Device Identity with user assignment
# Note: Hardware hash must be collected before using this resource
# See: https://learn.microsoft.com/en-us/autopilot/add-devices for collection methods

resource "microsoft365_graph_beta_device_management_windows_autopilot_device_identity" "example_with_user" {
  # Required field - must be unique per device
  serial_number = "EXAMPLESERIAL123"

  # Optional fields for device identification and organization
  group_tag                 = "Finance-Dept" # Used for targeting specific deployment profiles
  purchase_order_identifier = "PO-12345"     # For tracking procurement information
  product_key               = "XXXXX-XXXXX-XXXXX-XXXXX-XXXXX"
  display_name              = "Finance-Laptop-01"
  manufacturer              = "Example Manufacturer"
  model                     = "Example Model"

  # User assignment configuration - enables personalized setup
  # Warning: Ensure the UPN exists or device may become inaccessible
  user_assignment {
    user_principal_name = "finance-user@example.com"
    # addressable_user_name is computed and returned by the API
  }

  timeouts {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
}

# Example 2: Windows Autopilot Device Identity without user assignment
resource "microsoft365_graph_beta_device_management_windows_autopilot_device_identity" "example_without_user" {
  serial_number             = "EXAMPLESERIAL456"
  group_tag                 = "Example-Group-2"
  purchase_order_identifier = "PO-67890"
  product_key               = "YYYYY-YYYYY-YYYYY-YYYYY-YYYYY"
  display_name              = "Example Device without User"
  manufacturer              = "Example Manufacturer"
  model                     = "Example Model"

  # No user_assignment block means no user will be assigned to this device

  timeouts = {
    create = "30m"
    read   = "10m"
    update = "30m"
    delete = "10m"
  }
} 