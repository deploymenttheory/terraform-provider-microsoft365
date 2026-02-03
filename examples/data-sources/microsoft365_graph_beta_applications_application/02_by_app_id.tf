# Look up an application by its application (client) ID
# This is useful when you know the app ID but not the object ID
data "microsoft365_graph_beta_applications_application" "by_app_id" {
  app_id = "00000003-0000-0000-c000-000000000000" # Example: Microsoft Graph app ID
}

# Output the application details
output "app_by_app_id" {
  value = {
    id                    = data.microsoft365_graph_beta_applications_application.by_app_id.id
    display_name          = data.microsoft365_graph_beta_applications_application.by_app_id.display_name
    sign_in_audience      = data.microsoft365_graph_beta_applications_application.by_app_id.sign_in_audience
    identifier_uris       = data.microsoft365_graph_beta_applications_application.by_app_id.identifier_uris
    publisher_domain      = data.microsoft365_graph_beta_applications_application.by_app_id.publisher_domain
  }
}
