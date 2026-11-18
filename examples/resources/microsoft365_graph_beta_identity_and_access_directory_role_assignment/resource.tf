resource "microsoft365_graph_beta_identity_and_access_directory_role_assignment" "example" {
  role_definition_id = "88d8e3e3-8f55-4a1e-953a-9b9898b8876b" # Directory Readers template ID
  principal_id       = "11111111-1111-4111-8111-111111111111" # User, service principal, or role-assignable group object ID
}
