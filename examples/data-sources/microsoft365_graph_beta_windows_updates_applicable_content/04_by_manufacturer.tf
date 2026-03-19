# Example: Get updates from a specific manufacturer

data "microsoft365_graph_beta_windows_updates_applicable_content" "intel_updates" {
  audience_id  = "12345678-1234-1234-1234-123456789012"
  manufacturer = "Intel"
}

# Output Intel update details
output "intel_updates" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.intel_updates.applicable_content : {
      display_name = content.catalog_entry.display_name
      driver_class = content.catalog_entry.driver_class
      version      = content.catalog_entry.version
      release_date = content.catalog_entry.release_date_time
    }
  ]
}
