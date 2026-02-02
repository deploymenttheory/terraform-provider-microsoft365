# Service Principal Owner configuration for unit testing - Service Principal owner type
resource "microsoft365_graph_beta_applications_service_principal_owner" "test_service_principal" {
  service_principal_id = "22222222-2222-2222-2222-222222222222"
  owner_id             = "33333333-3333-3333-3333-333333333333"
  owner_object_type    = "ServicePrincipal"
}
