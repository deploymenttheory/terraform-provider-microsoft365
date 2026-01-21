---
page_title: "microsoft365_utility_licensing_service_plan_reference Data Source - terraform-provider-microsoft365"
subcategory: "Utility"

description: |-
  Queries Microsoft 365 licensing service plan reference data from the embedded licensing database. This data source is used to look up license SKU and service plan GUIDs by human-readable names for assignment configuration.
---

# microsoft365_utility_licensing_service_plan_reference

Queries Microsoft 365 licensing service plan reference data to find license products (SKUs) and service plans using human-readable names, GUIDs, or string IDs.

This utility datasource provides access to Microsoft's official licensing reference data, allowing you to dynamically look up license SKU GUIDs without hardcoding them in your Terraform configurations. This is particularly useful for license assignment operations where you need the `skuId` but want to use human-readable identifiers.

## Background

Microsoft 365 licenses are identified in three ways:

- **Product Name**: The human-readable name displayed in admin portals (e.g., "Microsoft 365 E3")
- **String ID**: Used by PowerShell v1.0 cmdlets and the `skuPartNumber` property in Microsoft Graph API (e.g., "ENTERPRISEPACK")
- **GUID**: Used by the `skuId` property in Microsoft Graph API (e.g., "6fd2c87f-b296-42f0-b197-1e91e994b900")

Each license product includes multiple service plans (individual Microsoft 365 services like Exchange Online, Teams, SharePoint, etc.). This datasource allows you to:

1. Search for products by name and retrieve their identifiers
2. Look up specific products by String ID or GUID
3. Find which products include specific service plans
4. Discover relationships between licenses and their included services

The data is automatically kept up-to-date through a weekly automated pipeline that fetches the latest information from Microsoft's official documentation.

## Search Modes

The datasource supports six search modes. **Exactly one** search parameter must be specified:

### Product Search
- **`product_name`**: Case-insensitive partial match (e.g., "Microsoft 365 E3" finds all E3 variants)
- **`string_id`**: Exact match by String ID (e.g., "ENTERPRISEPACK")
- **`guid`**: Exact match by Product GUID

### Service Plan Search
- **`service_plan_name`**: Case-insensitive partial match (e.g., "Exchange Online")
- **`service_plan_id`**: Case-insensitive partial match (e.g., "EXCHANGE")
- **`service_plan_guid`**: Exact match by Service Plan GUID

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.37.0-alpha | Experimental | Initial release |

## Example Usage

### Basic License Lookup

```terraform
# Example: Search for Microsoft 365 E3 license by product name
data "microsoft365_utility_licensing_service_plan_reference" "m365_e3" {
  product_name = "Microsoft 365 E3"
}

# Output the matching product details
output "m365_e3_details" {
  value = {
    product_name = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].product_name
    string_id    = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].string_id
    guid         = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].guid
  }
}

# Output the service plans included in Microsoft 365 E3
output "m365_e3_service_plans" {
  value = data.microsoft365_utility_licensing_service_plan_reference.m365_e3.matching_products[0].service_plans_included
}
```

### License Assignment with Dynamic Lookup

```terraform
# Example: Look up a license by its String ID
# String IDs are used in PowerShell v1.0 and the skuPartNumber property in Microsoft Graph
data "microsoft365_utility_licensing_service_plan_reference" "enterprise_pack" {
  string_id = "ENTERPRISEPACK" # Microsoft 365 E3
}

# Use the GUID in a license assignment
resource "microsoft365_graph_beta_users_user_license_assignment" "example" {
  user_id = "user-id-here"
  sku_id  = data.microsoft365_utility_licensing_service_plan_reference.enterprise_pack.matching_products[0].guid
}
```

### Find Products Containing Specific Service Plans

```terraform
# Example: Find which products include Exchange Online service plans
data "microsoft365_utility_licensing_service_plan_reference" "exchange_plans" {
  service_plan_name = "Exchange Online"
}

# Output the service plans and which products include them
output "exchange_service_plans" {
  value = [
    for plan in data.microsoft365_utility_licensing_service_plan_reference.exchange_plans.matching_service_plans : {
      service_plan_id   = plan.id
      service_plan_name = plan.name
      included_in       = [for sku in plan.included_in_skus : sku.product_name]
    }
  ]
}
```

### Practical Example: Complete User License Assignment

```terraform
# Practical Example: Dynamic license assignment using human-readable names
# This demonstrates how to avoid hardcoding GUIDs in your configurations

# Look up Microsoft 365 E3 (no Teams) by product name
data "microsoft365_utility_licensing_service_plan_reference" "m365_e3_no_teams" {
  product_name = "Microsoft 365 E3 (no Teams)"
}

# Create a user
resource "microsoft365_graph_beta_users_user" "example_user" {
  display_name        = "Example User"
  user_principal_name = "example.user@yourdomain.com"
  mail_nickname       = "example.user"
  account_enabled     = true
  usage_location      = "US"

  password_profile = {
    password                           = "SecureP@ssw0rd123!"
    force_change_password_next_sign_in = true
  }
}

# Assign the license using the dynamically looked-up GUID
resource "microsoft365_graph_beta_users_user_license_assignment" "example_license" {
  user_id = microsoft365_graph_beta_users_user.example_user.id

  # No hardcoded GUID - always up-to-date from Microsoft's reference data
  sku_id = data.microsoft365_utility_licensing_service_plan_reference.m365_e3_no_teams.matching_products[0].guid
}

# Output the license details for verification
output "assigned_license" {
  value = {
    product_name        = data.microsoft365_utility_licensing_service_plan_reference.m365_e3_no_teams.matching_products[0].product_name
    sku_id              = data.microsoft365_utility_licensing_service_plan_reference.m365_e3_no_teams.matching_products[0].guid
    service_plans_count = length(data.microsoft365_utility_licensing_service_plan_reference.m365_e3_no_teams.matching_products[0].service_plans_included)
  }
}
```

## Additional Resources

- [Microsoft Licensing Service Plan Reference](https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference)
- [Microsoft Graph API - subscribedSku Resource Type](https://learn.microsoft.com/en-us/graph/api/resources/subscribedsku)
- [Assign Licenses to Users (Microsoft Graph)](https://learn.microsoft.com/en-us/graph/api/user-assignlicense)
- [Group-Based Licensing in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity/users/licensing-groups-assign)


