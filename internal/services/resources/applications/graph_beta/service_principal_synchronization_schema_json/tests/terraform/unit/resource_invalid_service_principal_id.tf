resource "microsoft365_graph_beta_applications_service_principal_synchronization_schema_json" "test" {
  service_principal_id = "not-a-guid"
  job_id               = "scim.00000000000000000000000000000001.22222222-2222-2222-2222-222222222222"

  synchronization_rules_json = jsonencode([])
  directories_json           = jsonencode([])
}
