# Example 1: Get all Windows Update catalog entries
data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "all_updates" {
  filter_type = "all"
}

output "all_catalog_entries" {
  value = data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.all_updates.entries
}

# Example 2: Filter by catalog entry type (quality updates only)
data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "quality_updates" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

output "quality_update_count" {
  value = length(data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.quality_updates.entries)
}

# Example 3: Filter by catalog entry type (feature updates only)
data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "feature_updates" {
  filter_type  = "catalog_entry_type"
  filter_value = "featureUpdate"
}

output "feature_updates_list" {
  value = [
    for entry in data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.feature_updates.entries : {
      id      = entry.id
      name    = entry.display_name
      version = entry.version
      release = entry.release_date_time
    }
  ]
}

# Example 4: Filter by display name (search for specific update)
data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "security_updates" {
  filter_type  = "display_name"
  filter_value = "SecurityUpdate"
}

output "security_updates" {
  value = data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.security_updates.entries
}

# Example 5: Get a specific catalog entry by ID
data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "specific_update" {
  filter_type  = "id"
  filter_value = "c1dec151-c151-c1de-51c1-dec151c1dec1"
}

output "specific_update_details" {
  value = length(data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.specific_update.entries) > 0 ? data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.specific_update.entries[0] : null
}

# Example 6: Get quality updates with CVE information
data "microsoft365_graph_beta_device_management_windows_update_catalog_enteries" "quality_updates_with_cves" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

output "updates_with_exploited_cves" {
  value = [
    for entry in data.microsoft365_graph_beta_device_management_windows_update_catalog_enteries.quality_updates_with_cves.entries :
    entry if entry.cve_severity_information != null && length(entry.cve_severity_information.exploited_cves) > 0
  ]
}
