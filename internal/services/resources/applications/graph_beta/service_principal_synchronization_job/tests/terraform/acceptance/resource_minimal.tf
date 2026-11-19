resource "microsoft365_graph_beta_applications_service_principal_synchronization_job" "test" {
  service_principal_id = "%s"
  template_id          = "scim"
  enabled              = false
}
