# Run scheduled lifecycle workflows every 6 hours.
# The email settings keep their current tenant values.
resource "microsoft365_graph_beta_identity_and_access_lifecycle_workflow_settings" "example" {
  workflow_schedule_interval_in_hours = 6
}
