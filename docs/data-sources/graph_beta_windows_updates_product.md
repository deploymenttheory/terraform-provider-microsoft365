---
page_title: "microsoft365_graph_beta_windows_updates_product Data Source - terraform-provider-microsoft365"
subcategory: "Windows Updates"

description: |-
  Retrieves Windows Update product information from Microsoft Graph. This data source can search by catalog ID or KB number using the /admin/windows/updates/products/FindByCatalogId or /admin/windows/updates/products/FindByKbNumber endpoints.
---

# microsoft365_graph_beta_windows_updates_product (Data Source)

Retrieves Windows Update product information from Microsoft Graph. This data source can search by catalog ID or KB number using the `/admin/windows/updates/products/FindByCatalogId` or `/admin/windows/updates/products/FindByKbNumber` endpoints.

## Microsoft Documentation

- [product resource type](https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-product?view=graph-rest-beta)
- [product: findByCatalogId](https://learn.microsoft.com/en-us/graph/api/windowsupdates-product-findbycatalogid?view=graph-rest-beta)
- [product: findByKbNumber](https://learn.microsoft.com/en-us/graph/api/windowsupdates-product-findbykbnumber?view=graph-rest-beta)

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

### Find by Catalog ID

```terraform
# Find Windows Update product by catalog ID
# This example shows how to get a catalog ID from catalog entries,
# then use it to retrieve detailed product information including revisions and known issues

# Step 1: Get a specific catalog entry (e.g., latest quality update)
data "microsoft365_graph_beta_windows_updates_catalog_enteries" "latest_quality_update" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

# Step 2: Use the catalog ID to get detailed product information
data "microsoft365_graph_beta_windows_updates_product" "by_catalog_id" {
  search_type  = "catalog_id"
  search_value = data.microsoft365_graph_beta_windows_updates_catalog_enteries.latest_quality_update.entries[0].id
}

output "product_info" {
  description = "Product information retrieved by catalog ID"
  value = length(data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products) > 0 ? {
    id           = data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products[0].id
    name         = data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products[0].name
    group_name   = data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products[0].group_name
    revisions    = length(data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products[0].revisions)
    known_issues = length(data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products[0].known_issues)
  } : null
}

output "friendly_names" {
  description = "Friendly names for the product"
  value       = length(data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products) > 0 ? data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products[0].friendly_names : []
}

output "all_revisions" {
  description = "All revisions for the product with OS build details"
  value = length(data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products) > 0 ? [
    for revision in data.microsoft365_graph_beta_windows_updates_product.by_catalog_id.products[0].revisions : {
      id           = revision.id
      display_name = revision.display_name
      version      = revision.version
      os_build = try({
        major_version         = revision.os_build.major_version
        minor_version         = revision.os_build.minor_version
        build_number          = revision.os_build.build_number
        update_build_revision = revision.os_build.update_build_revision
      }, null)
    }
  ] : []
}
```

### Find by KB Number

```terraform
# Find Windows Update product by KB number
# This example shows how to search for a product using a KB article number
# KB numbers are typically known from Microsoft security bulletins or support articles

data "microsoft365_graph_beta_windows_updates_product" "by_kb_number" {
  search_type  = "kb_number"
  search_value = "5029332" # Example: KB5029332
}

output "product_by_kb" {
  description = "Product information retrieved by KB number"
  value = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products) > 0 ? {
    id           = data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].id
    name         = data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].name
    group_name   = data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].group_name
    revisions    = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].revisions)
    known_issues = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].known_issues)
  } : null
}

output "kb_articles" {
  description = "Knowledge Base articles for all revisions"
  value = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products) > 0 ? [
    for revision in data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].revisions : {
      revision_id = revision.id
      kb_id       = try(revision.knowledge_base_article.id, null)
      kb_url      = try(revision.knowledge_base_article.url, null)
    }
  ] : []
}

output "known_issues" {
  description = "All known issues for the product"
  value = length(data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products) > 0 ? [
    for issue in data.microsoft365_graph_beta_windows_updates_product.by_kb_number.products[0].known_issues : {
      id                 = issue.id
      title              = issue.title
      status             = issue.status
      description        = issue.description
      web_view_url       = issue.web_view_url
      start_date_time    = issue.start_date_time
      resolved_date_time = issue.resolved_date_time
    }
  ] : []
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `search_type` (String) Type of search to perform. Valid values are: `catalog_id`, `kb_number`.
- `search_value` (String) Value to search by. For catalog_id, provide the catalog identifier. For kb_number, provide the KB article number (e.g., '5029332').

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `products` (Attributes List) The list of Windows Update products that match the search criteria. (see [below for nested schema](#nestedatt--products))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).


<a id="nestedatt--products"></a>
### Nested Schema for `products`

Read-Only:

- `friendly_names` (List of String) The friendly names of the product, e.g., 'Version 22H2 (OS build 22621)'.
- `group_name` (String) The name of the product group, e.g., 'Windows 11'.
- `id` (String) The unique identifier for the product.
- `known_issues` (Attributes List) Known issues related to the product. (see [below for nested schema](#nestedatt--products--known_issues))
- `name` (String) The name of the product, e.g., 'Windows 11, version 22H2'.
- `revisions` (Attributes List) Product revisions associated with the search criteria. (see [below for nested schema](#nestedatt--products--revisions))

<a id="nestedatt--products--known_issues"></a>
### Nested Schema for `products.known_issues`

Read-Only:

- `description` (String) The description of the known issue.
- `id` (String) The unique identifier for the known issue.
- `last_updated_date_time` (String) The last updated date and time of the known issue in RFC3339 format.
- `originating_knowledge_base_article` (Attributes) The KB article that originated the known issue. (see [below for nested schema](#nestedatt--products--known_issues--originating_knowledge_base_article))
- `resolved_date_time` (String) The resolved date and time of the known issue in RFC3339 format.
- `resolving_knowledge_base_article` (Attributes) The KB article that resolved the known issue. (see [below for nested schema](#nestedatt--products--known_issues--resolving_knowledge_base_article))
- `safeguard_hold_ids` (List of String) List of safeguard hold IDs associated with the known issue.
- `start_date_time` (String) The start date and time of the known issue in RFC3339 format.
- `status` (String) The status of the known issue.
- `title` (String) The title of the known issue.
- `web_view_url` (String) The URL to view the known issue in the admin portal.

<a id="nestedatt--products--known_issues--originating_knowledge_base_article"></a>
### Nested Schema for `products.known_issues.originating_knowledge_base_article`

Read-Only:

- `id` (String) The KB article ID.
- `url` (String) The URL to the KB article.


<a id="nestedatt--products--known_issues--resolving_knowledge_base_article"></a>
### Nested Schema for `products.known_issues.resolving_knowledge_base_article`

Read-Only:

- `id` (String) The KB article ID.
- `url` (String) The URL to the KB article.



<a id="nestedatt--products--revisions"></a>
### Nested Schema for `products.revisions`

Read-Only:

- `catalog_entry` (Attributes) The catalog entry associated with this revision. (see [below for nested schema](#nestedatt--products--revisions--catalog_entry))
- `display_name` (String) The display name of the product revision.
- `id` (String) The unique identifier for the product revision, e.g., '10.0.22621.2215'.
- `knowledge_base_article` (Attributes) The knowledge base article associated with this revision. (see [below for nested schema](#nestedatt--products--revisions--knowledge_base_article))
- `os_build` (Attributes) The OS build information. (see [below for nested schema](#nestedatt--products--revisions--os_build))
- `release_date_time` (String) The release date and time in RFC3339 format.
- `version` (String) The version of the product revision, e.g., '22H2'.

<a id="nestedatt--products--revisions--catalog_entry"></a>
### Nested Schema for `products.revisions.catalog_entry`

Read-Only:

- `catalog_name` (String) The catalog name.
- `deployable_until_date_time` (String) The date and time until which the update can be deployed, in RFC3339 format.
- `display_name` (String) The display name of the catalog entry.
- `id` (String) The catalog entry identifier.
- `is_expeditable` (Boolean) Indicates whether the update can be expedited.
- `quality_update_cadence` (String) The release cadence of the quality update.
- `quality_update_classification` (String) The classification of the quality update.
- `release_date_time` (String) The release date and time in RFC3339 format.
- `short_name` (String) The short name of the update.


<a id="nestedatt--products--revisions--knowledge_base_article"></a>
### Nested Schema for `products.revisions.knowledge_base_article`

Read-Only:

- `id` (String) The KB article ID, e.g., 'KB5029351'.
- `url` (String) The URL to the KB article.


<a id="nestedatt--products--revisions--os_build"></a>
### Nested Schema for `products.revisions.os_build`

Read-Only:

- `build_number` (Number) The build number.
- `major_version` (Number) The major version number.
- `minor_version` (Number) The minor version number.
- `update_build_revision` (Number) The update build revision number.
