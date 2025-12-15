data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "startswith(publisherDisplayName, 'Microsoft')"
  odata_top    = 10
}

output "odata_packages_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items)
  description = "Number of packages returned by OData filter"
}

output "odata_first_display_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].display_name : null
  description = "Display name of first package"
}

output "odata_first_publisher" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].publisher : null
  description = "Publisher (should start with 'Microsoft')"
}

output "odata_first_version" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].display_version : null
  description = "Display version"
}

output "odata_first_file_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].file_name : null
  description = "File name"
}

output "odata_first_size" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].size : null
  description = "Package size in bytes"
}

output "odata_first_catalog_id" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].mobile_app_catalog_package_id : null
  description = "Catalog package ID"
}

output "odata_first_install_command" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].install_command_line : null
  description = "Install command"
}

output "odata_first_uninstall_command" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].uninstall_command_line : null
  description = "Uninstall command"
}

output "odata_first_architectures" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].allowed_architectures : null
  description = "Allowed architectures"
}

output "odata_first_min_windows_release" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].minimum_supported_windows_release : null
  description = "Minimum supported Windows release"
}

output "odata_first_rules_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].rules) : 0
  description = "Number of detection/requirement rules"
}

output "odata_first_return_codes_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].return_codes) : 0
  description = "Number of return codes"
}

output "odata_first_install_experience" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items) > 0 ? {
    run_as_account            = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].install_experience.run_as_account
    max_run_time_in_minutes   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].install_experience.max_run_time_in_minutes
    device_restart_behavior   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items[0].install_experience.device_restart_behavior
  } : null
  description = "Install experience settings"
}

output "odata_packages_summary" {
  value = [
    for pkg in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items : {
      display_name                  = pkg.display_name
      publisher                     = pkg.publisher
      display_version               = pkg.display_version
      file_name                     = pkg.file_name
      size                          = pkg.size
      mobile_app_catalog_package_id = pkg.mobile_app_catalog_package_id
      has_rules                     = length(pkg.rules) > 0
      has_return_codes              = length(pkg.return_codes) > 0
      has_msi_info                  = pkg.msi_information != null
    }
  ]
  description = "Summary of all packages returned by OData filter"
}
