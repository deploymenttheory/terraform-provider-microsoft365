# Example 1: Get all mobile app catalog packages
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "all_packages" {
  filter_type = "all"
  timeouts = {
    read = "30s"
  }
}

# Example output for all packages
output "all_catalog_package" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all_packages.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all_packages.items : []
}

# More focused output showing just key information
output "all_packages_summary" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all_packages.items != null ? [
    for package in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all_packages.items : {
      id                        = package.id
      product_id                = package.product_id
      product_display_name      = package.product_display_name
      publisher_display_name    = package.publisher_display_name
      version_display_name      = package.version_display_name
      branch_display_name       = package.branch_display_name
      applicable_architectures  = package.applicable_architectures
      package_auto_update_capable = package.package_auto_update_capable
    }
  ] : []
}

# Example 2: Get a specific package by product ID
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "by_product_id" {
  filter_type  = "id"
  filter_value = "3a6307ef-6991-faf1-01e1-35e1557287aa" # Replace with actual product ID

  timeouts = {
    read = "30s"
  }
}

# Output for by_product_id
output "package_by_id" {
  value = try(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_id.items[0], null)
}

# Example 3: Get packages by product name (partial match)
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "by_product_name" {
  filter_type  = "product_name"
  filter_value = "7-Zip" # This will find all packages with "7-Zip" in the product name

  timeouts = {
    read = "30s"
  }
}

# Output for by_product_name
output "packages_by_product_name" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items : []
}

# Example 4: Get packages by publisher name (partial match)
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "by_publisher_name" {
  filter_type  = "publisher_name"
  filter_value = "Microsoft" # This will find all packages with "Microsoft" in the publisher name

  timeouts = {
    read = "30s"
  }
}

# Output for by_publisher_name
output "packages_by_publisher_name" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_publisher_name.items : []
}

# Example 5: Get packages using OData filter
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_filter" {
  filter_type    = "odata"
  odata_filter   = "productDisplayName eq '7-Zip'"
  odata_count    = true
  odata_orderby  = "productDisplayName"
  odata_top      = 10

  timeouts = {
    read = "30s"
  }
}

# Output for OData filter
output "packages_odata_filter" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items : []
}

# Example 6: Advanced OData query with multiple parameters
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_advanced" {
  filter_type    = "odata"
  odata_filter   = "contains(productDisplayName, 'Microsoft')"
  odata_top      = 5
  odata_skip     = 0
  odata_select   = "id,productId,productDisplayName,publisherDisplayName,versionDisplayName"
  odata_orderby  = "productDisplayName"
  odata_count    = true
  odata_search   = "\"productDisplayName:Microsoft\""

  timeouts = {
    read = "30s"
  }
}

# Output for advanced OData query
output "packages_odata_advanced" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced.items : []
}

# Example 7: Search-only OData query
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_search_only" {
  filter_type  = "odata"
  odata_search = "\"productDisplayName:Microsoft\""

  timeouts = {
    read = "30s"
  }
}

# Output for search-only query
output "packages_search_only" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_search_only.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_search_only.items : []
}
