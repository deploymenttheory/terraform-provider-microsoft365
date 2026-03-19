---
page_title: "microsoft365_graph_beta_windows_updates_catalog_enteries Data Source - terraform-provider-microsoft365"
subcategory: "Windows Updates"

description: |-
  Retrieves Windows Update catalog entries from Microsoft Graph using the /admin/windows/updates/catalog/entries endpoint. This data source returns feature update and quality update catalog entries with deployment information.
---

# microsoft365_graph_beta_windows_updates_catalog_enteries (Data Source)

Retrieves Windows Update catalog entries from Microsoft Graph using the `/admin/windows/updates/catalog/entries` endpoint. This data source returns feature update and quality update catalog entries with deployment information.

## Microsoft Documentation

- [catalog resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-catalog?view=graph-rest-beta)
- [List catalog entries](https://learn.microsoft.com/en-us/graph/api/windowsupdates-catalog-list-entries?view=graph-rest-beta)
- [featureUpdateCatalogEntry resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-featureupdatecatalogentry?view=graph-rest-beta)
- [qualityUpdateCatalogEntry resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-qualityupdatecatalogentry?view=graph-rest-beta)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this data source:

**Required:**
- `WindowsUpdates.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.50.0-alpha | Experimental | Initial release |

## Example Usage

### Get All Windows Update Catalog Entries

```terraform
# Get all Windows Update catalog entries
# This retrieves all available catalog entries (both feature and quality updates)

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "all" {
  filter_type = "all"
}

output "total_entries" {
  description = "Total number of catalog entries available"
  value       = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.all.entries)
}

output "entry_types" {
  description = "Distinct catalog entry types found"
  value = distinct([
    for entry in data.microsoft365_graph_beta_windows_updates_catalog_enteries.all.entries :
    entry.catalog_entry_type
  ])
}

output "recent_updates" {
  description = "The 5 most recent updates"
  value = [
    for entry in slice(data.microsoft365_graph_beta_windows_updates_catalog_enteries.all.entries, 0, min(5, length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.all.entries))) : {
      id                 = entry.id
      display_name       = entry.display_name
      catalog_entry_type = entry.catalog_entry_type
      release_date_time  = entry.release_date_time
    }
  ]
}
```

### Get Feature Updates Only

```terraform
# Get feature updates only
# This retrieves only feature update catalog entries (e.g., Windows 11 22H2, 23H2)

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "feature_updates" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

output "feature_update_count" {
  description = "Number of feature updates available"
  value       = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_updates.entries)
}

output "available_versions" {
  description = "All available Windows feature update versions"
  value = [
    for entry in data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_updates.entries :
    {
      id           = entry.id
      version      = entry.version
      display_name = entry.display_name
      release_date = entry.release_date_time
    }
  ]
}

output "latest_feature_update" {
  description = "The most recent feature update"
  value = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_updates.entries) > 0 ? {
    id           = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_updates.entries[0].id
    version      = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_updates.entries[0].version
    display_name = data.microsoft365_graph_beta_windows_updates_catalog_enteries.feature_updates.entries[0].display_name
  } : null
}
```

### Get Quality Updates Only

```terraform
# Get quality updates only
# This retrieves only quality update catalog entries (security and non-security updates)

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_updates" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

output "quality_update_count" {
  description = "Number of quality updates available"
  value       = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries)
}

output "security_updates" {
  description = "All security quality updates"
  value = [
    for entry in data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries :
    {
      id                            = entry.id
      display_name                  = entry.display_name
      short_name                    = entry.short_name
      catalog_name                  = entry.catalog_name
      quality_update_classification = entry.quality_update_classification
      is_expeditable                = entry.is_expeditable
      release_date_time             = entry.release_date_time
    } if entry.quality_update_classification == "security"
  ]
}

output "latest_quality_update" {
  description = "The most recent quality update with CVE information"
  value = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries) > 0 ? {
    id                            = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].id
    display_name                  = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].display_name
    short_name                    = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].short_name
    quality_update_classification = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].quality_update_classification
    is_expeditable                = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].is_expeditable
    cve_max_severity              = try(data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].cve_severity_information.max_severity, null)
    cve_max_base_score            = try(data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].cve_severity_information.max_base_score, null)
  } : null
}
```

### Filter by Windows Update Catalog Entry Display Name

```terraform
# Filter catalog entries by display name
# This searches for catalog entries containing specific text in their display name

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "windows_11_updates" {
  filter_type  = "display_name"
  filter_value = "Windows 11"
}

