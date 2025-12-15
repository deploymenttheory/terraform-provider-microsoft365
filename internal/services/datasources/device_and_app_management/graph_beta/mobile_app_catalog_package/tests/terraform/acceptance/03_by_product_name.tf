data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "by_product_name" {
  filter_type  = "product_name"
  filter_value = "7-Zip"
}

output "product_name_packages_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items)
  description = "Number of packages matching product name"
}

output "product_name_first_display_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].display_name : null
  description = "Display name of first matching package"
}

output "product_name_first_publisher" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].publisher : null
  description = "Publisher of first matching package"
}

output "product_name_first_version" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].display_version : null
  description = "Version of first matching package"
}

output "product_name_first_file_name" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].file_name : null
  description = "File name of first matching package"
}

output "product_name_first_size" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].size : null
  description = "Size of first matching package"
}

output "product_name_first_catalog_id" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].mobile_app_catalog_package_id : null
  description = "Catalog package ID of first matching package"
}

output "product_name_first_install_command" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].install_command_line : null
  description = "Install command of first matching package"
}

output "product_name_first_uninstall_command" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].uninstall_command_line : null
  description = "Uninstall command of first matching package"
}

output "product_name_first_rules_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].rules) : 0
  description = "Number of rules in first matching package"
}

output "product_name_first_return_codes_count" {
  value       = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items) > 0 ? length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items[0].return_codes) : 0
  description = "Number of return codes in first matching package"
}

output "product_name_all_packages" {
  value = [
    for pkg in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items : {
      id                            = pkg.id
      display_name                  = pkg.display_name
      publisher                     = pkg.publisher
      display_version               = pkg.display_version
      file_name                     = pkg.file_name
      size                          = pkg.size
      mobile_app_catalog_package_id = pkg.mobile_app_catalog_package_id
      allowed_architectures         = pkg.allowed_architectures
      rules_count                   = length(pkg.rules)
      return_codes_count            = length(pkg.return_codes)
    }
  ]
  description = "All matching packages with key details"
}
