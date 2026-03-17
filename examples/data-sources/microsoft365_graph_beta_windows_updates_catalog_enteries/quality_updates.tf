# Get quality updates only
# This retrieves only quality update catalog entries (security and non-security updates)

data "microsoft365_graph_beta_windows_updates_catalog_enteries" "quality_updates" {
  filter_type  = "catalog_entry_type"
  filter_value = "qualityUpdate"
}

output "quality_update_count" {
  description = "Number of quality updates available"
  value       = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries)
}

output "security_updates" {
  description = "All security quality updates"
  value = [
    for entry in data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries :
    {
      id                            = entry.id
      display_name                  = entry.display_name
      short_name                    = entry.short_name
      catalog_name                  = entry.catalog_name
      quality_update_classification = entry.quality_update_classification
      is_expeditable                = entry.is_expeditable
      release_date_time             = entry.release_date_time
    } if entry.quality_update_classification == "security"
  ]
}

output "latest_quality_update" {
  description = "The most recent quality update with CVE information"
  value = length(data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries) > 0 ? {
    id                            = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].id
    display_name                  = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].display_name
    short_name                    = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].short_name
    quality_update_classification = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].quality_update_classification
    is_expeditable                = data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].is_expeditable
    cve_max_severity              = try(data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].cve_severity_information.max_severity, null)
    cve_max_base_score            = try(data.microsoft365_graph_beta_windows_updates_catalog_enteries.quality_updates.entries[0].cve_severity_information.max_base_score, null)
  } : null
}
