# Test 01: Get all Windows Update catalog entries
# This test retrieves all available catalog entries from the Windows Update service

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "test" {
  filter_type = "all"
}

output "total_entries" {
  value = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.test.entries)
}

output "entry_types" {
  value = distinct([
    for entry in data.microsoft365_graph_beta_windows_updates_catalog_enteries.test.entries :
    entry.catalog_entry_type
  ])
}
