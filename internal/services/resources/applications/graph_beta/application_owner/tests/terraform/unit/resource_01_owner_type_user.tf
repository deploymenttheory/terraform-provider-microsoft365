# Application Owner configuration for unit testing - User owner type
resource "microsoft365_graph_beta_applications_application_owner" "test_user" {
  application_id    = "11111111-1111-1111-1111-111111111111"
  owner_id          = "user-11111111-1111-1111-1111-111111111111"
  owner_object_type = "User"
}
