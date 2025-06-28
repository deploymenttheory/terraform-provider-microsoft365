# Example: Retrieve all Cloud PC audit events

data "microsoft365_graph_beta_windows_365_cloud_pc_audit_event" "all" {
  filter_type = "all"
}

# Output: List all audit event IDs
output "all_audit_event_ids" {
  value = [for event in data.microsoft365_graph_beta_windows_365_cloud_pc_audit_event.all.items : event.id]
}

# Output: Show all details for the first audit event (if present)
output "first_audit_event_details" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_audit_event.all.items[0]
}

# Output: Show nested actor and resource details for the first audit event
output "first_audit_event_actor" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_audit_event.all.items[0].actor
}

output "first_audit_event_resources" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_audit_event.all.items[0].resources
}

# Example: Retrieve a specific audit event by ID (using the ID from the sample JSON)
data "microsoft365_graph_beta_windows_365_cloud_pc_audit_event" "by_id" {
  filter_type  = "id"
  filter_value = "250473f5-029f-4037-813d-ba4768201d61"
}

output "audit_event_by_id" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_audit_event.by_id.items[0]
}

# Example: Retrieve audit events by display name substring (using the displayName from the sample JSON)
data "microsoft365_graph_beta_windows_365_cloud_pc_audit_event" "by_display_name" {
  filter_type  = "display_name"
  filter_value = "Delete OnPremisesConnection"
}

output "audit_events_by_display_name" {
  value = data.microsoft365_graph_beta_windows_365_cloud_pc_audit_event.by_display_name.items
} 