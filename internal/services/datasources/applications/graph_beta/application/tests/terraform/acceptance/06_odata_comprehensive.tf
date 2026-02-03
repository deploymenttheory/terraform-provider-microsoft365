resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name     = "acc-test-app-comprehensive-${random_string.suffix.result}"
  sign_in_audience = "AzureADMyOrg"
  
  tags = [
    "terraform-test",
    "comprehensive-test",
    "MyCustomTag"
  ]
}

data "microsoft365_graph_beta_applications_application" "odata_comprehensive" {
  odata_query = "tags/any(t:t eq 'MyCustomTag') and displayName eq '${microsoft365_graph_beta_applications_application.test.display_name}'"
  
  depends_on = [microsoft365_graph_beta_applications_application.test]
}
