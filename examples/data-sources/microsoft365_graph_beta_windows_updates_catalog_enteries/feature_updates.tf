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
