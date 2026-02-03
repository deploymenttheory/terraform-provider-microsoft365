resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name     = "acc-test-app-advanced-${random_string.suffix.result}"
  sign_in_audience = "AzureADMyOrg"
  
  tags = [
    "terraform-test"
  ]
}

data "microsoft365_graph_beta_applications_application" "odata_advanced" {
  odata_query = "appId eq '${microsoft365_graph_beta_applications_application.test.app_id}'"
  
  depends_on = [microsoft365_graph_beta_applications_application.test]
}
