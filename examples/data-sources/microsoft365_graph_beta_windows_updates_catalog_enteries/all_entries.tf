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
