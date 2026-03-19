# Example: Get all applicable content for a deployment audience

data "microsoft365_graph_beta_windows_updates_applicable_content" "all" {
  audience_id = "12345678-1234-1234-1234-123456789012"
}

# Output the count of applicable content
output "total_applicable_content" {
  value = length(data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content)
}

# Output details of the first applicable content entry
output "first_content" {
  value = length(data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content) > 0 ? {
    catalog_entry_id = data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content[0].catalog_entry_id
    display_name     = data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content[0].catalog_entry.display_name
    matched_devices  = length(data.microsoft365_graph_beta_windows_updates_applicable_content.all.applicable_content[0].matched_devices)
  } : null
}
