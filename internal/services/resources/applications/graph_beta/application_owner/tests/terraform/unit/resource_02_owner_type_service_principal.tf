# Application Owner configuration for unit testing - Service Principal owner type
resource "microsoft365_graph_beta_applications_application_owner" "test_service_principal" {
  application_id    = "22222222-2222-2222-2222-222222222222"
  owner_id          = "sp-11111111-1111-1111-1111-111111111111"
  owner_object_type = "ServicePrincipal"
}
