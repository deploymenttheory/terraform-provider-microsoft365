resource "microsoft365_graph_beta_applications_service_principal_app_role_assigned_to" "test_minimal" {
  resource_object_id                 = "22222222-2222-2222-2222-222222222222"
  app_role_id                        = "df021288-bdef-4463-88db-98f22de89214"
  target_service_principal_object_id = "11111111-1111-1111-1111-111111111111"
}

