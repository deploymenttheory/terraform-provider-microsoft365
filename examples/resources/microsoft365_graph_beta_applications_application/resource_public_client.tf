resource "microsoft365_graph_beta_applications_application" "mobile_app" {
  display_name              = "my-mobile-application"
  description               = "Mobile or desktop application (public client)"
  sign_in_audience          = "AzureADMyOrg"
  is_fallback_public_client = true

  public_client = {
    redirect_uris = [
      "http://localhost",
      "ms-appx-web://microsoft.aad.brokerplugin/my-mobile-app"
    ]
  }

  required_resource_access = []
}
