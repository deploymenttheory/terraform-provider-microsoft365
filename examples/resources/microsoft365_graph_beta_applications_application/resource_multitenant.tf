resource "microsoft365_graph_beta_applications_application" "multitenant" {
  display_name     = "my-multitenant-app"
  description      = "Multitenant application with personal Microsoft account support"
  sign_in_audience = "AzureADandPersonalMicrosoftAccount"

  identifier_uris = [
    "https://mycompany.com/my-multitenant-app"
  ]

  web = {
    home_page_url = "https://contoso.com"
    redirect_uris = [
      "https://contoso.com/signin-oidc"
    ]
    implicit_grant_settings = {
      enable_access_token_issuance = false
      enable_id_token_issuance     = true
    }
  }

  spa = {
    redirect_uris = [
      "https://contoso.com/spa"
    ]
  }

  required_resource_access = []

  tags = [
    "multitenant",
    "production"
  ]
}
