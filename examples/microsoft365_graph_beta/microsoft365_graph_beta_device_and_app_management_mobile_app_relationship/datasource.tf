# This example demonstrates how to retrieve information about mobile app relationships in Intune

# Example 1: Retrieve all mobile app relationships
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "all" {
  filter_type = "all"
}

# Example 2: Retrieve a specific mobile app relationship by ID
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_id" {
  filter_type  = "id"
  filter_value = "00000000-0000-0000-0000-000000000001_00000000-0000-0000-0000-000000000002" # Replace with a valid relationship ID
}

# Example 3: Retrieve mobile app relationships by source app ID
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_source" {
  filter_type  = "source_id"
  filter_value = "00000000-0000-0000-0000-000000000001" # Replace with a valid source app ID
}

# Example 4: Retrieve mobile app relationships by target app ID
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_target" {
  filter_type  = "target_id"
  filter_value = "00000000-0000-0000-0000-000000000002" # Replace with a valid target app ID
}

# Example 5: Use OData filtering to find relationships where the target is a child
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "child_relationships" {
  filter_type  = "odata"
  odata_filter = "targetType eq 'child'"
  odata_top    = 10 # Limit to 10 results
}

# Example 6: Use OData filtering to find relationships where the target is a parent
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "parent_relationships" {
  filter_type  = "odata"
  odata_filter = "targetType eq 'parent'"
  odata_top    = 10 # Limit to 10 results
}

# Example 7: Output information about a specific relationship
output "relationship_details" {
  value = {
    id                  = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_id.items[0].id
    source_id           = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_id.items[0].source_id
    source_display_name = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_id.items[0].source_display_name
    target_id           = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_id.items[0].target_id
    target_display_name = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_id.items[0].target_display_name
    target_type         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_id.items[0].target_type
  }
} 