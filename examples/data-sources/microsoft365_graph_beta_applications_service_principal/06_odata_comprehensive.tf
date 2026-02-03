# Retrieve a service principal using comprehensive OData query filters
# This example demonstrates filtering for SAML-based service principals
data "microsoft365_graph_beta_applications_service_principal" "odata_comprehensive" {
  odata_query = "preferredSingleSignOnMode eq 'saml' and servicePrincipalType eq 'Application'"
}
