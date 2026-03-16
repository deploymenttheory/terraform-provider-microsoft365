# Test 04: Filter by display name
# This test searches for catalog entries by display name

data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "test" {
  filter_type  = "display_name"
  filter_value = "Windows 11"
}

output "matching_entries" {
  value = length(data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.test.entries)
}

output "entry_names" {
  value = [
    for entry in data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.test.entries :
    entry.display_name
  ]
}
