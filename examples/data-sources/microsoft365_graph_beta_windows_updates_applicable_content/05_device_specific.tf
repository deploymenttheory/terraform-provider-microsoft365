# Example: Get applicable content for a specific device

data "microsoft365_graph_beta_windows_updates_applicable_content" "device_updates" {
  audience_id = "12345678-1234-1234-1234-123456789012"
  device_id   = "aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee"
}

# Output updates applicable to this specific device
output "device_applicable_updates" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.device_updates.applicable_content : {
      catalog_entry_id = content.catalog_entry_id
      display_name     = content.catalog_entry.display_name
      driver_class     = content.catalog_entry.driver_class
      manufacturer     = content.catalog_entry.manufacturer
    }
  ]
}
