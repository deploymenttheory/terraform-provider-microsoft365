resource "microsoft365_graph_beta_applications_oauth2_permission_grant" "tenant_consent" {
  client_id    = "11111111-1111-4111-8111-111111111111" # Client service principal object ID
  resource_id  = "22222222-2222-4222-8222-222222222222" # API service principal object ID
  consent_type = "AllPrincipals"
  scopes       = ["User.Read"]
}

resource "microsoft365_graph_beta_applications_oauth2_permission_grant" "user_consent" {
  client_id    = "11111111-1111-4111-8111-111111111111"
  resource_id  = "22222222-2222-4222-8222-222222222222"
  consent_type = "Principal"
  principal_id = "33333333-3333-4333-8333-333333333333" # User object ID
  scopes       = ["User.Read"]
}
