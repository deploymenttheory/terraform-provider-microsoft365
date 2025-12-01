---
page_title: "microsoft365_utility_licensing_service_plan_reference Data Source - terraform-provider-microsoft365"
subcategory: "Utility"

description: |-
  Queries Microsoft 365 licensing service plan reference data. This utility data source allows you to search for license products (SKUs) and service plans using human-readable names, GUIDs, or string IDs. The data is sourced from Microsoft's official licensing service plan reference documentation.
  Search Modes:
  By Product Name: Use product_name to search for license products (e.g., "Microsoft 365 E3")By Product Identifier: Use string_id or guid to look up a specific productBy Service Plan Name: Use service_plan_name to find service plans (e.g., "Exchange Online")By Service Plan Identifier: Use service_plan_id or service_plan_guid for specific service plans
  Only one search parameter should be specified at a time. Results include both the matching items and their relationships (e.g., which products include a specific service plan, or which service plans are included in a product).
  Reference: Microsoft Licensing Service Plan Reference https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference
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

## Argument Reference

**Note:** Exactly one of the following search parameters must be specified:

### Product Search Parameters

* `product_name` - (Optional) Search for products by name (case-insensitive partial match). Returns all products whose names contain this string. Example: `"Microsoft 365 E3"` or `"Office 365"`.

* `string_id` - (Optional) Look up a product by its String ID (exact match, case-insensitive). String IDs are used by PowerShell v1.0 and the `skuPartNumber` property in Microsoft Graph. Example: `"ENTERPRISEPACK"`, `"SPE_E5"`.

* `guid` - (Optional) Look up a product by its GUID (exact match). GUIDs are used by the `skuId` property in Microsoft Graph. Must be in the format `00000000-0000-0000-0000-000000000000`.

### Service Plan Search Parameters

* `service_plan_name` - (Optional) Search for service plans by friendly name (case-insensitive partial match). Returns all service plans whose names contain this string. Example: `"Exchange Online"`, `"Microsoft Teams"`.

* `service_plan_id` - (Optional) Search for service plans by ID (case-insensitive partial match). Returns all service plans whose IDs contain this string. Example: `"EXCHANGE"`, `"TEAMS"`.

* `service_plan_guid` - (Optional) Look up a service plan by its GUID (exact match). Must be in the format `00000000-0000-0000-0000-000000000000`.

