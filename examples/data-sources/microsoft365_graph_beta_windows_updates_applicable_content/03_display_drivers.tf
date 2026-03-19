# Example: Get only display driver updates

data "microsoft365_graph_beta_windows_updates_applicable_content" "display_drivers" {
  audience_id        = "12345678-1234-1234-1234-123456789012"
  catalog_entry_type = "driver"
  driver_class       = "Display"
}

# Output display driver details
output "display_drivers" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.display_drivers.applicable_content : {
      display_name = content.catalog_entry.display_name
      manufacturer = content.catalog_entry.manufacturer
      version      = content.catalog_entry.version
      provider     = content.catalog_entry.provider
      device_count = length(content.matched_devices)
    }
  ]
}
