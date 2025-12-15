# Mobile App Catalog Package Data Source Examples
#
# This data source retrieves complete win32CatalogApp details from Microsoft Intune,
# including installation commands, detection rules, return codes, and MSI information.
#
# The data source performs two API calls:
# 1. Searches for mobile app catalog packages based on your filter
# 2. Converts each package to a complete win32CatalogApp with all deployment details
#
# IMPORTANT: Microsoft Graph API Limitations for mobileAppCatalogPackages endpoint:
# ✅ Supported OData features:
#    - $filter with startswith() function (e.g., startswith(publisherDisplayName, 'value'))
#    - $top for limiting results
# ❌ Not supported/problematic OData features:
#    - $skip - causes 500 errors and timeouts
#    - $select - causes 500 errors and timeouts
#    - $orderby - returns no results when combined with $filter
#    - $count - returns no results when combined with $filter
#    - $search - not reliably supported
#    - eq operator in filters - not reliable, use startswith() instead
#
# For best results, use the simple filter types (all, id, product_name, publisher_name)
# or OData with only $filter (using startswith()) and $top parameters.

# ============================================================================
# Example 1: Get all mobile app catalog packages with complete details
# ============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "all_packages" {
  filter_type = "all"
  timeouts = {
    read = "5m" # Increased timeout due to conversion API calls
  }
}

# Output showing complete app details including install commands and detection rules
output "all_packages_detailed" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all_packages.items != null ? [
    for pkg in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all_packages.items : {
      # Basic app information
      id                            = pkg.id
      display_name                  = pkg.display_name
      publisher                     = pkg.publisher
      description                   = pkg.description
      display_version               = pkg.display_version
      
      # Win32 specific deployment details
      file_name                     = pkg.file_name
      size                          = pkg.size
      install_command_line          = pkg.install_command_line
      uninstall_command_line        = pkg.uninstall_command_line
      setup_file_path               = pkg.setup_file_path
      allowed_architectures         = pkg.allowed_architectures
      minimum_supported_windows_release = pkg.minimum_supported_windows_release
      
      # Catalog reference
      mobile_app_catalog_package_id = pkg.mobile_app_catalog_package_id
      
      # Detection rules
      rules_count                   = length(pkg.rules)
      has_detection_rules           = length(pkg.rules) > 0
      
      # Installation settings
      run_as_account                = pkg.install_experience != null ? pkg.install_experience.run_as_account : null
      max_run_time_in_minutes       = pkg.install_experience != null ? pkg.install_experience.max_run_time_in_minutes : null
      
      # Return codes
      return_codes_count            = length(pkg.return_codes)
      
      # MSI information (if applicable)
      has_msi_info                  = pkg.msi_information != null
      msi_product_code              = pkg.msi_information != null ? pkg.msi_information.product_code : null
    }
  ] : []
  description = "All mobile app catalog packages with deployment details"
}

# ============================================================================
# Example 2: Get a specific package by product ID (7-Zip example)
# ============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "seven_zip" {
  filter_type  = "id"
  filter_value = "3a6307ef-6991-faf1-01e1-35e1557287aa" # 7-Zip product ID
  
  timeouts = {
    read = "2m"
  }
}

# Output showing complete 7-Zip deployment configuration
output "seven_zip_full_details" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items) > 0 ? {
    # Application metadata
    id                      = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].id
    display_name            = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].display_name
    description             = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].description
    publisher               = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].publisher
    developer               = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].developer
    display_version         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].display_version
    
    # URLs
    privacy_url             = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].privacy_information_url
    information_url         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].information_url
    
    # Installation details
    file_name               = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].file_name
    size_bytes              = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].size
    install_command         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].install_command_line
    uninstall_command       = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].uninstall_command_line
    setup_file_path         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].setup_file_path
    
    # Requirements
    allowed_architectures   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].allowed_architectures
    minimum_windows_release = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].minimum_supported_windows_release
    allow_uninstall         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].allow_available_uninstall
    
    # Detection rules
    detection_rules = [
      for rule in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].rules : {
        type               = rule.odata_type
        rule_type          = rule.rule_type
        path               = rule.path
        file_or_folder     = rule.file_or_folder_name
        key_path           = rule.key_path
        value_name         = rule.value_name
        operation_type     = rule.operation_type
        operator           = rule.operator
        comparison_value   = rule.comparison_value
        check_32bit_on_64  = rule.check_32bit_on_64system
      }
    ]
    
    # Install experience
    install_experience = {
      run_as_account          = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].install_experience.run_as_account
      max_run_time_minutes    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].install_experience.max_run_time_in_minutes
      device_restart_behavior = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].install_experience.device_restart_behavior
    }
    
    # Return codes
    return_codes = [
      for rc in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].return_codes : {
        code = rc.return_code
        type = rc.type
      }
    ]
    
    # MSI information (if applicable)
    msi_info = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].msi_information != null ? {
      product_code    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].msi_information.product_code
      product_version = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].msi_information.product_version
      upgrade_code    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].msi_information.upgrade_code
      requires_reboot = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].msi_information.requires_reboot
      package_type    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].msi_information.package_type
      product_name    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].msi_information.product_name
    } : null
    
    # Catalog reference
    catalog_package_id = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].mobile_app_catalog_package_id
  } : null
  description = "Complete 7-Zip deployment configuration"
}

