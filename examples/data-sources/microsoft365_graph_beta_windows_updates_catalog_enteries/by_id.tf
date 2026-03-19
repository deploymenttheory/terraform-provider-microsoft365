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