* `timeouts` - (Optional) Timeout configuration block. See [Timeouts](#timeouts) below.

## Attributes Reference

* `id` - The computed ID of this datasource operation in the format `{search_type}:{search_value}`.

* `matching_products` - (Computed) List of products matching the search criteria. Populated when searching by `product_name`, `string_id`, or `guid`. Each product contains:
  - `product_name` - The product name as displayed in management portals
  - `string_id` - The product String ID used by PowerShell and Graph API
  - `guid` - The product GUID used by Graph API
  - `service_plans_included` - List of service plans included in this product:
    - `id` - Service plan ID
    - `name` - Service plan friendly name
    - `guid` - Service plan GUID

* `matching_service_plans` - (Computed) List of service plans matching the search criteria. Populated when searching by `service_plan_name`, `service_plan_id`, or `service_plan_guid`. Each service plan contains:
  - `id` - Service plan ID
  - `name` - Service plan friendly name
  - `guid` - Service plan GUID
  - `included_in_skus` - List of products (SKUs) that include this service plan:
    - `product_name` - Product name
    - `string_id` - Product String ID
    - `guid` - Product GUID

## Timeouts

The `timeouts` block supports:

* `read` - (Optional) Timeout for reading data. Defaults to 3 minutes. Note: This datasource performs local lookups against embedded data, so operations are typically instantaneous.

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `guid` (String) Look up a product by its GUID (exact match). GUIDs are used by the `skuId` property in Microsoft Graph. Example: `"6fd2c87f-b296-42f0-b197-1e91e994b900"`.
- `product_name` (String) Search for products by name (case-insensitive partial match). Returns all products whose names contain this string. Example: `"Microsoft 365 E3"` or `"Office 365"`.
- `service_plan_guid` (String) Look up a service plan by its GUID (exact match). Example: `"113feb6c-3fe4-4440-bddc-54d774bf0318"`.
- `service_plan_id` (String) Search for service plans by ID (case-insensitive partial match). Returns all service plans whose IDs contain this string. Example: `"EXCHANGE"`, `"TEAMS"`.
- `service_plan_name` (String) Search for service plans by friendly name (case-insensitive partial match). Returns all service plans whose names contain this string. Example: `"Exchange Online"`, `"Microsoft Teams"`.
- `string_id` (String) Look up a product by its String ID (exact match, case-insensitive). String IDs are used by PowerShell v1.0 and the `skuPartNumber` property in Microsoft Graph. Example: `"ENTERPRISEPACK"`, `"SPE_E3"`.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The ID of this data source operation.
- `matching_products` (Attributes List) List of products matching the search criteria. Populated when searching by product name, string_id, guid, or when searching for service plans (returns products that include the matching service plans). (see [below for nested schema](#nestedatt--matching_products))
- `matching_service_plans` (Attributes List) List of service plans matching the search criteria. Populated when searching by service plan name, id, or guid. Each entry includes a list of products (SKUs) that include this service plan. (see [below for nested schema](#nestedatt--matching_service_plans))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--matching_products"></a>
### Nested Schema for `matching_products`

Read-Only:

- `guid` (String) The product GUID used by skuId in Graph API.
- `product_name` (String) The product name as displayed in management portals (e.g., "Microsoft 365 E3").
- `service_plans_included` (Attributes List) Service plans included in this product. (see [below for nested schema](#nestedatt--matching_products--service_plans_included))
- `string_id` (String) The product String ID used by PowerShell v1.0 and skuPartNumber in Graph API.

<a id="nestedatt--matching_products--service_plans_included"></a>
### Nested Schema for `matching_products.service_plans_included`

Read-Only:

- `guid` (String) The service plan GUID.
- `id` (String) The service plan ID.
- `name` (String) The service plan friendly name.



<a id="nestedatt--matching_service_plans"></a>
### Nested Schema for `matching_service_plans`

Read-Only:

- `guid` (String) The service plan GUID.
- `id` (String) The service plan ID.
- `included_in_skus` (Attributes List) List of products (SKUs) that include this service plan. (see [below for nested schema](#nestedatt--matching_service_plans--included_in_skus))
- `name` (String) The service plan friendly name.

<a id="nestedatt--matching_service_plans--included_in_skus"></a>
### Nested Schema for `matching_service_plans.included_in_skus`

Read-Only:

- `guid` (String) The product GUID.
- `product_name` (String) The product name.
- `string_id` (String) The product String ID.

## Use Cases

This datasource supports multiple licensing automation scenarios:

1. **Dynamic License Assignment** - Look up license SKU GUIDs by human-readable names instead of hardcoding GUIDs
2. **License Discovery** - Find all available licenses that include specific service plans (e.g., all licenses with Exchange Online)
3. **Service Plan Analysis** - Understand which services are included in a license product
4. **Configuration Validation** - Verify that licenses in your tenant match expected String IDs or GUIDs
5. **Documentation Generation** - Dynamically generate documentation of available licenses and their capabilities

## Best Practices

1. **Use Product Names**: Search by `product_name` for better readability and maintainability (e.g., `"Microsoft 365 E3"` instead of `"ENTERPRISEPACK"`)
2. **Validate Results**: Check that `matching_products` or `matching_service_plans` contains at least one result before accessing `[0]`
3. **Handle Multiple Matches**: Product name searches may return multiple results (e.g., "Microsoft 365 E3" returns both standard and "no Teams" variants)
4. **Use Locals**: Store looked-up GUIDs in local values for reuse across multiple resources
5. **Document Your Choice**: Comment why you're using a specific license (e.g., "Microsoft 365 E3 (no Teams) to avoid Teams conflicts")

## Important Notes

### No Authentication Required
This datasource performs local lookups against embedded reference data and doesn't make any API calls. Unlike other datasources in this provider, it doesn't require Microsoft Graph API credentials.

### Data Freshness
The embedded licensing data is automatically updated weekly by a GitHub Actions workflow that fetches the latest information from Microsoft's documentation. The data source always reflects the state of the embedded JSON file at provider build time.

### Partial Matching
When searching by `product_name`, `service_plan_name`, or `service_plan_id`, the datasource performs case-insensitive partial matching. This means `"Microsoft 365 E3"` will match:
- "Microsoft 365 E3"
- "Microsoft 365 E3 (no Teams)"
- "Microsoft 365 E3 - Unattended License"

Use exact identifiers (`string_id` or `guid`) when you need precise matching.

### Multiple Results
Some searches may return multiple results. Always check the count and access specific array indices appropriately. Consider using filters or selecting the appropriate match based on your requirements.

## Common Licenses

Here are some commonly used licenses that can be looked up with this datasource:

| Product Name | String ID | Common Use Case |
|--------------|-----------|-----------------|
| Microsoft 365 E3 | ENTERPRISEPACK | Enterprise productivity suite |
| Microsoft 365 E3 (no Teams) | Microsoft_365_E3_(no_Teams) | Enterprise suite without Teams |
| Microsoft 365 E5 | ENTERPRISEPREMIUM | Enterprise with advanced security |
| Microsoft 365 Business Premium | SPB | Small/medium business suite |
| Exchange Online Plan 1 | EXCHANGESTANDARD | Email only |
| Exchange Online Plan 2 | EXCHANGE_S_ENTERPRISE | Email with advanced features |

## Data Updates

The licensing reference data is maintained through an automated pipeline:

- **Update Frequency**: Weekly (every Monday at 2 AM UTC)
- **Source**: [Microsoft Licensing Service Plan Reference](https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference)
- **Process**: Automated PR created when changes are detected
- **Review**: Changes are reviewed before merging into the provider

To manually update the data, run:
```bash
python3 scripts/pipeline/get_licensing_service_plan_reference.py \
  -f json \
  -o internal/services/datasources/utility/licensing_service_plan_reference/data/licensing_service_plan_reference.json
```

## Additional Resources

- [Microsoft Licensing Service Plan Reference](https://learn.microsoft.com/en-us/entra/identity/users/licensing-service-plan-reference)
- [Microsoft Graph API - subscribedSku Resource Type](https://learn.microsoft.com/en-us/graph/api/resources/subscribedsku)
- [Assign Licenses to Users (Microsoft Graph)](https://learn.microsoft.com/en-us/graph/api/user-assignlicense)
- [Group-Based Licensing in Microsoft Entra ID](https://learn.microsoft.com/en-us/entra/identity/users/licensing-groups-assign)