# ============================================================================
# Example 3: Get packages by product name (partial match)
# ============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "by_product_name" {
  filter_type  = "product_name"
  filter_value = "Docker" # This will find all packages with "Docker" in the name
  
  timeouts = {
    read = "3m"
  }
}

# Output showing key deployment details for matched packages
output "docker_packages" {
  value = [
    for pkg in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.by_product_name.items : {
      display_name          = pkg.display_name
      display_version       = pkg.display_version
      publisher             = pkg.publisher
      file_name             = pkg.file_name
      size_mb               = floor(pkg.size / 1024 / 1024)
      install_command       = pkg.install_command_line
      architectures         = pkg.allowed_architectures
      detection_rules_count = length(pkg.rules)
      catalog_package_id    = pkg.mobile_app_catalog_package_id
    }
  ]
  description = "Docker packages with deployment details"
}

# ============================================================================
# Example 4: Get packages by publisher name (partial match)
# ============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "microsoft_apps" {
  filter_type  = "publisher_name"
  filter_value = "Microsoft" # Finds all Microsoft published apps
  
  timeouts = {
    read = "5m"
  }
}

# Output summarizing Microsoft applications
output "microsoft_apps_summary" {
  value = [
    for pkg in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.microsoft_apps.items : {
      name              = pkg.display_name
      version           = pkg.display_version
      file_name         = pkg.file_name
      size_mb           = floor(pkg.size / 1024 / 1024)
      min_windows       = pkg.minimum_supported_windows_release
      has_msi_info      = pkg.msi_information != null
      detection_rules   = length(pkg.rules)
      return_codes      = length(pkg.return_codes)
    }
  ]
  description = "Summary of Microsoft published applications"
}

# ============================================================================
# Example 5: Using OData filter with startswith()
# ============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_google" {
  filter_type  = "odata"
  odata_filter = "startswith(publisherDisplayName, 'Google')"
  odata_top    = 5
  
  timeouts = {
    read = "3m"
  }
}

output "google_apps" {
  value = [
    for pkg in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_google.items : {
      display_name      = pkg.display_name
      version           = pkg.display_version
      install_command   = pkg.install_command_line
      uninstall_command = pkg.uninstall_command_line
      file_name         = pkg.file_name
      architectures     = pkg.allowed_architectures
    }
  ]
  description = "Google applications with install commands"
}

# ============================================================================
# Example 6: Practical use case - Extract detection rules for a specific app
# ============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "app_for_rules" {
  filter_type  = "product_name"
  filter_value = "7-Zip"
  
  timeouts = {
    read = "2m"
  }
}

# Output only the detection rules in a readable format
output "detection_rules_example" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.app_for_rules.items) > 0 ? {
    app_name = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.app_for_rules.items[0].display_name
    rules = [
      for rule in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.app_for_rules.items[0].rules : {
        type = rule.odata_type == "#microsoft.graph.win32LobAppFileSystemRule" ? "File System" : "Registry"
        location = rule.odata_type == "#microsoft.graph.win32LobAppFileSystemRule" ? (
          "${rule.path}\\${rule.file_or_folder_name}"
        ) : (
          "${rule.key_path} [${rule.value_name}]"
        )
        operation    = rule.operation_type
        operator     = rule.operator
        expected_value = rule.comparison_value
      }
    ]
  } : null
  description = "Detection rules in readable format"
}

# ============================================================================
# Example 7: Extract installation requirements
# ============================================================================
output "installation_requirements" {
  value = length(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items) > 0 ? {
    app_name                = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].display_name
    file_name               = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].file_name
    size_mb                 = floor(data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].size / 1024 / 1024)
    allowed_architectures   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].allowed_architectures
    minimum_windows_release = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].minimum_supported_windows_release
    minimum_disk_space_mb   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].minimum_free_disk_space_in_mb
    minimum_memory_mb       = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].minimum_memory_in_mb
    run_as                  = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].install_experience.run_as_account
    max_runtime_minutes     = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].install_experience.max_run_time_in_minutes
  } : null
  description = "Application installation requirements"
}

# ============================================================================
# Example 8: Compare multiple package versions
# ============================================================================
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "all_versions" {
  filter_type = "all"
  
  timeouts = {
    read = "5m"
  }
}

# Group packages by product name to show version comparisons
output "packages_by_product" {
  value = {
    for pkg in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.all_versions.items :
    pkg.display_name => {
      version           = pkg.display_version
      publisher         = pkg.publisher
      size_mb           = floor(pkg.size / 1024 / 1024)
      catalog_id        = pkg.mobile_app_catalog_package_id
      file_name         = pkg.file_name
    }...
  }
  description = "All packages grouped by product name"
}
