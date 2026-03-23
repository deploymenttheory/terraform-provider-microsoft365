---
page_title: "microsoft365_graph_identity_and_access_subscribed_skus Data Source - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Retrieves Microsoft 365 license SKUs from Microsoft Entra ID using the /subscribedSkus endpoint. Supports flexible lookup by SKU ID, SKU part number, account ID, account name, or applies-to filter.
---

# microsoft365_graph_identity_and_access_subscribed_skus (Data Source)

Retrieves Microsoft 365 license SKUs from Microsoft Entra ID using the `/subscribedSkus` endpoint. Supports flexible lookup by SKU ID, SKU part number, account ID, account name, or applies-to filter.

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `Directory.Read.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.15.0-alpha | Experimental | Initial release |
| v0.35.0-alpha | Experimental | Renamed from graph_directory_management_subscribed_skus |
| v0.51.0-alpha | Experimental | Added `account_id` and `account_name` filters |

## Example Usage

### List All Subscribed SKUs

Retrieve all license SKUs available in your tenant.

```terraform
data "microsoft365_graph_identity_and_access_subscribed_skus" "all" {
  list_all = true

  timeouts = {
    read = "30s"
  }
}

output "all_skus_count" {
  value       = length(data.microsoft365_graph_identity_and_access_subscribed_skus.all.items)
  description = "Total number of subscribed SKUs"
}
```

### Get Specific SKU by ID

Query a specific SKU using its composite ID (format: `{accountId}_{skuId}`).

```terraform
data "microsoft365_graph_identity_and_access_subscribed_skus" "e5_sku" {
  sku_id = "2fd6bb84-ad40-4ec5-9369-a215b25c9952_06ebc4ee-1bb5-47dd-8120-11324bc54e06" // tenant id and sku id

  timeouts = {
    read = "30s"
  }
}

output "e5_sku_details" {
  value = length(data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items) > 0 ? {
    sku_part_number   = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].sku_part_number
    consumed_units    = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].consumed_units
    enabled_units     = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].prepaid_units.enabled
    capability_status = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].capability_status
    service_plans     = data.microsoft365_graph_identity_and_access_subscribed_skus.e5_sku.items[0].service_plans
  } : null
  description = "Microsoft 365 E5 SKU details including service plans"
}
```

### Filter by SKU Part Number

Find SKUs by part number using case-insensitive partial matching (e.g., "E5" matches "SPE_E5", "Microsoft_365_E3", etc.).

```terraform
data "microsoft365_graph_identity_and_access_subscribed_skus" "e5_by_part_number" {
  sku_part_number = "E5"

  timeouts = {
    read = "30s"
  }
}

output "e5_skus_summary" {
  value = [
    for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.e5_by_part_number.items : {
      sku_part_number    = sku.sku_part_number
      consumed_units     = sku.consumed_units
      enabled_units      = sku.prepaid_units.enabled
      available_licenses = sku.prepaid_units.enabled - sku.consumed_units
    }
  ]
  description = "All SKUs containing 'E5' in their part number"
}
```

### Filter by Account ID (Tenant ID)

Retrieve all SKUs for a specific account. The `account_id` typically matches your tenant ID.

```terraform
data "microsoft365_graph_identity_and_access_subscribed_skus" "by_account" {
  account_id = "f97aeefc-af85-414d-8ae4-b457f90efc40" // your tenant id

  timeouts = {
    read = "30s"
  }
}

output "tenant_license_summary" {
  value = {
    total_skus = length(data.microsoft365_graph_identity_and_access_subscribed_skus.tenant_skus.items)
    skus = [
      for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.tenant_skus.items : {
        name               = sku.sku_part_number
        total_licenses     = sku.prepaid_units.enabled
        used_licenses      = sku.consumed_units
        available_licenses = sku.prepaid_units.enabled - sku.consumed_units
        status             = sku.capability_status
      }
    ]
  }
  description = "Complete license summary for the tenant (account_id matches tenant_id)"
}
```

### Filter by Account Name

Find SKUs by account name using case-insensitive partial matching.

```terraform
data "microsoft365_graph_identity_and_access_subscribed_skus" "org_skus" {
  account_name = "DeploymentTheory" // your tenant name

  timeouts = {
    read = "30s"
  }
}

output "org_license_usage" {
  value = [
    for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.org_skus.items : {
      sku_part_number    = sku.sku_part_number
      account_name       = sku.account_name
      total_licenses     = sku.prepaid_units.enabled
      used_licenses      = sku.consumed_units
      available_licenses = sku.prepaid_units.enabled - sku.consumed_units
      utilization_pct    = sku.prepaid_units.enabled > 0 ? (sku.consumed_units / sku.prepaid_units.enabled) * 100 : 0
    }
  ]
  description = "License usage summary filtered by account name (partial match)"
}
```

### Filter by User-Assignable SKUs

Retrieve only SKUs that can be assigned to individual users (excludes company-level licenses).

```terraform
data "microsoft365_graph_identity_and_access_subscribed_skus" "user_assignable" {
  applies_to = "User"

  timeouts = {
    read = "30s"
  }
}

