# Example: Get only content that has matched devices

data "microsoft365_graph_beta_windows_updates_applicable_content" "with_matches" {
  audience_id        = "12345678-1234-1234-1234-123456789012"
  include_no_matches = false
}

# Output content with device matches
output "content_with_matches" {
  value = [
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.with_matches.applicable_content : {
      display_name    = content.catalog_entry.display_name
      matched_devices = length(content.matched_devices)
      device_ids      = [for device in content.matched_devices : device.device_id]
    }
  ]
}

# Output total devices that have applicable content
output "total_devices_with_updates" {
  value = length(distinct(flatten([
    for content in data.microsoft365_graph_beta_windows_updates_applicable_content.with_matches.applicable_content : [
      for device in content.matched_devices : device.device_id
    ]
  ])))
}
