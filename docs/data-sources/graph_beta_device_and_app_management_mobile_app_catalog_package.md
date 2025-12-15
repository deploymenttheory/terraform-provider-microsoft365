---
page_title: "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package Data Source - terraform-provider-microsoft365"
subcategory: "Device and App Management"

description: |-
  Retrieves mobile app catalog packages from Microsoft Intune using the /deviceAppManagement/MobileAppCatalogPackage endpoint. This data source enables querying mobile app catalog packages with advanced filtering capabilities including OData queries for filtering by product name, publisher, and other properties.
---

# microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package (Data Source)

Retrieves mobile app catalog packages from Microsoft Intune using the `/deviceAppManagement/MobileAppCatalogPackage` endpoint. This data source enables querying mobile app catalog packages with advanced filtering capabilities including OData queries for filtering by product name, publisher, and other properties.

## Microsoft Documentation

- [mobileAppCatalogPackage resource type](https://learn.microsoft.com/en-us/intune/intune-service/apps/apps-enterprise-app-management)

## API Permissions

The following API permissions are required in order to use this data source.

### Microsoft Graph

- **Application**: `DeviceManagementApps.Read.All`, `DeviceManagementApps.ReadWrite.All`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.32.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
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
      id              = pkg.id
      display_name    = pkg.display_name
      publisher       = pkg.publisher
      description     = pkg.description
      display_version = pkg.display_version

      # Win32 specific deployment details
      file_name                         = pkg.file_name
      size                              = pkg.size
      install_command_line              = pkg.install_command_line
      uninstall_command_line            = pkg.uninstall_command_line
      setup_file_path                   = pkg.setup_file_path
      allowed_architectures             = pkg.allowed_architectures
      minimum_supported_windows_release = pkg.minimum_supported_windows_release

      # Catalog reference
      mobile_app_catalog_package_id = pkg.mobile_app_catalog_package_id

      # Detection rules
      rules_count         = length(pkg.rules)
      has_detection_rules = length(pkg.rules) > 0

      # Installation settings
      run_as_account          = pkg.install_experience != null ? pkg.install_experience.run_as_account : null
      max_run_time_in_minutes = pkg.install_experience != null ? pkg.install_experience.max_run_time_in_minutes : null

      # Return codes
      return_codes_count = length(pkg.return_codes)

      # MSI information (if applicable)
      has_msi_info     = pkg.msi_information != null
      msi_product_code = pkg.msi_information != null ? pkg.msi_information.product_code : null
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
    id              = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].id
    display_name    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].display_name
    description     = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].description
    publisher       = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].publisher
    developer       = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].developer
    display_version = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].display_version

    # URLs
    privacy_url     = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].privacy_information_url
    information_url = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].information_url

    # Installation details
    file_name         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].file_name
    size_bytes        = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].size
    install_command   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].install_command_line
    uninstall_command = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].uninstall_command_line
    setup_file_path   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].setup_file_path

    # Requirements
    allowed_architectures   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].allowed_architectures
    minimum_windows_release = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].minimum_supported_windows_release
    allow_uninstall         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].allow_available_uninstall

    # Detection rules
    detection_rules = [
      for rule in data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.seven_zip.items[0].rules : {
        type              = rule.odata_type
        rule_type         = rule.rule_type
        path              = rule.path
        file_or_folder    = rule.file_or_folder_name
        key_path          = rule.key_path
        value_name        = rule.value_name
        operation_type    = rule.operation_type
        operator          = rule.operator
        comparison_value  = rule.comparison_value
        check_32bit_on_64 = rule.check_32bit_on_64system
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
      name            = pkg.display_name
      version         = pkg.display_version
      file_name       = pkg.file_name
      size_mb         = floor(pkg.size / 1024 / 1024)
      min_windows     = pkg.minimum_supported_windows_release
      has_msi_info    = pkg.msi_information != null
      detection_rules = length(pkg.rules)
      return_codes    = length(pkg.return_codes)
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
        operation      = rule.operation_type
        operator       = rule.operator
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
      version    = pkg.display_version
      publisher  = pkg.publisher
      size_mb    = floor(pkg.size / 1024 / 1024)
      catalog_id = pkg.mobile_app_catalog_package_id
      file_name  = pkg.file_name
    }...
  }
  description = "All packages grouped by product name"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `product_name`, `publisher_name`, `odata`.

### Optional

