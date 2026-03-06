# Test 02: Get feature updates only
# This test retrieves only feature update catalog entries

data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "test" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

output "feature_update_count" {
  value = length(data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.test.entries)
}

output "feature_versions" {
  value = [
    for entry in data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.test.entries :
    entry.version if entry.version != null
  ]
}
