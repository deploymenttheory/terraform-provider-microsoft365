data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "all" {
  filter_type = "all"
}

output "all_packages_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items)
  description = "Total number of packages retrieved"
}

output "first_package_id" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].id : null
  description = "ID of the first package"
}

output "first_package_display_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].display_name : null
  description = "Display name of the first package"
}

output "first_package_publisher" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].publisher : null
  description = "Publisher of the first package"
}

output "first_package_file_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].file_name : null
  description = "File name of the first package"
}

output "first_package_size" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].size : null
  description = "Size in bytes of the first package"
}

output "first_package_display_version" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].display_version : null
  description = "Display version of the first package"
}

output "first_package_catalog_package_id" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].mobile_app_catalog_package_id : null
  description = "Mobile app catalog package ID"
}

output "first_package_install_command" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].install_command_line : null
  description = "Install command line"
}

output "first_package_uninstall_command" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].uninstall_command_line : null
  description = "Uninstall command line"
}

output "first_package_allowed_architectures" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].allowed_architectures : null
  description = "Allowed architectures"
}

output "first_package_rules_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].rules) : 0
  description = "Number of detection/requirement rules"
}

output "first_package_return_codes_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].return_codes) : 0
  description = "Number of return codes"
}

output "first_package_install_experience" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? {
    run_as_account          = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].install_experience.run_as_account
    max_run_time_in_minutes = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].install_experience.max_run_time_in_minutes
    device_restart_behavior = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].install_experience.device_restart_behavior
  } : null
  description = "Install experience settings"
}

output "first_package_msi_information" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items) > 0 ? (
    data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].msi_information != null ? {
      product_code    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].msi_information.product_code
      product_version = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].msi_information.product_version
      package_type    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].msi_information.package_type
      requires_reboot = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all.items[0].msi_information.requires_reboot
    } : null
  ) : null
  description = "MSI information (if applicable)"
}
