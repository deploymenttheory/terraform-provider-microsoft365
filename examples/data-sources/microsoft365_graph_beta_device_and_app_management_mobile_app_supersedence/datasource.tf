# This example demonstrates how to retrieve information about mobile app supersedence relationships in Intune

# Example 1: Retrieve a specific supersedence relationship by ID
data "microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence" "by_id" {
  id = "00000000-0000-0000-0000-000000000001_00000000-0000-0000-0000-000000000002"
}

# Example 2: Use the data source to output information about the supersedence relationship
output "supersedence_relationship_details" {
  value = {
    id                     = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.id
    source_id              = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.source_id
    source_display_name    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.source_display_name
    source_display_version = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.source_display_version
    target_id              = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.target_id
    target_display_name    = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.target_display_name
    target_display_version = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.target_display_version
    supersedence_type      = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.supersedence_type
    superseded_app_count   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.superseded_app_count
    superseding_app_count  = data.microsoft365_graph_beta_device_and_app_management_mobile_app_supersedence.by_id.superseding_app_count
  }
}

# Note: Replace the ID with an actual supersedence relationship ID from your Intune environment
# The ID format is typically source_app_id_target_app_id 