output "user_license_inventory" {
  value = [
    for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.user_assignable.items : {
      sku_part_number    = sku.sku_part_number
      sku_id             = sku.sku_id
      consumed_units     = sku.consumed_units
      enabled_units      = sku.prepaid_units.enabled
      available_units    = sku.prepaid_units.enabled - sku.consumed_units
      capability_status  = sku.capability_status
      service_plan_count = length(sku.service_plans)
    }
  ]
  description = "User-assignable SKUs with license availability"
}
```

### Filter by Company-Level SKUs

Retrieve only SKUs that apply at the company level (not assignable to individual users).

```terraform
data "microsoft365_graph_identity_and_access_subscribed_skus" "company_level" {
  applies_to = "Company"

  timeouts = {
    read = "30s"
  }
}

output "company_licenses" {
  value = [
    for sku in data.microsoft365_graph_identity_and_access_subscribed_skus.company_level.items : {
      sku_part_number   = sku.sku_part_number
      sku_id            = sku.sku_id
      account_name      = sku.account_name
      capability_status = sku.capability_status
      consumed_units    = sku.consumed_units
    }
  ]
  description = "Company-level SKUs (not assignable to individual users)"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `account_id` (String) Filter by account ID (GUID). This typically matches your tenant ID. Conflicts with other lookup attributes.
- `account_name` (String) Filter by account name. Case-insensitive partial match. Conflicts with other lookup attributes.
- `applies_to` (String) Filter by target class. Possible values: 'User', 'Company'. Conflicts with other lookup attributes.
- `list_all` (Boolean) Retrieve all subscribed SKUs. Conflicts with specific lookup attributes.
- `sku_id` (String) The unique identifier of a specific SKU (format: accountId_skuId). Conflicts with other lookup attributes.
- `sku_part_number` (String) Filter by SKU part number (e.g., 'ENTERPRISEPREMIUM', 'AAD_PREMIUM'). Case-insensitive partial match. Conflicts with other lookup attributes.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) The unique identifier for this data source operation.
- `items` (Attributes List) List of subscribed SKUs matching the query criteria. (see [below for nested schema](#nestedatt--items))

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

- `account_id` (String) The unique ID of the account this SKU belongs to.
- `account_name` (String) The name of the account this SKU belongs to.
- `applies_to` (String) The target class for this SKU. Only SKUs with target class 'User' are assignable. Possible values: 'User', 'Company'.
- `capability_status` (String) The status of the SKU. Possible values: 'Enabled', 'Warning', 'Suspended', 'Deleted', 'LockedOut'.
- `consumed_units` (Number) The number of licenses that have been assigned.
- `id` (String) The unique identifier for the subscribed SKU object.
- `prepaid_units` (Attributes) Information about the number and status of prepaid licenses. (see [below for nested schema](#nestedatt--items--prepaid_units))
- `service_plans` (Attributes List) Information about the service plans that are available with the SKU. (see [below for nested schema](#nestedatt--items--service_plans))
- `sku_id` (String) The unique identifier (GUID) for the service SKU.
- `sku_part_number` (String) The SKU part number; for example: 'AAD_PREMIUM' or 'RMSBASIC'.
- `subscription_ids` (List of String) A list of all subscription IDs associated with this SKU.

<a id="nestedatt--items--prepaid_units"></a>
### Nested Schema for `items.prepaid_units`

Read-Only:

- `enabled` (Number) The number of units that are enabled.
- `locked_out` (Number) The number of units that are locked out.
- `suspended` (Number) The number of units that are suspended.
- `warning` (Number) The number of units that are in warning state.


<a id="nestedatt--items--service_plans"></a>
### Nested Schema for `items.service_plans`

Read-Only:

- `applies_to` (String) The object the service plan can be assigned to.
- `provisioning_status` (String) The provisioning status of the service plan.
- `service_plan_id` (String) The unique identifier of the service plan.
- `service_plan_name` (String) The name of the service plan.

## Important Notes

- **License Information**: This data source retrieves information about commercial subscriptions that an organization has acquired using the [List subscribedSkus](https://learn.microsoft.com/en-us/graph/api/subscribedsku-list?view=graph-rest-1.0&tabs=http) Microsoft Graph API.
- **Flexible Lookup**: Supports multiple lookup methods including SKU ID, SKU part number, account ID, account name, applies to (User/Company), or list all.
- **Mutually Exclusive Filters**: Only one lookup method can be used at a time (enforced by schema validators).
- **Account ID**: The `account_id` attribute typically matches your Microsoft 365 tenant ID. All SKUs for a given tenant will share the same account ID.
- **SKU ID Format**: The `sku_id` is a composite identifier in the format `{accountId}_{skuId}` (e.g., `2fd6bb84-ad40-4ec5-9369-a215b25c9952_06ebc4ee-1bb5-47dd-8120-11324bc54e06`).
- **Local Filtering**: Since the Graph API v1.0 `/subscribedSkus` endpoint does not support OData filtering, all filters except `sku_id` perform local filtering after retrieving all SKUs.
- **License Usage**: The data includes consumed units vs prepaid units to help with license management.
- **Service Plans**: Each SKU includes detailed service plan information with provisioning status.
- **No User Permissions**: This provider only supports application permissions, not delegated user permissions.
