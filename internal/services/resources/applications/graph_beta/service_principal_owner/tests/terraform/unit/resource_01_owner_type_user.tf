# Service Principal Owner configuration for unit testing - User owner type
resource "microsoft365_graph_beta_applications_service_principal_owner" "test_user" {
  service_principal_id = "11111111-1111-1111-1111-111111111111"
  owner_id             = "22222222-2222-2222-2222-222222222222"
  owner_object_type    = "User"
}
