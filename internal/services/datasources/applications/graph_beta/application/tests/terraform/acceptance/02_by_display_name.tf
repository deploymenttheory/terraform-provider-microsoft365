resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name     = "acc-test-app-displayname-${random_string.suffix.result}"
  sign_in_audience = "AzureADMyOrg"
  
  tags = [
    "terraform-test"
  ]
}

data "microsoft365_graph_beta_applications_application" "by_display_name" {
  display_name = microsoft365_graph_beta_applications_application.test.display_name
  
  depends_on = [microsoft365_graph_beta_applications_application.test]
}
