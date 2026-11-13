# Each lookup must match exactly one relationship (dependency or supersedence).

# Example 1: Look up a relationship by its ID ({source_id}_{target_id})
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_relationship_id" {
  relationship_id = "00000000-0000-0000-0000-000000000001_00000000-0000-0000-0000-000000000002"
}

# Example 2: Look up the relationship between a source app and a target app
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_target_id" {
  source_id = "00000000-0000-0000-0000-000000000001"
  target_id = "00000000-0000-0000-0000-000000000002"
}

# Example 3: Look up a source app's relationship by the target app's display name
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "by_target_display_name" {
  source_id           = "00000000-0000-0000-0000-000000000001"
  target_display_name = "Microsoft Visual C++ Redistributable"
}

# Example 4: Use an OData filter on the source app's relationships
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "dependency" {
  source_id   = "00000000-0000-0000-0000-000000000001"
  odata_query = "isof('microsoft.graph.mobileAppDependency')"
}

# Example 5: Find the app that depends on or supersedes this app (an incoming relationship)
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "parent" {
  source_id   = "00000000-0000-0000-0000-000000000002"
  odata_query = "targetType eq 'parent'"
}

# Example 6: Narrow a lookup with target_type
data "microsoft365_graph_beta_device_and_app_management_mobile_app_relationship" "child_only" {
  source_id   = "00000000-0000-0000-0000-000000000001"
  target_id   = "00000000-0000-0000-0000-000000000002"
  target_type = "child"
}

output "relationship_details" {
  value = {
    id                  = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.id
    relationship_type   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.relationship_type
    source_display_name = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.source_display_name
    target_display_name = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.target_display_name
    target_type         = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.target_type
    dependency_type     = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.dependency_type
    supersedence_type   = data.microsoft365_graph_beta_device_and_app_management_mobile_app_relationship.by_target_id.supersedence_type
  }
}