- `filter_value` (String) Value to filter by. Not required when filter_type is 'all' or 'odata'.
- `odata_count` (Boolean) OData $count parameter to include count of total results. Only used when filter_type is 'odata'.
- `odata_expand` (String) OData $expand parameter to include related entities. Only used when filter_type is 'odata'.
- `odata_filter` (String) OData $filter parameter for filtering results. Only used when filter_type is 'odata'. Example: productDisplayName eq 'Microsoft Office'.
- `odata_orderby` (String) OData $orderby parameter to sort results. Only used when filter_type is 'odata'. Example: productDisplayName.
- `odata_search` (String) OData $search parameter for full-text search. Only used when filter_type is 'odata'.
- `odata_select` (String) OData $select parameter to specify which fields to include. Only used when filter_type is 'odata'.
- `odata_skip` (Number) OData $skip parameter for pagination. Only used when filter_type is 'odata'.
- `odata_top` (Number) OData $top parameter to limit the number of results. Only used when filter_type is 'odata'.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `items` (Attributes List) The list of win32 catalog applications with full details that match the filter criteria. (see [below for nested schema](#nestedatt--items))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--items"></a>
### Nested Schema for `items`

Read-Only:

- `allow_available_uninstall` (Boolean) Indicates whether the app can be uninstalled from the available context.
- `allowed_architectures` (String) The allowed target architectures for this app.
- `applicable_architectures` (String) The Windows architecture(s) on which this app can run.
- `committed_content_version` (String) The internal committed content version.
- `created_date_time` (String) The date and time the app was created.
- `dependent_app_count` (Number) The total number of dependencies the child app has.
- `description` (String) The description of the application.
- `developer` (String) The developer of the app.
- `display_name` (String) The display name of the application.
- `display_version` (String) The display version of the application.
- `file_name` (String) The name of the main installation file.
- `id` (String) The unique identifier for the application.
- `information_url` (String) The more information URL.
- `install_command_line` (String) The command line to install this app.
- `install_experience` (Attributes) The install experience for this app. (see [below for nested schema](#nestedatt--items--install_experience))
- `is_assigned` (Boolean) Indicates whether the app is assigned to at least one group.
- `is_featured` (Boolean) Indicates whether the app is marked as featured by the admin.
- `last_modified_date_time` (String) The date and time the app was last modified.
- `minimum_cpu_speed_in_mhz` (Number) The minimum CPU speed required to install this app.
- `minimum_free_disk_space_in_mb` (Number) The minimum free disk space required to install this app.
- `minimum_memory_in_mb` (Number) The minimum memory required to install this app.
- `minimum_number_of_processors` (Number) The minimum number of processors required to install this app.
- `minimum_supported_windows_release` (String) The minimum supported Windows release version.
- `mobile_app_catalog_package_id` (String) The mobile app catalog package ID.
- `msi_information` (Attributes) The MSI information for MSI-based apps. (see [below for nested schema](#nestedatt--items--msi_information))
- `notes` (String) Notes for the app.
- `owner` (String) The owner of the app.
- `privacy_information_url` (String) The privacy statement URL.
- `publisher` (String) The publisher of the application.
- `publishing_state` (String) The publishing state for the app.
- `return_codes` (Attributes List) The return codes for post installation behavior. (see [below for nested schema](#nestedatt--items--return_codes))
- `role_scope_tag_ids` (List of String) List of scope tag IDs for this app.
- `rules` (Attributes List) The detection and requirement rules for this app. (see [below for nested schema](#nestedatt--items--rules))
- `setup_file_path` (String) The relative path of the setup file in the app package.
- `size` (Number) The total size of the application in bytes.
- `superseded_app_count` (Number) The total number of apps this app is directly or indirectly superseded by.
- `superseding_app_count` (Number) The total number of apps this app directly or indirectly supersedes.
- `uninstall_command_line` (String) The command line to uninstall this app.
- `upload_state` (Number) The upload state.

<a id="nestedatt--items--install_experience"></a>
### Nested Schema for `items.install_experience`

Read-Only:

- `device_restart_behavior` (String) Device restart behavior.
- `max_run_time_in_minutes` (Number) The maximum run time in minutes.
- `run_as_account` (String) Indicates the account context to execute the app.


<a id="nestedatt--items--msi_information"></a>
### Nested Schema for `items.msi_information`

Read-Only:

- `package_type` (String) The MSI package type.
- `product_code` (String) The MSI product code.
- `product_name` (String) The MSI product name.
- `product_version` (String) The MSI product version.
- `publisher` (String) The MSI publisher.
- `requires_reboot` (Boolean) Whether the MSI app requires reboot.
- `upgrade_code` (String) The MSI upgrade code.


<a id="nestedatt--items--return_codes"></a>
### Nested Schema for `items.return_codes`

Read-Only:

- `return_code` (Number) The return code.
- `type` (String) The type of return code.


<a id="nestedatt--items--rules"></a>
### Nested Schema for `items.rules`

Read-Only:

- `check_32bit_on_64system` (Boolean) Indicates whether to check 32-bit on a 64-bit system.
- `comparison_value` (String) The comparison value for the rule.
- `file_or_folder_name` (String) The file or folder name to detect.
- `key_path` (String) The registry key path for registry rules.
- `odata_type` (String) The OData type of the rule.
- `operation_type` (String) The operation type for the rule.
- `operator` (String) The operator for the rule.
- `path` (String) The file or folder path for file system rules.
- `rule_type` (String) The type of rule (detection or requirement).
- `value_name` (String) The registry value name for registry rules.
