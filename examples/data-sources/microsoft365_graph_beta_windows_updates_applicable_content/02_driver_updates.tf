# Example: Get only driver updates for a deployment audience

data "microsoft365_graph_beta_windows_updates_applicable_content" "drivers" {
  audience_id        = "12345678-1234-1234-1234-123456789012"
  catalog_entry_type = "driver"
}

# Output driver count
output "driver_count" {
  value = length(data.microsoft365_graph_beta_windows_updates_applicable_content.drivers.applicable_content)
}

# Output driver manufacturers
output "driver_manufacturers" {
  value = distinct([
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.drivers.applicable_content :
    content.catalog_entry.manufacturer
    if content.catalog_entry.manufacturer != null
  ])
}
