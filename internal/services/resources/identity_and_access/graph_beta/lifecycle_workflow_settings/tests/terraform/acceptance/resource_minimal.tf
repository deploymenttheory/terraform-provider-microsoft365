# Minimal configuration: only the workflow schedule interval is managed.
# The email settings keep their current tenant values.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "test" {
  workflow_schedule_interval_in_hours = 5
}
