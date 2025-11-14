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
      id                          = package.id
      product_id                  = package.product_id
      product_display_name        = package.product_display_name
      publisher_display_name      = package.publisher_display_name
      version_display_name        = package.version_display_name
      branch_display_name         = package.branch_display_name
      applicable_architectures    = package.applicable_architectures
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

# Example 5: Get packages using OData filter with startswith()
# Note: The mobileAppCatalogPackages endpoint has limited OData support
# Working: $filter with startswith(), $top
# Not working: $orderby, $count, $skip, $select, $search when combined with $filter
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_filter" {
  filter_type  = "odata"
  odata_filter = "startswith(publisherDisplayName, 'Microsoft')"
  odata_top    = 10

  timeouts = {
    read = "30s"
  }
}

# Output for OData filter
output "packages_odata_filter" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_filter.items : []
}

# Example 6: OData query with different publisher filter
# WARNING: Do not use $skip or $select as they cause 500 errors/timeouts
# Do not combine $filter with $orderby or $count as they return no results
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_advanced" {
  filter_type  = "odata"
  odata_filter = "startswith(publisherDisplayName, 'Google')"
  odata_top    = 5

  timeouts = {
    read = "30s"
  }
}

# Output for advanced OData query
output "packages_odata_advanced" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_advanced.items : []
}

# Example 7: Filter by product display name
# Note: Use startswith() instead of $search which is not reliably supported
data "microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package" "odata_by_product" {
  filter_type  = "odata"
  odata_filter = "startswith(productDisplayName, 'Microsoft Edge')"
  odata_top    = 10

  timeouts = {
    read = "30s"
  }
}

# Output for product filter query
output "packages_by_product" {
  value = data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_by_product.items != null ? data.microsoft365_graph_beta_device_and_app_management_mobile_app_catalog_package.odata_by_product.items : []
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

- `items` (Attributes List) The list of mobile app catalog packages that match the filter criteria. (see [below for nested schema](#nestedatt--items))

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

- `applicable_architectures` (String) The applicable architectures for the package (e.g., x64, x86, ARM64).
- `branch_display_name` (String) The display name of the branch.
- `id` (String) The unique identifier for the mobile app catalog package.
- `locales` (List of String) The list of supported locales for the package.
- `package_auto_update_capable` (Boolean) Indicates whether the package supports automatic updates.
- `product_display_name` (String) The display name of the product.
- `product_id` (String) The unique identifier for the product.
- `publisher_display_name` (String) The display name of the publisher.
- `version_display_name` (String) The display name of the version.
