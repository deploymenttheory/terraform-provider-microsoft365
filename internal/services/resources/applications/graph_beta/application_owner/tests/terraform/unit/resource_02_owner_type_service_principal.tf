# Application Owner configuration for unit testing - Service Principal owner type
resource "microsoft365_graph_beta_applications_application_owner" "test_service_principal" {
  application_id    = "33333333-3333-3333-3333-333333333333"
  owner_id          = "44444444-4444-4444-4444-444444444444"
  owner_object_type = "ServicePrincipal"
}
