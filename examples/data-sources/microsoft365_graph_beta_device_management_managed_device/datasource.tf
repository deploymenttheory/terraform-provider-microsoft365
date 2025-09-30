data "microsoft365_graph_beta_device_management_managed_device" "all_devices" {
  filter_type = "all"
}

output "first_device" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0]
}

output "first_device_id" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].id
}

output "first_device_name" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].device_name
}

output "first_device_operating_system" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].operating_system
}

output "first_device_user_principal_name" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].user_principal_name
}

output "first_device_serial_number" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].serial_number
}

output "first_device_model" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].model
}

output "first_device_manufacturer" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].manufacturer
}

output "first_device_compliance_state" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].compliance_state
}

output "first_device_enrolled_date_time" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].enrolled_date_time
}

output "first_device_last_sync_date_time" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].last_sync_date_time
}

output "first_device_owner_type" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].owner_type
}

output "first_device_management_state" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].management_state
}

output "first_device_jail_broken" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].jail_broken
}

output "first_device_os_version" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].os_version
}

output "first_device_aad_registered" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].aad_registered
}

output "first_device_device_type" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].device_type
}

output "first_device_email_address" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].email_address
}

output "first_device_is_supervised" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].is_supervised
}

output "first_device_is_encrypted" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].is_encrypted
}

output "first_device_user_id" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].user_id
}

output "first_device_device_registration_state" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].device_registration_state
}

output "first_device_device_category_display_name" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].device_category_display_name
}

output "first_device_azure_ad_device_id" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].azure_ad_device_id
}

output "first_device_azure_active_directory_device_id" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].azure_active_directory_device_id
}

output "first_device_managed_device_owner_type" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].managed_device_owner_type
}

output "first_device_management_agent" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].management_agent
}

output "first_device_eas_activated" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].eas_activated
}

output "first_device_eas_device_id" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].eas_device_id
}

output "first_device_eas_activation_date_time" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].eas_activation_date_time
}

output "first_device_lost_mode_state" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].lost_mode_state
}

output "first_device_activation_lock_bypass_code" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].activation_lock_bypass_code
}

output "first_device_exchange_last_successful_sync_date_time" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].exchange_last_successful_sync_date_time
}

output "first_device_exchange_access_state" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].exchange_access_state
}

output "first_device_exchange_access_state_reason" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].exchange_access_state_reason
}

output "first_device_remote_assistance_session_url" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].remote_assistance_session_url
}

output "first_device_remote_assistance_session_error_details" {
  value = data.microsoft365_graph_beta_device_management_managed_device.all_devices.items[0].remote_assistance_session_error_details
} 