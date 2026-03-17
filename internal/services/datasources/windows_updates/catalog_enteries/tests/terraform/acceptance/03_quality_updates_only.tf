# Test 03: Get quality updates only
# This test retrieves only quality update catalog entries

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "test" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

output "quality_update_count" {
  value = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.test.entries)
}

output "security_updates" {
  value = [
    for entry in data.microsoft365_graph_beta_windows_updates_catalog_enteries.test.entries :
    entry if entry.quality_update_classification == "security"
  ]
}
