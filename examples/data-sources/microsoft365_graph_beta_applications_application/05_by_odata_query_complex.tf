# Look up an application using a complex OData query with multiple conditions
# This example finds a single-tenant application with a specific name pattern
data "microsoft365_graph_beta_applications_application" "by_odata_complex" {
  odata_query = "displayName eq 'Production API' and signInAudience eq 'AzureADMyOrg'"
}

# Output comprehensive application details
output "app_by_odata_complex" {
  value = {
    id                            = data.microsoft365_graph_beta_applications_application.by_odata_complex.id
    app_id                        = data.microsoft365_graph_beta_applications_application.by_odata_complex.app_id
    display_name                  = data.microsoft365_graph_beta_applications_application.by_odata_complex.display_name
    sign_in_audience              = data.microsoft365_graph_beta_applications_application.by_odata_complex.sign_in_audience
    identifier_uris               = data.microsoft365_graph_beta_applications_application.by_odata_complex.identifier_uris
    is_fallback_public_client     = data.microsoft365_graph_beta_applications_application.by_odata_complex.is_fallback_public_client
    is_device_only_auth_supported = data.microsoft365_graph_beta_applications_application.by_odata_complex.is_device_only_auth_supported
  }
}