output "matching_entries_count" {
  description = "Number of entries matching 'Windows 11'"
  value       = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.windows_11_updates.entries)
}

output "matching_updates" {
  description = "All catalog entries matching the display name filter"
  value = [
    for entry in data.microsoft365_graph_beta_windows_updates_catalog_enteries.windows_11_updates.entries : {
      id                 = entry.id
      display_name       = entry.display_name
      catalog_entry_type = entry.catalog_entry_type
      release_date_time  = entry.release_date_time
      version            = entry.catalog_entry_type == "featureUpdate" ? entry.version : null
      short_name         = entry.catalog_entry_type == "qualityUpdate" ? entry.short_name : null
    }
  ]
}
```

### Get Specific Windows Update Catalog Entry by ID

```terraform
# Get a specific catalog entry by ID
# This retrieves a single catalog entry using its unique identifier

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "specific_update" {
  filter_type  = "id"
  filter_value = "c1dec151-c151-c1de-51c1-dec151c1dec1"
}

output "update_details" {
  description = "Details of the specific catalog entry"
  value = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.specific_update.entries) > 0 ? {
    id                         = data.microsoft365_graph_beta_windows_updates_catalog_enteries.specific_update.entries[0].id
    display_name               = data.microsoft365_graph_beta_windows_updates_catalog_enteries.specific_update.entries[0].display_name
    catalog_entry_type         = data.microsoft365_graph_beta_windows_updates_catalog_enteries.specific_update.entries[0].catalog_entry_type
    release_date_time          = data.microsoft365_graph_beta_windows_updates_catalog_enteries.specific_update.entries[0].release_date_time
    deployable_until_date_time = data.microsoft365_graph_beta_windows_updates_catalog_enteries.specific_update.entries[0].deployable_until_date_time
  } : null
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter_type` (String) Type of filter to apply. Valid values are: `all`, `id`, `display_name`, `catalog_entry_type`.

### Optional

- `filter_value` (String) Value to filter by. Not required when filter_type is 'all'. For catalog_entry_type, use 'featureUpdate' or 'qualityUpdate'.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `entries` (Attributes List) The list of Windows Update Catalog Entries that match the filter criteria. (see [below for nested schema](#nestedatt--entries))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--entries"></a>
### Nested Schema for `entries`

Read-Only:

- `catalog_entry_type` (String) The type of catalog entry. Values: 'featureUpdate' or 'qualityUpdate'.
- `catalog_name` (String) The catalog name of the quality update (quality updates only).
- `cve_severity_information` (Attributes) CVE severity information for the quality update (quality updates only). (see [below for nested schema](#nestedatt--entries--cve_severity_information))
- `deployable_until_date_time` (String) The date and time until which the update can be deployed, in RFC3339 format. Null if no expiration.
- `display_name` (String) The display name of the catalog entry.
- `id` (String) The unique identifier for the catalog entry.
- `is_expeditable` (Boolean) Indicates whether the quality update can be expedited (quality updates only).
- `quality_update_cadence` (String) The release cadence of the quality update, e.g., 'monthly' (quality updates only).
- `quality_update_classification` (String) The classification of the quality update, e.g., 'security' (quality updates only).
- `release_date_time` (String) The release date and time of the catalog entry in RFC3339 format.
- `short_name` (String) The short name of the quality update (quality updates only).
- `version` (String) The version of the feature update (feature updates only).

<a id="nestedatt--entries--cve_severity_information"></a>
### Nested Schema for `entries.cve_severity_information`

Read-Only:

- `exploited_cves` (Attributes List) List of exploited CVEs. (see [below for nested schema](#nestedatt--entries--cve_severity_information--exploited_cves))
- `max_base_score` (Number) The maximum CVSS base score.
- `max_severity` (String) The maximum severity level of CVEs, e.g., 'critical'.

<a id="nestedatt--entries--cve_severity_information--exploited_cves"></a>
### Nested Schema for `entries.cve_severity_information.exploited_cves`

Read-Only:

- `number` (String) The CVE number, e.g., 'CVE-2023-32046'.
- `url` (String) The URL to the CVE details.
