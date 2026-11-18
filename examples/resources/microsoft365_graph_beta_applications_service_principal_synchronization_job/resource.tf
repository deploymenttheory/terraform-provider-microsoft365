resource "microsoft365_graph_beta_applications_service_principal_synchronization_job" "example" {
  service_principal_id = "11111111-1111-4111-8111-111111111111"
  template_id          = "scim"
  enabled              = false
}

# The sibling schema resource consumes example.job_id, not example.id.
# Configure credentials and schema before explicitly enabling provisioning.
