# Test with default app role (when no specific role is defined)
resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "test_default_role" {
  resource_object_id                 = "22222222-2222-2222-2222-222222222222"
  app_role_id                        = "00000000-0000-0000-0000-000000000000"
  target_service_principal_object_id = "11111111-1111-1111-1111-111111111111"
}

