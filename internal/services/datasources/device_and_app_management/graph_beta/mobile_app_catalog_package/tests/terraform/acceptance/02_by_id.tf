data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "by_id" {
  filter_type  = "id"
  filter_value = "3a6307ef-6991-faf1-01e1-35e1557287aa"
}

output "by_id_package_found" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0
  description = "Whether the package was found"
}

output "by_id_package_id" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].id : null
  description = "Package ID"
}

output "by_id_display_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].display_name : null
  description = "Display name"
}

output "by_id_publisher" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].publisher : null
  description = "Publisher"
}

output "by_id_description" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].description : null
  description = "Package description"
}

output "by_id_file_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].file_name : null
  description = "Installer file name"
}

output "by_id_size" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].size : null
  description = "Package size in bytes"
}

output "by_id_display_version" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].display_version : null
  description = "Display version"
}

output "by_id_install_command" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].install_command_line : null
  description = "Install command line"
}

output "by_id_uninstall_command" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].uninstall_command_line : null
  description = "Uninstall command line"
}

output "by_id_setup_file_path" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].setup_file_path : null
  description = "Setup file path"
}

output "by_id_allowed_architectures" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].allowed_architectures : null
  description = "Allowed architectures"
}

output "by_id_minimum_windows_release" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].minimum_supported_windows_release : null
  description = "Minimum supported Windows release"
}

output "by_id_allow_available_uninstall" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].allow_available_uninstall : null
  description = "Allow available uninstall"
}

output "by_id_catalog_package_id" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].mobile_app_catalog_package_id : null
  description = "Mobile app catalog package ID"
}

output "by_id_rules" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? [
    for rule in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].rules : {
      odata_type              = rule.odata_type
      rule_type               = rule.rule_type
      path                    = rule.path
      file_or_folder_name     = rule.file_or_folder_name
      key_path                = rule.key_path
      value_name              = rule.value_name
      operation_type          = rule.operation_type
      operator                = rule.operator
      comparison_value        = rule.comparison_value
      check_32bit_on_64system = rule.check_32bit_on_64system
    }
  ] : []
  description = "Detection and requirement rules"
}

output "by_id_install_experience" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? {
    run_as_account          = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].install_experience.run_as_account
    max_run_time_in_minutes = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].install_experience.max_run_time_in_minutes
    device_restart_behavior = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].install_experience.device_restart_behavior
  } : null
  description = "Install experience settings"
}

output "by_id_return_codes" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? [
    for rc in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].return_codes : {
      return_code = rc.return_code
      type        = rc.type
    }
  ] : []
  description = "Return codes"
}

output "by_id_msi_information" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? (
    data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].msi_information != null ? {
      product_code    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].msi_information.product_code
      product_version = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].msi_information.product_version
      upgrade_code    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].msi_information.upgrade_code
      requires_reboot = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].msi_information.requires_reboot
      package_type    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].msi_information.package_type
      product_name    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].msi_information.product_name
      publisher       = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].msi_information.publisher
    } : null
  ) : null
  description = "MSI package information"
}

output "by_id_privacy_url" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].privacy_information_url : null
  description = "Privacy information URL"
}

output "by_id_information_url" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].information_url : null
  description = "Information URL"
}

output "by_id_developer" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_id.items[0].developer : null
  description = "Developer"
}
