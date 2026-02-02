# Service Principal configuration for unit testing - Minimal
resource "microsoft365_graph_beta_applications_service_principal" "test_minimal" {
  app_id = "11111111-1111-1111-1111-111111111111"
}
