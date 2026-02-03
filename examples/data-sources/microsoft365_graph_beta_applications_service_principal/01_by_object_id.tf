# Retrieve a service principal by its Object ID
data "microsoft365_graph_beta_applications_service_principal" "by_id" {
  object_id = "3b6f95b0-2064-4cc9-b5e5-1ab72af707b3"
}
