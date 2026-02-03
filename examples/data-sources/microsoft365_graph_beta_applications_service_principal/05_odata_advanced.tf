# Retrieve a service principal using OData query with multiple filter conditions
data "microsoft365_graph_beta_applications_service_principal" "odata_advanced" {
  odata_query = "servicePrincipalType eq 'Application' and accountEnabled eq true"
}
