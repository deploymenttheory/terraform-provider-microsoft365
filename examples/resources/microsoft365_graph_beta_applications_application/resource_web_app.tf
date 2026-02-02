resource "microsoft365_graph_beta_applications_application" "web_app" {
  display_name     = "my-web-application"
  description      = "Web application with OIDC authentication"
  sign_in_audience = "AzureADMyOrg"

  identifier_uris = [
    "https://mycompany.com/my-web-app"
  ]

  web = {
    home_page_url = "https://my-web-app.azurewebsites.net"
    logout_url    = "https://my-web-app.azurewebsites.net/signout"
    redirect_uris = [
      "https://my-web-app.azurewebsites.net/signin-oidc"
    ]
    implicit_grant_settings = {
      enable_access_token_issuance = false
      enable_id_token_issuance     = true
    }
  }

  required_resource_access = []
}
