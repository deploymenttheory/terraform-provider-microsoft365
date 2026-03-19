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
