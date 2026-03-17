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
