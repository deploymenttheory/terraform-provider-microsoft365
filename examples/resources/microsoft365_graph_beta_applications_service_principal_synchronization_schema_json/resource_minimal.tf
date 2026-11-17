# Manages the provisioning schema of an existing synchronization job.
#
# Export the job's current schema with
#   GET https://graph.microsoft.com/beta/servicePrincipals/{id}/synchronization/jobs/{jobId}/schema
# and save its "synchronizationRules" and "directories" arrays as JSON files next to this configuration.
resource "microsoft365_graph_beta_applications_service_principal_synchronization_schema_json" "example" {
  service_principal_id = "00000000-0000-0000-0000-000000000000"
  job_id               = "scim.00000000000000000000000000000000.00000000-0000-0000-0000-000000000000"

  synchronization_rules_json = file("${path.module}/synchronization_rules.json")
  directories_json           = file("${path.module}/directories.json")
}
