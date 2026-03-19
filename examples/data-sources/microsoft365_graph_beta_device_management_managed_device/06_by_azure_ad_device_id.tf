# Example: Get managed device by Azure AD Device ID
data "microsoft365_graph_beta_device_management_managed_device" "by_aad_id" {
  azure_ad_device_id = "aaaaaaaa-0000-0000-0000-000000000001"
}

# Output: Device information by Azure AD Device ID
output "device_by_aad_id" {
  value = length(data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items) > 0 ? {
    intune_id          = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].id
    azure_ad_device_id = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].azure_ad_device_id
    device_name        = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].device_name
    registration_state = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].device_registration_state
    user               = data.microsoft365_graph_beta_device_management_managed_device.by_aad_id.items[0].user_principal_name
  } : null
  description = "Device information using Azure AD Device ID lookup"
}
