# Example: Get only quality/security updates

data "microsoft365_graph_beta_windows_updates_applicable_content" "quality_updates" {
  audience_id        = "12345678-1234-1234-1234-123456789012"
  catalog_entry_type = "quality"
}

# Output quality update details
output "quality_updates" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.quality_updates.applicable_content : {
      display_name         = content.catalog_entry.display_name
      release_date         = content.catalog_entry.release_date_time
      deployable_until     = content.catalog_entry.deployable_until_date_time
      matched_device_count = length(content.matched_devices)
    }
  ]
}
