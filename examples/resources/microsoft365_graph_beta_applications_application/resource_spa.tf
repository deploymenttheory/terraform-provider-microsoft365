resource "microsoft365_graph_beta_applications_application" "spa" {
  display_name     = "my-single-page-app"
  description      = "Single Page Application (React, Angular, Vue)"
  sign_in_audience = "AzureADMultipleOrgs"

  identifier_uris = [
    "https://mycompany.com/my-spa"
  ]

  spa = {
    redirect_uris = [
      "http://localhost:3000",
      "https://my-spa.azurestaticapps.net"
    ]
  }

  required_resource_access = []
}
