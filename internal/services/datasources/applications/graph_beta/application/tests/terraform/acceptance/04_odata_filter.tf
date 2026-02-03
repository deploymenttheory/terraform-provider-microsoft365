resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

resource "microsoft365_graph_beta_applications_application" "test" {
  display_name     = "acc-test-app-odata-${random_string.suffix.result}"
  sign_in_audience = "AzureADMyOrg"

  tags = [
    "terraform-test",
    "odata-test"
  ]
}

data "microsoft365_graph_beta_applications_application" "odata_filter" {
  odata_query = "displayName eq '${microsoft365_graph_beta_applications_application.test.display_name}' and signInAudience eq 'AzureADMyOrg'"

  depends_on = [microsoft365_graph_beta_applications_application.test]
}
