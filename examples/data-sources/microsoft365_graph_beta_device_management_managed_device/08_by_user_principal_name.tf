# Example: Get managed devices by user principal name (email)
data "microsoft365_graph_beta_device_management_managed_device" "by_user" {
  user_principal_name = "user@contoso.com"
}

# Output: Devices assigned to the specified user
output "user_devices" {
  value = [
    for device in data.microsoft365_graph_beta_device_management_managed_device.by_user.items :
    {
      id                = device.id
      device_name       = device.device_name
      user_display_name = device.user_display_name
      user_email        = device.user_principal_name
      operating_system  = device.operating_system
      compliance_state  = device.compliance_state
      last_sync         = device.last_sync_date_time
    }
  ]
  description = "List of devices assigned to the specified user"
}
